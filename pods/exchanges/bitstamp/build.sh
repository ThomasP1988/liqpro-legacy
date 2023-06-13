#! /bin/sh
eval $(minikube -p minikube docker-env)
go mod vendor
docker build . -t exchanges/bitstamp
rm -rf ./vendor