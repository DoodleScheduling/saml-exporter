# SAML exporter

Installs the [SAML metadata exporter](https://github.com/doodlescheduling/saml-exporter) for [Prometheus](https://prometheus.io/).

## Installing the Chart

To install the chart with the release name `saml-exporter`:

```console
helm upgrade saml-exporter --install oci://ghcr.io/doodlescheduling/charts/saml-exporter
```

This command deploys the MongoDB Exporter with the default configuration. The [configuration](#configuration) section lists the parameters that can be configured during installation.

## Using the Chart

The chart comes with a ServiceMonitor for use with the [Prometheus Operator](https://github.com/helm/charts/tree/master/stable/prometheus-operator).
If you're not using the Prometheus Operator, you can disable the ServiceMonitor by setting `serviceMonitor.enabled` to `false` and instead
populate the `podAnnotations` as below:

```yaml
podAnnotations:
  prometheus.io/scrape: "true"
  prometheus.io/port: "metrics"
  prometheus.io/path: "/metrics"
```

## Configuration

To see all configurable options with detailed comments, visit the chart's values.yaml, or run the configuration command:

```sh
helm show values oci://ghcr.io/doodlescheduling/charts/saml-exporter
```
