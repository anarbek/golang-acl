# golang-acl
jwt acl example in golang

## build docker image
Docker build -t nr-golang-acl:1.0 .
## docker compose
cd /path/to/golang-acl/inc/DockerRelated/DockerComposeLocal
docker-compose -f docker-compose-local.yaml up

## minikube windows
minikube -p minikube docker-env | Invoke-Expression
kubectl create deployment hello-minikube5 --image=nr-golang-acl:1.0
kubectl expose deployment hello-minikube5 --type=NodePort --port=8081
minikube service hello-minikube5
#here app is sending response to 8081 port of pod
#and we access app using 8082 port from browser
kubectl port-forward service/hello-minikube5 8082:8081

## kubectl commands
kubectl get pods
kubectl describe pod [nameof-pod]