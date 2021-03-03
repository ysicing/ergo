#!/bin/bash

version=$(cat version.txt)
macsha=$(cat dist/ergo_darwin_amd64.sha256sum | awk '{print $1}')
m1sha=$(cat dist/ergo_darwin_arm64.sha256sum | awk '{print $1}')
linuxsha=$(cat dist/ergo_linux_amd64.sha256sum | awk '{print $1}')

cat > ergo.rb <<EOF
class Ergo < Formula
    desc "Devops tools 运维工具"
    homepage "https://github.com/ysicing/ergo"
    version "${version}"
    bottle :unneeded

    if OS.mac?
      if Hardware::CPU.arm?
        url "https://github.com/ysicing/ergo/releases/download/#{version}/ergo_darwin_arm64"
        sha256 "${m1sha}"
      else
        url "https://github.com/ysicing/ergo/releases/download/#{version}/ergo_darwin_amd64"
        sha256 "${macsha}"
      end  
    elsif OS.linux?
      if Hardware::CPU.intel?
        url "https://github.com/ysicing/ergo/releases/download/#{version}/ergo_linux_amd64"
        sha256 "${linuxsha}"
      end
    end

    def install
      if Hardware::CPU.intel?
        bin.install "ergo_darwin_amd64" => "ergo"
      else
        bin.install "ergo_darwin_arm64" => "ergo"
      end 
    end
  end
EOF

docker build -t ysicing/taprb:ergo -f hack/brew/Dockerfile .
docker push ysicing/taprb:ergo