#! /bin/sh
# kubectl delete persistentvolume mongodb-standalone
kubectl delete -f secret.yml
kubectl delete -f configmap.yml
kubectl delete -f statefulsets.yml
kubectl delete -f service.yml
kubectl delete -f storageclass.yml
kubectl delete -f persistent-volume-claim.yml
kubectl delete -f persistent-volume.yml