FROM gcr.io/distroless/static:nonroot@sha256:42d15c647a762d3ce3a67eab394220f5268915d6ddba9006871e16e4698c3a24
WORKDIR /
COPY saml-exporter saml-exporter
EXPOSE      9412

ENTRYPOINT ["/saml-exporter"]
