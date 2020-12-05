#!/bin/bash

version=$(cat version.txt)
macsha=$(cat dist/ergo_darwin_amd64.sha256sum | awk '{print $1}')
linuxsha=$(cat dist/ergo_linux_amd64.sha256sum | awk '{print $1}')

dist/ergo_linux_amd64 | grep $version &&

cat > ergo.rb <<EOF
class Ergo < Formula
    desc "Devops tools 运维工具"
    homepage "https://github.com/ysicing/ergo"
    version "${version}"
    bottle :unneeded

    if OS.mac?
      url "https://github.com/ysicing/ergo/releases/download/#{version}/ergo_darwin_amd64"
      sha256 "${macsha}"
    elsif OS.linux?
      if Hardware::CPU.intel?
        url "https://github.com/ysicing/ergo/releases/download/#{version}/ergo_linux_amd64"
        sha256 "${linuxsha}"
      end
    end

    def install
      bin.install "ergo_darwin_amd64" => "ergo"
    end
  end
EOF

docker build -t ysicing/taprb:ergo -f hack/brew/Dockerfile .
docker push ysicing/taprb:ergo