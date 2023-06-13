#! /bin/sh
kubectl apply -f storageclass.yml
kubectl apply -f persistent-volume.yml
kubectl apply -f persistent-volume-claim.yml
kubectl apply -f secret.yml
kubectl apply -f configmap.yml
kubectl apply -f statefulsets.yml
kubectl apply -f service.yml