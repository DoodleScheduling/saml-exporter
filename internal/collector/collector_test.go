//go:build unit

package collector

import (
	"errors"
	"io"
	"net/http"
	"net/url"
	"strings"
	"testing"
	"testing/iotest"

	"github.com/doodlescheduling/saml-exporter/internal/transport"
	"github.com/go-logr/logr"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/testutil"
	"github.com/tj/assert"
)

type test struct {
	name     string
	response http.RoundTripper
	expected string
}

func TestInitializeMetrics(t *testing.T) {
	var tests = []test{
		{
			name: "Simple service provider SAML metadata",
			response: transport.NewMock(&http.Response{
				StatusCode: 200,
				Body: io.NopCloser(strings.NewReader(`
				<EntityDescriptor
				xmlns="urn:oasis:names:tc:SAML:2.0:metadata"
				entityID="example-entity-id">
				<SPSSODescriptor
					AuthnRequestsSigned="false"
					WantAssertionsSigned="false"
					protocolSupportEnumeration=
						"urn:oasis:names:tc:SAML:2.0:protocol">
					<KeyDescriptor use="signing">
						<KeyInfo xmlns="http://www.w3.org/2000/09/xmldsig#">
							<X509Data>
								<X509Certificate>
			MIICYDCCAgqgAwIBAgICBoowDQYJKoZIhvcNAQEEBQAwgZIxCzAJBgNVBAYTAlVTMRMwEQYDVQQI
			EwpDYWxpZm9ybmlhMRQwEgYDVQQHEwtTYW50YSBDbGFyYTEeMBwGA1UEChMVU3VuIE1pY3Jvc3lz
			dGVtcyBJbmMuMRowGAYDVQQLExFJZGVudGl0eSBTZXJ2aWNlczEcMBoGA1UEAxMTQ2VydGlmaWNh
			dGUgTWFuYWdlcjAeFw0wNjExMDIxOTExMzRaFw0xMDA3MjkxOTExMzRaMDcxEjAQBgNVBAoTCXNp
			cm9lLmNvbTEhMB8GA1UEAxMYbG9hZGJhbGFuY2VyLTkuc2lyb2UuY29tMIGfMA0GCSqGSIb3DQEB
			AQUAA4GNADCBiQKBgQCjOwa5qoaUuVnknqf5pdgAJSEoWlvx/jnUYbkSDpXLzraEiy2UhvwpoBgB
			EeTSUaPPBvboCItchakPI6Z/aFdH3Wmjuij9XD8r1C+q//7sUO0IGn0ORycddHhoo0aSdnnxGf9V
			tREaqKm9dJ7Yn7kQHjo2eryMgYxtr/Z5Il5F+wIDAQABo2AwXjARBglghkgBhvhCAQEEBAMCBkAw
			DgYDVR0PAQH/BAQDAgTwMB8GA1UdIwQYMBaAFDugITflTCfsWyNLTXDl7cMDUKuuMBgGA1UdEQQR
			MA+BDW1hbGxhQHN1bi5jb20wDQYJKoZIhvcNAQEEBQADQQB/6DOB6sRqCZu2OenM9eQR0gube85e
			nTTxU4a7x1naFxzYXK1iQ1vMARKMjDb19QEJIEJKZlDK4uS7yMlf1nFS
								</X509Certificate>
							</X509Data>
						</KeyInfo>
					</KeyDescriptor>
					<KeyDescriptor use="encryption">
						<KeyInfo xmlns="http://www.w3.org/2000/09/xmldsig#">
							<X509Data>
								<X509Certificate>
			MIICTDCCAfagAwIBAgICBo8wDQYJKoZIhvcNAQEEBQAwgZIxCzAJBgNVBAYTAlVTMRMwEQYDVQQI
			EwpDYWxpZm9ybmlhMRQwEgYDVQQHEwtTYW50YSBDbGFyYTEeMBwGA1UEChMVU3VuIE1pY3Jvc3lz
			dGVtcyBJbmMuMRowGAYDVQQLExFJZGVudGl0eSBTZXJ2aWNlczEcMBoGA1UEAxMTQ2VydGlmaWNh
			dGUgTWFuYWdlcjAeFw0wNjExMDcyMzU2MTdaFw0xMDA4MDMyMzU2MTdaMCMxITAfBgNVBAMTGGxv
			YWRiYWxhbmNlci05LnNpcm9lLmNvbTCBnzANBgkqhkiG9w0BAQEFAAOBjQAwgYkCgYEAw574iRU6
			HsSO4LXW/OGTXyfsbGv6XRVOoy3v+J1pZ51KKejcDjDJXNkKGn3/356AwIaqbcymWd59T0zSqYfR
			Hn+45uyjYxRBmVJseLpVnOXLub9jsjULfGx0yjH4w+KsZSZCXatoCHbj/RJtkzuZY6V9to/hkH3S
			InQB4a3UAgMCAwEAAaNgMF4wEQYJYIZIAYb4QgEBBAQDAgZAMA4GA1UdDwEB/wQEAwIE8DAfBgNV
			HSMEGDAWgBQ7oCE35Uwn7FsjS01w5e3DA1CrrjAYBgNVHREEETAPgQ1tYWxsYUBzdW4uY29tMA0G
			CSqGSIb3DQEBBAUAA0EAMlbfBg/ff0Xkv4DOR5LEqmfTZKqgdlD81cXynfzlF7XfnOqI6hPIA90I
			x5Ql0ejivIJAYcMGUyA+/YwJg2FGoA==
								</X509Certificate>
							</X509Data>
						</KeyInfo>
						<EncryptionMethod Algorithm=
							"https://www.w3.org/2001/04/xmlenc#aes128-cbc">
							<KeySize xmlns="https://www.w3.org/2001/04/xmlenc#">128</KeySize>
						</EncryptionMethod>
					</KeyDescriptor>
				</SPSSODescriptor>
			</EntityDescriptor>`)),
			}, nil),
			expected: `
			# HELP saml_x509_cert_not_valid_after SAML X509 certificate expiration date
			# TYPE saml_x509_cert_not_valid_after gauge
			saml_x509_cert_not_valid_after{entityid="example-entity-id",issuer_C="US",issuer_CN="Certificate Manager",issuer_L="Santa Clara",issuer_O="Sun Microsystems Inc.",issuer_ST="",serial_number="1674",subject_C="",subject_CN="loadbalancer-9.siroe.com",subject_L="",subject_O="siroe.com",url="//saml-metadata",use="signing"} 1.280430694e+09
			saml_x509_cert_not_valid_after{entityid="example-entity-id",issuer_C="US",issuer_CN="Certificate Manager",issuer_L="Santa Clara",issuer_O="Sun Microsystems Inc.",issuer_ST="",serial_number="1679",subject_C="",subject_CN="loadbalancer-9.siroe.com",subject_L="",subject_O="",url="//saml-metadata",use="encryption"} 1.280879777e+09
			
			# HELP saml_x509_cert_not_valid_before SAML X509 certificate validity period start
			# TYPE saml_x509_cert_not_valid_before gauge
			saml_x509_cert_not_valid_before{entityid="example-entity-id",issuer_C="US",issuer_CN="Certificate Manager",issuer_L="Santa Clara",issuer_O="Sun Microsystems Inc.",issuer_ST="",serial_number="1674",subject_C="",subject_CN="loadbalancer-9.siroe.com",subject_L="",subject_O="siroe.com",url="//saml-metadata",use="signing"} 1.162494694e+09
			saml_x509_cert_not_valid_before{entityid="example-entity-id",issuer_C="US",issuer_CN="Certificate Manager",issuer_L="Santa Clara",issuer_O="Sun Microsystems Inc.",issuer_ST="",serial_number="1679",subject_C="",subject_CN="loadbalancer-9.siroe.com",subject_L="",subject_O="",url="//saml-metadata",use="encryption"} 1.162943777e+09
			`,
		},
		{
			name: "Simple identity provider SAML metadata",
			response: transport.NewMock(&http.Response{
				StatusCode: 200,
				Body: io.NopCloser(strings.NewReader(`
				<EntityDescriptor
				xmlns="urn:oasis:names:tc:SAML:2.0:metadata"
				entityID="example-entity-id">
				<IDPSSODescriptor
					AuthnRequestsSigned="false"
					WantAssertionsSigned="false"
					protocolSupportEnumeration=
						"urn:oasis:names:tc:SAML:2.0:protocol">
					<KeyDescriptor use="signing">
						<KeyInfo xmlns="http://www.w3.org/2000/09/xmldsig#">
							<X509Data>
								<X509Certificate>
			MIICYDCCAgqgAwIBAgICBoowDQYJKoZIhvcNAQEEBQAwgZIxCzAJBgNVBAYTAlVTMRMwEQYDVQQI
			EwpDYWxpZm9ybmlhMRQwEgYDVQQHEwtTYW50YSBDbGFyYTEeMBwGA1UEChMVU3VuIE1pY3Jvc3lz
			dGVtcyBJbmMuMRowGAYDVQQLExFJZGVudGl0eSBTZXJ2aWNlczEcMBoGA1UEAxMTQ2VydGlmaWNh
			dGUgTWFuYWdlcjAeFw0wNjExMDIxOTExMzRaFw0xMDA3MjkxOTExMzRaMDcxEjAQBgNVBAoTCXNp
			cm9lLmNvbTEhMB8GA1UEAxMYbG9hZGJhbGFuY2VyLTkuc2lyb2UuY29tMIGfMA0GCSqGSIb3DQEB
			AQUAA4GNADCBiQKBgQCjOwa5qoaUuVnknqf5pdgAJSEoWlvx/jnUYbkSDpXLzraEiy2UhvwpoBgB
			EeTSUaPPBvboCItchakPI6Z/aFdH3Wmjuij9XD8r1C+q//7sUO0IGn0ORycddHhoo0aSdnnxGf9V
			tREaqKm9dJ7Yn7kQHjo2eryMgYxtr/Z5Il5F+wIDAQABo2AwXjARBglghkgBhvhCAQEEBAMCBkAw
			DgYDVR0PAQH/BAQDAgTwMB8GA1UdIwQYMBaAFDugITflTCfsWyNLTXDl7cMDUKuuMBgGA1UdEQQR
			MA+BDW1hbGxhQHN1bi5jb20wDQYJKoZIhvcNAQEEBQADQQB/6DOB6sRqCZu2OenM9eQR0gube85e
			nTTxU4a7x1naFxzYXK1iQ1vMARKMjDb19QEJIEJKZlDK4uS7yMlf1nFS
								</X509Certificate>
							</X509Data>
						</KeyInfo>
					</KeyDescriptor>
					<KeyDescriptor use="encryption">
						<KeyInfo xmlns="http://www.w3.org/2000/09/xmldsig#">
							<X509Data>
								<X509Certificate>
			MIICTDCCAfagAwIBAgICBo8wDQYJKoZIhvcNAQEEBQAwgZIxCzAJBgNVBAYTAlVTMRMwEQYDVQQI
			EwpDYWxpZm9ybmlhMRQwEgYDVQQHEwtTYW50YSBDbGFyYTEeMBwGA1UEChMVU3VuIE1pY3Jvc3lz
			dGVtcyBJbmMuMRowGAYDVQQLExFJZGVudGl0eSBTZXJ2aWNlczEcMBoGA1UEAxMTQ2VydGlmaWNh
			dGUgTWFuYWdlcjAeFw0wNjExMDcyMzU2MTdaFw0xMDA4MDMyMzU2MTdaMCMxITAfBgNVBAMTGGxv
			YWRiYWxhbmNlci05LnNpcm9lLmNvbTCBnzANBgkqhkiG9w0BAQEFAAOBjQAwgYkCgYEAw574iRU6
			HsSO4LXW/OGTXyfsbGv6XRVOoy3v+J1pZ51KKejcDjDJXNkKGn3/356AwIaqbcymWd59T0zSqYfR
			Hn+45uyjYxRBmVJseLpVnOXLub9jsjULfGx0yjH4w+KsZSZCXatoCHbj/RJtkzuZY6V9to/hkH3S
			InQB4a3UAgMCAwEAAaNgMF4wEQYJYIZIAYb4QgEBBAQDAgZAMA4GA1UdDwEB/wQEAwIE8DAfBgNV
			HSMEGDAWgBQ7oCE35Uwn7FsjS01w5e3DA1CrrjAYBgNVHREEETAPgQ1tYWxsYUBzdW4uY29tMA0G
			CSqGSIb3DQEBBAUAA0EAMlbfBg/ff0Xkv4DOR5LEqmfTZKqgdlD81cXynfzlF7XfnOqI6hPIA90I
			x5Ql0ejivIJAYcMGUyA+/YwJg2FGoA==
								</X509Certificate>
							</X509Data>
						</KeyInfo>
						<EncryptionMethod Algorithm=
							"https://www.w3.org/2001/04/xmlenc#aes128-cbc">
							<KeySize xmlns="https://www.w3.org/2001/04/xmlenc#">128</KeySize>
						</EncryptionMethod>
					</KeyDescriptor>
				</IDPSSODescriptor>
			</EntityDescriptor>`)),
			}, nil),
			expected: `
			# HELP saml_x509_cert_not_valid_after SAML X509 certificate expiration date
			# TYPE saml_x509_cert_not_valid_after gauge
			saml_x509_cert_not_valid_after{entityid="example-entity-id",issuer_C="US",issuer_CN="Certificate Manager",issuer_L="Santa Clara",issuer_O="Sun Microsystems Inc.",issuer_ST="",serial_number="1674",subject_C="",subject_CN="loadbalancer-9.siroe.com",subject_L="",subject_O="siroe.com",url="//saml-metadata",use="signing"} 1.280430694e+09
			saml_x509_cert_not_valid_after{entityid="example-entity-id",issuer_C="US",issuer_CN="Certificate Manager",issuer_L="Santa Clara",issuer_O="Sun Microsystems Inc.",issuer_ST="",serial_number="1679",subject_C="",subject_CN="loadbalancer-9.siroe.com",subject_L="",subject_O="",url="//saml-metadata",use="encryption"} 1.280879777e+09
			
			# HELP saml_x509_cert_not_valid_before SAML X509 certificate validity period start
			# TYPE saml_x509_cert_not_valid_before gauge
			saml_x509_cert_not_valid_before{entityid="example-entity-id",issuer_C="US",issuer_CN="Certificate Manager",issuer_L="Santa Clara",issuer_O="Sun Microsystems Inc.",issuer_ST="",serial_number="1674",subject_C="",subject_CN="loadbalancer-9.siroe.com",subject_L="",subject_O="siroe.com",url="//saml-metadata",use="signing"} 1.162494694e+09
			saml_x509_cert_not_valid_before{entityid="example-entity-id",issuer_C="US",issuer_CN="Certificate Manager",issuer_L="Santa Clara",issuer_O="Sun Microsystems Inc.",issuer_ST="",serial_number="1679",subject_C="",subject_CN="loadbalancer-9.siroe.com",subject_L="",subject_O="",url="//saml-metadata",use="encryption"} 1.162943777e+09
			`,
		},
		{
			name: "Identity provider metadata with an empty signing certificate",
			response: transport.NewMock(&http.Response{
				StatusCode: 200,
				Body: io.NopCloser(strings.NewReader(`
				<EntityDescriptor
				xmlns="urn:oasis:names:tc:SAML:2.0:metadata"
				entityID="example-entity-id">
				<IDPSSODescriptor
					AuthnRequestsSigned="false"
					WantAssertionsSigned="false"
					protocolSupportEnumeration=
						"urn:oasis:names:tc:SAML:2.0:protocol">
					<KeyDescriptor use="signing">
						<KeyInfo xmlns="http://www.w3.org/2000/09/xmldsig#">
							<X509Data>
								<X509Certificate></X509Certificate>
							</X509Data>
						</KeyInfo>
					</KeyDescriptor>
					<KeyDescriptor use="encryption">
						<KeyInfo xmlns="http://www.w3.org/2000/09/xmldsig#">
							<X509Data>
								<X509Certificate>
			MIICTDCCAfagAwIBAgICBo8wDQYJKoZIhvcNAQEEBQAwgZIxCzAJBgNVBAYTAlVTMRMwEQYDVQQI
			EwpDYWxpZm9ybmlhMRQwEgYDVQQHEwtTYW50YSBDbGFyYTEeMBwGA1UEChMVU3VuIE1pY3Jvc3lz
			dGVtcyBJbmMuMRowGAYDVQQLExFJZGVudGl0eSBTZXJ2aWNlczEcMBoGA1UEAxMTQ2VydGlmaWNh
			dGUgTWFuYWdlcjAeFw0wNjExMDcyMzU2MTdaFw0xMDA4MDMyMzU2MTdaMCMxITAfBgNVBAMTGGxv
			YWRiYWxhbmNlci05LnNpcm9lLmNvbTCBnzANBgkqhkiG9w0BAQEFAAOBjQAwgYkCgYEAw574iRU6
			HsSO4LXW/OGTXyfsbGv6XRVOoy3v+J1pZ51KKejcDjDJXNkKGn3/356AwIaqbcymWd59T0zSqYfR
			Hn+45uyjYxRBmVJseLpVnOXLub9jsjULfGx0yjH4w+KsZSZCXatoCHbj/RJtkzuZY6V9to/hkH3S
			InQB4a3UAgMCAwEAAaNgMF4wEQYJYIZIAYb4QgEBBAQDAgZAMA4GA1UdDwEB/wQEAwIE8DAfBgNV
			HSMEGDAWgBQ7oCE35Uwn7FsjS01w5e3DA1CrrjAYBgNVHREEETAPgQ1tYWxsYUBzdW4uY29tMA0G
			CSqGSIb3DQEBBAUAA0EAMlbfBg/ff0Xkv4DOR5LEqmfTZKqgdlD81cXynfzlF7XfnOqI6hPIA90I
			x5Ql0ejivIJAYcMGUyA+/YwJg2FGoA==
								</X509Certificate>
							</X509Data>
						</KeyInfo>
						<EncryptionMethod Algorithm=
							"https://www.w3.org/2001/04/xmlenc#aes128-cbc">
							<KeySize xmlns="https://www.w3.org/2001/04/xmlenc#">128</KeySize>
						</EncryptionMethod>
					</KeyDescriptor>
				</IDPSSODescriptor>
			</EntityDescriptor>`)),
			}, nil),
			expected: `
			# HELP saml_x509_cert_not_valid_after SAML X509 certificate expiration date
			# TYPE saml_x509_cert_not_valid_after gauge
			saml_x509_cert_not_valid_after{entityid="example-entity-id",issuer_C="US",issuer_CN="Certificate Manager",issuer_L="Santa Clara",issuer_O="Sun Microsystems Inc.",issuer_ST="",serial_number="1679",subject_C="",subject_CN="loadbalancer-9.siroe.com",subject_L="",subject_O="",url="//saml-metadata",use="encryption"} 1.280879777e+09

			# HELP saml_x509_cert_not_valid_before SAML X509 certificate validity period start
			# TYPE saml_x509_cert_not_valid_before gauge
			saml_x509_cert_not_valid_before{entityid="example-entity-id",issuer_C="US",issuer_CN="Certificate Manager",issuer_L="Santa Clara",issuer_O="Sun Microsystems Inc.",issuer_ST="",serial_number="1679",subject_C="",subject_CN="loadbalancer-9.siroe.com",subject_L="",subject_O="",url="//saml-metadata",use="encryption"} 1.162943777e+09

			# HELP saml_x509_read_errors_total Errors encountered while parsing SAML X509 certificates
			# TYPE saml_x509_read_errors_total counter
			saml_x509_read_errors_total{entityid="example-entity-id",url="//saml-metadata",use="signing"} 1
			`,
		},
		{
			name: "Identity provider metadata with an invalid signing certificate",
			response: transport.NewMock(&http.Response{
				StatusCode: 200,
				Body: io.NopCloser(strings.NewReader(`
				<EntityDescriptor
				xmlns="urn:oasis:names:tc:SAML:2.0:metadata"
				entityID="example-entity-id">
				<IDPSSODescriptor
					AuthnRequestsSigned="false"
					WantAssertionsSigned="false"
					protocolSupportEnumeration=
						"urn:oasis:names:tc:SAML:2.0:protocol">
					<KeyDescriptor use="signing">
						<KeyInfo xmlns="http://www.w3.org/2000/09/xmldsig#">
							<X509Data>
								<X509Certificate>This is an invalid cert</X509Certificate>
							</X509Data>
						</KeyInfo>
					</KeyDescriptor>
				</IDPSSODescriptor>
			</EntityDescriptor>`)),
			}, nil),
			expected: `
			# HELP saml_x509_read_errors_total Errors encountered while parsing SAML X509 certificates
			# TYPE saml_x509_read_errors_total counter
			saml_x509_read_errors_total{entityid="example-entity-id",url="//saml-metadata",use="signing"} 1
			`,
		},
		{
			name: "Invalid saml metadata",
			response: transport.NewMock(&http.Response{
				StatusCode: 500,
				Body:       io.NopCloser(strings.NewReader(`Invalid XML`)),
			}, nil),
			expected: `
			# HELP saml_metadata_errors_total Errors encountered while parsing SAML metadata
			# TYPE saml_metadata_errors_total counter
			saml_metadata_errors_total{url="//saml-metadata"} 1
			`,
		},
		{
			name: "Failed to read saml metadata response body",
			response: transport.NewMock(&http.Response{
				StatusCode: 500,
				Body:       io.NopCloser(iotest.ErrReader(errors.New("random error"))),
			}, nil),
			expected: `
			# HELP saml_metadata_errors_total Errors encountered while parsing SAML metadata
			# TYPE saml_metadata_errors_total counter
			saml_metadata_errors_total{url="//saml-metadata"} 1
			`,
		},
		{
			name: "Failed to execute http request",
			response: transport.NewMock(nil,
				errors.New("random error"),
			),
			expected: ``,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			certParse.Reset()
			metadataParse.Reset()

			reg := prometheus.NewRegistry()

			mockHttpClient := http.DefaultClient
			mockHttpClient.Transport = test.response

			emptyDependantCollector := New(logr.Discard(), mockHttpClient, []*url.URL{})
			c := New(logr.Discard(), mockHttpClient, []*url.URL{{Host: "saml-metadata"}}, emptyDependantCollector)
			reg.MustRegister(c)

			assert.NoError(t, testutil.GatherAndCompare(reg, strings.NewReader(test.expected)))
		})
	}
}
