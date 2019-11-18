# SpeechRecognition Golang Frontend

``` git clone https://github.com/BrianCoveney/GoSpeechRecognitionServices.git  ```

``` docker-compose up ```

Edit ```/etc/hosts``` file, and add the following:
``` 127.0.0.1   speech.local ```

Accessible locally:
http<nolink>://speech.local/

Accessible online:
[https://speech.briancoveney.com/](https://speech.briancoveney.com/)

## Making changes in code / web template, and updating image locally:

Remove image:
``` docker rmi -f bricov/speech_frontend ```

Make your changes and push to DockerHub:
``` docker push bricov/speech_frontend ```

Push code change to GitHub.

## Updating server:

Pull code change from GitHub.

Stop all docker containers:  
``` docker stop $(docker ps -q) ```

Remove Docker Images, Containers, and Volumes:  
``` docker rmi -f <image_name> / docker rm <container_name> / docker system prune / docker volume prune  ```

Start the containers in the background and leave them running:
``` docker-compose up -d ```

## Running Kubernetes Locally via Minikube

Set up Minikube on local machine by following the docs:  
[https://kubernetes.io/docs/setup/minikube/](https://kubernetes.io/docs/setup/minikube/)

Start running Minikube (delete and start if it hangs):  
```  minikube start -p brian-cluster --vm-driver=virtualbox  ```

Use the Kompose conversion tool which converts our ``` docker-compose.yaml ``` to Kubernetes ready resources:  
[https://github.com/kubernetes/kompose](https://github.com/kubernetes/kompose)

Run the following script which creates our deployment using ```kubectl create -f``` with our ```.yaml``` files:  
``` $./deploy.sh ```

Find the minikube ip for the cluster:  
``` minikube -p brian-cluster ip ```

The Minikube VM is exposed to the host system via a host-only IP address. Our ``` frontend-service.yaml ``` includes a ``` NodePort ``` set as `30008`, which provides access to the website locally.   

![alt text](https://github.com/BrianCoveney/SpeechRecognition-Golang-Frontend/blob/master/images/terminal_1.png)

We can then use [Postman](https://www.getpostman.com/) to test our API endpoints. Here we create a child using POST:     
![alt text](https://github.com/BrianCoveney/SpeechRecognition-Golang-Frontend/blob/master/images/postman_1.png)

Update image:    
``` kubectl set image deployment frontend frontend=bricov/speech_frontend:latest ```

Open in browser:  
``` http://192.168.99.100:30008 ```

Updating image after a change:    
``` kubectl set image deployment frontend frontend=bricov/speech_frontend:latest ```

## Running Kubernetes on Google Cloud Platform:  


Set a default project:  
``` gcloud config set project [PROJECT_ID] ```

Set a default compute zone:  
``` gcloud config set compute/zone [COMPUTE_ZONE] ```

Creat a GKE cluster:  
``` gcloud container clusters create [CLUSTER_NAME] ``` 

Authentication credentials for the cluster:  
``` gcloud container clusters get-credentials [CLUSTER_NAME] ``` 

Create the Deployment:  
``` kubectl create -f mongodb-repository-service.yaml,mongodb-repository-deployment.yaml,frontend-service.yaml,frontend-deployment.yaml ```

Exposing the Deployment to the internet. Note: Use NodePort for Minikube(local). Use LoadBalancer for GCP:  
``` kubectl expose deployment frontend --type=LoadBalancer --name=frontend-service ```

Inspecting the application:    
``` kubectl get service ```
![alt text](https://github.com/BrianCoveney/SpeechRecognition-Golang-Frontend/blob/master/images/terminal_2.png)

View the application:   
``` http://[EXTERNAL_IP]/ ```

We can then use [Postman](https://www.getpostman.com/) to test our API endpoints. Here we Update a child using PUT:     
![alt text](https://github.com/BrianCoveney/SpeechRecognition-Golang-Frontend/blob/master/images/postman_2.png)

This is a Functional Test of the API. We make a GET request and validate the resonse by writing JS tests in Postman:     
![alt text](https://github.com/BrianCoveney/SpeechRecognition-Golang-Frontend/blob/master/images/postman_3.png)

## Cleaning Up

Login:  
```gcloud compute os-login describe-profile```

Updates a kubeconfig file with appropriate credentials:  
``` gcloud container clusters get-credentials [CLUSTER_NAME]```

Delete services and clusters:   
``` kubectl delete service [SERVICE_NAME] ```   
``` gcloud container clusters delete [CLUSTER_NAME] ```
