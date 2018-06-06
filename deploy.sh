#!/usr/bin/env bash

kubectl delete service frontend
kubectl delete service mongodb

kubectl delete deployment frontend
kubectl delete deployment mongodb

#kubectl create -f mongodb-service.yaml,mongodb-deployment.yaml,frontend-service.yaml,frontend-deployment.yaml
