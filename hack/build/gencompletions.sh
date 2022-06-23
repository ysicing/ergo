#!/bin/sh
set -e
rm -rf completions
mkdir completions
for sh in bash zsh fish; do
	go run ./ergo.go completion "$sh" >"completions/ergo.$sh"
done
