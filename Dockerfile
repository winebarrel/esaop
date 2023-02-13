FROM --platform=linux/x86_64 debian:bullseye-slim

RUN apt-get update && \
  apt-get install -y \
  gettext-base \
  ca-certificates \
  curl

ARG ESAOP_VERSION=0.2.1
RUN curl -sSfL https://github.com/winebarrel/esaop/releases/download/v${ESAOP_VERSION}/esaop_${ESAOP_VERSION}_linux_amd64.tar.gz \
  | tar zxf -

COPY dockerfiles /

ENTRYPOINT ["/entrypoint.sh"]
