FROM gcr.io/distroless/static:nonroot@sha256:92d40eea0b5307a94f2ebee3e94095e704015fb41e35fc1fcbd1d151cc282222
WORKDIR /
COPY saml-exporter saml-exporter
EXPOSE      9412

ENTRYPOINT ["/saml-exporter"]
