#!/usr/bin/env bash

version=$(cat version.txt)

cp -a dist/ergo_linux_amd64 hack/deb/opsergo/usr/local/bin/ergo

chmod +x hack/deb/opsergo/usr/local/bin/ergo

pushd hack/deb/opsergo
    [ -f "DEBIAN/md5sums" ] && rm -rf DEBIAN/md5sums
    [ -f "usr/local/bin/.gitkeep" ] && rm -rf usr/local/bin/.gitkeep
    md5sum usr/local/bin/ergo >> DEBIAN/md5sums
popd

docker run --rm -v ${PWD}/hack/deb:/builder ysicing/ergodeb opsergo ${version}

if [ -z "$DEBTOKEN" ]; then
    exit 1
fi

opsergodeb="hack/deb/opsergo.${version}_amd64.deb"

[ -f "${opsergodeb}" ] && (
  curl -F package=@${opsergodeb} https://${DEBTOKEN}@push.fury.io/ysicing/
)