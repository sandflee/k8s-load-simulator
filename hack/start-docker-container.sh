#!/bin/sh


apiserver="http://10.234.136.224:8081"
log_dir="/data/home/gaia/moon/simluator/log"

docker run -d --net=host -v $log_dir:/k8s/log docker.oa.com:8080/g_moonfang/k8s-load-simulator:v1 k8s-load-simulator --apiserver=$apiserver --nodeNum=1000 --v=3 --log_dir=/k8s/log

