#!/bin/bash

if [ "$1" == "bash" -o "$1" == "ash" ];then
    exec /bin/bash
fi

exec $@
