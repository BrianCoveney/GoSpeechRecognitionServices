# SpeechRecognition Golang Frontend

``` git clone https://github.com/BrianCoveney/GoSpeechRecognitionServices.git  ```

``` docker-compose up ```

Edit ```/etc/hosts``` file, and add the following:
``` 127.0.0.1   speech.local ```

Accessible locally:
http<nolink>://speech.local/

Accessible online:
[http://speech.briancoveney.com/](http://speech.briancoveney.com/)


## Running Kubernetes Locally via Minikube

Run the following command:  
``` $./deploy.sh ```

Find the frontend service with the exposed port of 30008:  
``` $minikube service list ```

e.g:  

| default | frontend | http<span></span>://192.168.99.100:30008 |
  

Open in browser:  
``` http://192.168.99.100:30008 ```

## Making changes and updating image

Find local docker image:  
``` docker images ```

Remove image:  
``` docker rmi -f bricov/speech_frontend ```

Make your changes and push to DockerHub:  
``` docker push bricov/speech_frontend ```

Update image:    
``` kubectl set image deployment frontend frontend=bricov/speech_frontend:latest ```
