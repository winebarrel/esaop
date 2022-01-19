FROM debian:bullseye-slim

RUN apt-get update && \
  apt-get install -y \
  gettext-base \
  ca-certificates \
  curl

ARG OPEN_ESA_VERSION=0.1.0
RUN curl -sSfL https://github.com/winebarrel/openesa/releases/download/v${OPEN_ESA_VERSION}/openesa_${OPEN_ESA_VERSION}_linux_amd64.tar.gz \
  | tar zxf -

COPY dockerfiles /

ENTRYPOINT ["/entrypoint.sh"]
