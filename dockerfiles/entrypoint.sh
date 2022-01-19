#!/bin/bash
set -e

envsubst '
  $ESAOP_ADDR
  $ESAOP_PORT
  $ESAOP_TAEM
  $ESAOP_SESSION_SECRET
  $ESAOP_COOKIE_SECURE
  $ESAOP_OAUTH2_CLIENT_ID
  $ESAOP_OAUTH2_CLIENT_SECRET
  $ESAOP_OAUTH2_REDIRECT_HOST
' < /esaop.toml.template > /esaop.toml

exec /esaop
