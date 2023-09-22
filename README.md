# Prometheus SAML Metadata exporter
[![release](https://github.com/doodlescheduling/saml-exporter/actions/workflows/release.yaml/badge.svg)](https://github.com/doodlescheduling/saml-exporter/actions/workflows/release.yaml)
[![Go Report Card](https://goreportcard.com/badge/github.com/doodlescheduling/saml-exporter)](https://goreportcard.com/report/github.com/doodlescheduling/saml-exporter)
[![OpenSSF Scorecard](https://api.securityscorecards.dev/projects/github.com/DoodleScheduling/saml-exporter/badge)](https://api.securityscorecards.dev/projects/github.com/DoodleScheduling/saml-exporter)
[![Coverage Status](https://coveralls.io/repos/github/DoodleScheduling/saml-exporter/badge.svg?branch=master)](https://coveralls.io/github/DoodleScheduling/saml-exporter?branch=master)
[![Artifact Hub](https://img.shields.io/endpoint?url=https://artifacthub.io/badge/repository/saml-exporter)](https://artifacthub.io/packages/search?repo=saml-exporter)

SAML Metadata exporter for [Prometheus](https://prometheus.io).

## Features

* Tests if the SAML endpoint is reachable and exposes related http metrics
* Exposes metrics related to all encryption and signing x509 certificates
* Supports multiple SAML endpoints

## Installation

Get the exporter either as a binaray from the latest release or packaged as a [Docker image](https://github.com/doodlescheduling/saml-exporter/pkgs/container/saml-exporter).

### Helm Chart
For kubernetes users there is an official helm chart.
Please read the installation instructions [here](https://github.com/doodlescheduling/saml-exporter/blob/master/chart/saml-exporter/README.md).

```sh
helm template saml-exporter oci://ghcr.io/doodlescheduling/charts/saml-exporter --set samlMetadataURLSlice='{http://idp/metadata}'
```

### Docker
You can run the exporter using docker:
```sh
docker run ghcr.io/doodlescheduling/saml-exporter:latest http://idp/metadata
```

## Usage

```
saml-exporter
```

Use the `-help` flag to get help information.

## Access metrics
The metrics are by default exposed at `/metrics`.

```
curl localhost:9412/metrics
```

## Exporter configuration

The exporter can be configured by either command line flags (`saml-exporter -h`) or by defining env variables.

| Env variable             | Description                              | Default |
|--------------------------|------------------------------------------|---------|
| URL                      | Comma separated list of http URL to SAML metadata  | `` |
| LOG_LEVEL                | Log level                                | `info` |
| LOG_ENCODING             | Log format                               | `json` |
| BIND                     | Bind address for the HTTP server         | `:9412` |
| METRICS_PATH             | Metrics endpoint                         | `/metrics` |
| HEALTH_PATH              | Health probe endpoint                    | `/health` |
| USER_AGENT               | HTTP request User agent                  | `saml-exporter (go-http-client)` |

## Exposed metrics 

| Name                     | Description                              | Type | Labels |
|--------------------------|------------------------------------------|---------|-----------|
| `saml_exporter_build_info`    | Build info SAML exporter            | `Gauge` | `"branch", "goversion", "revision", "revision"` |
| `http_client_request`    | HTTP client request                      | `Counter` | `"host", "code", "method"` |
| `saml_metadata_errors`   | Errors encountered while parsing SAML metadata | `Counter` | `"url"` |
| `saml_x509_read_errors`  | Errors encountered while parsing SAML X509 certificates  | `Counter` | `"entityid", "use"` |
| `saml_x509_cert_not_after` | SAML X509 certificate expiration date  | `Gauge` | `"entityid", "use", "serial_number", "issuer_C", "issuer_CN", "issuer_L", "issuer_O", "issuer_ST", "subject_C", "subject_CN", "subject_L", "subject_O"` |
| `saml_x509_cert_not_before` | SAML X509 certificate not valid before  | `Gauge` | `"entityid", "use", "serial_number", "issuer_C", "issuer_CN", "issuer_L", "issuer_O", "issuer_ST", "subject_C", "subject_CN", "subject_L", "subject_O"` |

## Grafana dashboard

This exporter comes with a read to use grafana dashboard, see ./grafana/dashboard.json
**Note**: The helm chart as well as the kustomize base will deploy the grafana dashboard as a ConfigMap.
