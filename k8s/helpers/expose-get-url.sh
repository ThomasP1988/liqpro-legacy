#! /bin/sh
kubectl expose rc auth-service-super-token --type=NodePort
minikube service auth-service-super-token --url