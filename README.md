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

Run the following command:  
``` $./deploy.sh ```

Find the frontend service with the exposed port of 30008:  
``` $minikube service list ```

e.g:  

| default | frontend | http<span></span>://192.168.99.100:30008 |

Update image:    
``` kubectl set image deployment frontend frontend=bricov/speech_frontend:latest ```

Open in browser:  
``` http://192.168.99.100:30008 ```

Updating image after a change:    
``` kubectl set image deployment frontend frontend=bricov/speech_frontend:latest ```


