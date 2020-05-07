#!/bin/bash

ergo init vm --vmname istio --vmnum 3

ergo init debian  --ip 11.11.11.111 --ip 11.11.11.112 --ip 11.11.11.113

ergo install docker --ip 11.11.11.111 --ip 11.11.11.112 --ip 11.11.11.113

ergo install tools --ip 11.11.11.111 --ip 11.11.11.112 --ip 11.11.11.113

# 安装master单节点
ergo install k8s --enablenfs=true --mip 11.11.11.111
# 安装多节点
ergo install k8s --enablenfs=true --mip 11.11.11.111 --wip 11.11.11.112-11.11.11.113
