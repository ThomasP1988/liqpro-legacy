#! /bin/sh
export IDENTITY_SCHEMAS=$PWD/identity.traits.schema.json
helm install auth-service -f config.yml ory/kratos 
export POD_NAME=$(kubectl get pods --namespace default -l "app.kubernetes.io/name=kratos,app.kubernetes.io/instance=auth-service" -o jsonpath="{.items[0].metadata.name}")
echo "Visit http://127.0.0.1:80 to use your application"
kubectl port-forward $POD_NAME 4849:4433
export KRATOS_PUBLIC_URL=http://127.0.0.1:4849/
export POD_NAME=$(kubectl get pods --namespace default -l "app.kubernetes.io/name=kratos,app.kubernetes.io/instance=auth-service" -o jsonpath="{.items[0].metadata.name}")
echo "Visit http://127.0.0.1:80 to use your application"
kubectl port-forward $POD_NAME 4849:4434
export KRATOS_ADMIN_URL=http://127.0.0.1:4849/