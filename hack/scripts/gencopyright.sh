#!/bin/bash

set -xe

addlicense -f licenses/licenses.tpl -ignore web/** -ignore "**/*.md" -ignore vendor/** -ignore "**/*.yml" -ignore "**/*.yaml" -ignore "**/*.sh" ./**
