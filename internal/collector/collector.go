package collector

import (
	"crypto/x509"
	"encoding/pem"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"sync"

	"github.com/crewjam/saml"
	"github.com/go-logr/logr"
	"github.com/prometheus/client_golang/prometheus"
)

type Collector struct {
	logger     logr.Logger
	httpClient *http.Client
	urls       []*url.URL
	collectors []prometheus.Collector
}

func New(logger logr.Logger, httpClient *http.Client, urls []*url.URL, collectors ...prometheus.Collector) *Collector {
	return &Collector{
		logger:     logger,
		httpClient: httpClient,
		urls:       urls,
		collectors: collectors,
	}
}

func (c *Collector) Describe(ch chan<- *prometheus.Desc) {
	notAfter.Describe(ch)
	certParse.Describe(ch)
	metadataParse.Describe(ch)
	for _, c := range c.collectors {
		c.Describe(ch)
	}
}

func (c *Collector) Collect(ch chan<- prometheus.Metric) {
	var wg sync.WaitGroup

	for _, u := range c.urls {
		wg.Add(1)
		go func(u *url.URL) {
			defer wg.Done()

			req := &http.Request{
				URL:    u,
				Method: http.MethodGet,
			}

			res, err := c.httpClient.Do(req)
			if err != nil {
				c.logger.Error(err, "http request failed", "url", u.String())
				return
			}

			b, err := io.ReadAll(res.Body)
			if err != nil {
				c.logger.Error(err, "read body from http response failed", "url", u.String())
				incSAMLParse(ch, u)
				return
			}

			entity := &saml.EntityDescriptor{}
			if xml.Unmarshal(b, entity) != nil {
				c.logger.Error(err, "parse saml metadata failed", "url", u.String())
				incSAMLParse(ch, u)
				return
			}

			for _, ssoDescriptor := range entity.IDPSSODescriptors {
				for _, keyDescriptor := range ssoDescriptor.KeyDescriptors {
					c.collectExpireMetric(ch, u, entity, keyDescriptor)
				}
			}

			for _, ssoDescriptor := range entity.SPSSODescriptors {
				for _, keyDescriptor := range ssoDescriptor.KeyDescriptors {
					c.collectExpireMetric(ch, u, entity, keyDescriptor)
				}
			}
		}(u)
	}

	wg.Wait()

	for _, c := range c.collectors {
		c.Collect(ch)
	}
}

func incSAMLParse(ch chan<- prometheus.Metric, url *url.URL) {
	metric := metadataParse.With(prometheus.Labels{
		"url": url.String(),
	})

	metric.Inc()
	ch <- metric
}

func incCertParse(ch chan<- prometheus.Metric, url string, entity *saml.EntityDescriptor, keyDescriptor saml.KeyDescriptor) {
	metric := certParse.With(prometheus.Labels{
		"url":      url,
		"entityid": entity.EntityID,
		"use":      keyDescriptor.Use,
	})

	metric.Inc()
	ch <- metric
}

func (c *Collector) collectExpireMetric(ch chan<- prometheus.Metric, url *url.URL, entity *saml.EntityDescriptor, keyDescriptor saml.KeyDescriptor) {
	u := url.String()
	for _, x509Cert := range keyDescriptor.KeyInfo.X509Data.X509Certificates {
		block, _ := pem.Decode([]byte(fmt.Sprintf("-----BEGIN CERTIFICATE-----\n%s\n-----END CERTIFICATE-----", x509Cert.Data)))

		if block == nil {
			c.logger.Info("found empty pem block", "entityid", entity.EntityID)
			incCertParse(ch, u, entity, keyDescriptor)
			continue
		}

		cert, err := x509.ParseCertificate(block.Bytes)

		if err != nil {
			c.logger.Error(err, "failed to parse x509 certificate", "entityid", entity.EntityID)
			incCertParse(ch, u, entity, keyDescriptor)
			continue
		}

		labels := prometheus.Labels{
			"url":           u,
			"entityid":      entity.EntityID,
			"use":           keyDescriptor.Use,
			"serial_number": cert.SerialNumber.String(),
			"issuer_C":      strings.Join(cert.Issuer.Country, ","),
			"issuer_CN":     cert.Issuer.CommonName,
			"issuer_L":      strings.Join(cert.Issuer.Locality, ","),
			"issuer_O":      strings.Join(cert.Issuer.Organization, ","),
			"issuer_ST":     strings.Join(cert.Issuer.StreetAddress, ","),
			"subject_C":     strings.Join(cert.Subject.Country, ","),
			"subject_CN":    cert.Subject.CommonName,
			"subject_L":     strings.Join(cert.Subject.Locality, ","),
			"subject_O":     strings.Join(cert.Subject.Organization, ","),
		}

		metricNotAfter := notAfter.With(labels)
		metricNotAfter.Set(float64(cert.NotAfter.Unix()))
		ch <- metricNotAfter

		metricNotBefore := notBefore.With(labels)
		metricNotBefore.Set(float64(cert.NotBefore.Unix()))
		ch <- metricNotBefore
	}
}
