FROM debian:bullseye-slim

RUN apt-get update && \
  apt-get install -y \
  gettext-base \
  ca-certificates \
  curl

ARG OPENESA_VERSION=0.1.1
RUN curl -sSfL https://github.com/winebarrel/openesa/releases/download/v${OPENESA_VERSION}/openesa_${OPENESA_VERSION}_linux_amd64.tar.gz \
  | tar zxf -

COPY dockerfiles /

ENTRYPOINT ["/entrypoint.sh"]
