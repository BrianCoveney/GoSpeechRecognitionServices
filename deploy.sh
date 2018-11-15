#!/usr/bin/env bash

kubectl delete service frontend
kubectl delete service mongodb-repository
kubectl delete deployment frontend
kubectl delete deployment mongodb-repository

kubectl create -f mongodb-repository-service.yaml,mongodb-repository-deployment.yaml,frontend-service.yaml,frontend-deployment.yaml
