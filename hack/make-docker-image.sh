#!/bin/sh

docker build -t docker.oa.com:8080/g_moonfang/k8s-load-simulator:v1 ./
docker push docker.oa.com:8080/g_moonfang/k8s-load-simulator:v1 
