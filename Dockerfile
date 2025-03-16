FROM gcr.io/distroless/static:nonroot@sha256:b35229a3a6398fe8f86138c74c611e386f128c20378354fc5442811700d5600d
WORKDIR /
COPY saml-exporter saml-exporter
EXPOSE      9412

ENTRYPOINT ["/saml-exporter"]
