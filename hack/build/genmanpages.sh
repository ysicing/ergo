#!/bin/sh
set -e
rm -rf manpages
mkdir manpages
go run ./ergo.go man | gzip -c -9 >manpages/ergo.1.gz
