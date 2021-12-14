#!/bin/bash

version=$(cat version.txt)
macosamd64sha=$(cat dist/ergo_darwin_amd64.sha256sum | awk '{print $1}')
macosarm64sha=$(cat dist/ergo_darwin_arm64.sha256sum | awk '{print $1}')
linuxamd64sha=$(cat dist/ergo_linux_amd64.sha256sum | awk '{print $1}')
linuxarm64sha=$(cat dist/ergo_linux_arm64.sha256sum | awk '{print $1}')


cat > ergo.rb <<EOF
class Ergo < Formula
    desc "Devops tools 运维工具Ergo"
    homepage "https://github.com/ysicing/ergo"
    version "${version}"

    if OS.mac?
      if Hardware::CPU.arm?
        url "https://github.com/ysicing/ergo/releases/download/#{version}/ergo_darwin_arm64"
        sha256 "${macosarm64sha}"
      else
        url "https://github.com/ysicing/ergo/releases/download/#{version}/ergo_darwin_amd64"
        sha256 "${macosamd64sha}"
      end  
    elsif OS.linux?
      if Hardware::CPU.intel?
        url "https://github.com/ysicing/ergo/releases/download/#{version}/ergo_linux_amd64"
        sha256 "${linuxamd64sha}"
      end
      if Hardware::CPU.arm? && Hardware::CPU.is_64_bit?
        url "https://github.com/ysicing/ergo/releases/download/#{version}/ergo_linux_arm64"
        sha256 "${linuxarm64sha}"
      end
    end

    def install
      if OS.mac?
        if Hardware::CPU.intel?
          bin.install "ergo_darwin_amd64" => "ergo"
        else
          bin.install "ergo_darwin_arm64" => "ergo"
        end 
      elsif OS.linux?
        if Hardware::CPU.intel?
          bin.install "ergo_linux_amd64" => "ergo"
        else
          bin.install "ergo_linux_arm64" => "ergo"
        end 
      end

      # Install bash completion
      output = Utils.safe_popen_read(bin/"ergo", "completion", "bash")
      (bash_completion/"ergo").write output

      # Install zsh completion
      output = Utils.safe_popen_read(bin/"ergo", "completion", "zsh")
      (zsh_completion/"_ergo").write output
      
      fish_output = Utils.safe_popen_read(bin/"ergo", "completion", "fish")
      (fish_completion/"ergo.fish").write fish_output
    end

    test do
      assert_match "ergo vervion v#{version}", shell_output("#{bin}/ergo version")
    end
end
EOF

docker build -t ysicing/taprb:ergo -f hack/brew/Dockerfile .
docker push ysicing/taprb:ergo