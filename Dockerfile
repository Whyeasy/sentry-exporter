FROM alpine

COPY sentry-exporter /usr/bin/
ENTRYPOINT ["/usr/bin/sentry-exporter"]