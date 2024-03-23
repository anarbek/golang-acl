# golang-acl
jwt acl example in golang

## build docker image
Docker build -t nr-golang-acl:1.3 .
#cd G:\CustomProjects\GolangRelated\golang-acl\src
#Docker build -t nr-golang-acl:1.3 .

## docker compose
cd /path/to/golang-acl/inc/DockerRelated/DockerComposeLocal
docker-compose -f docker-compose-local.yaml up

## minikube windows
cd G:\CustomProjects\GolangRelated\golang-acl\inc\DockerRelated\MinikubeRelated
minikube -p minikube docker-env | Invoke-Expression
kubectl create deployment golang-acl-web-interface --image=nr-golang-acl:1.3
kubectl expose deployment golang-acl-web-interface --type=NodePort --port=8081
minikube service golang-acl-web-interface
#here app is sending response to 8081 port of pod
#and we access app using 8082 port from browser
kubectl port-forward service/golang-acl-web-interface 8082:8081

## kubectl commands
cd G:\CustomProjects\GolangRelated\golang-acl\inc\DockerRelated\MinikubeRelated
kubectl get pods
kubectl describe pod [nameof-pod]
build with yaml: kubectl apply -f golang-acl.yaml
kubectl expose deployment golang-acl-web-interface --type=NodePort --port=8081
minikube service golang-acl-web-interface

#using yaml in kubectl:
bring up with nodeport: kubectl apply -f golang-acl-nodeport.yaml
bring up with load-balancer: kubectl apply -f golang-acl-load-balancer.yaml
minikube service golang-acl-web-balancer

#curl example:
curl -L -H "Authorization: Bearer <token>" http://127.0.0.1:56990/rndAuth
