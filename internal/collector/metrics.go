package collector

import "github.com/prometheus/client_golang/prometheus"

var (
	notBefore = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "saml_x509_cert_not_valid_before",
		Help: "SAML X509 certificate validity period start",
	}, []string{"entityid", "use", "serial_number", "issuer_C", "issuer_CN", "issuer_L", "issuer_O", "issuer_ST", "subject_C", "subject_CN", "subject_L", "subject_O"})

	notAfter = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "saml_x509_cert_not_valid_after",
		Help: "SAML X509 certificate expiration date",
	}, []string{"entityid", "use", "serial_number", "issuer_C", "issuer_CN", "issuer_L", "issuer_O", "issuer_ST", "subject_C", "subject_CN", "subject_L", "subject_O"})

	certParse = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "saml_x509_read_errors_total",
		Help: "Errors encountered while parsing SAML X509 certificates",
	}, []string{"entityid", "use"})

	metadataParse = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "saml_metadata_errors_total",
		Help: "Errors encountered while parsing SAML metadata",
	}, []string{"url"})
)
