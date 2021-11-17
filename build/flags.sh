#!/bin/sh
PKGCONFIG="github.com/hirokimoto/crypto-auto/config"

LD_FLAG_MESSAGE="-X '${PKGCONFIG}.ApplicationVersion=${VERSION}'"

LDFLAGS="${LD_FLAG_MESSAGE}"