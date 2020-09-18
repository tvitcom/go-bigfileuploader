#!/bin/sh

set -a
test -f ./.env && . ./.env
set +a

#go run main.go
bigfileuploader >> server.log&
