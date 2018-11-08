#!/bin/sh
set -e
(
   $GOPATH/bin/go-bindata -o ./web.go -pkg web -nomemcopy dist/...
)