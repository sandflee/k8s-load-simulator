#!/bin/sh

glide update
if [ $? -ne 0 ];then
    echo "glide upadate failed"
    exit -1
fi
glide install --strip-vendor --strip-vcs
