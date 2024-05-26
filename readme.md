To connect to Evans
evans --host localhost --port 8080 -r repl

To list all the service running on the secfic port
lsof -i :8080

To kill the PID
kill -9 pid

For showing logs inside any service
docker logs deploy-authentication-service-1

To list all in a folder
ls -l

To go inside the docker container
docker exec -it postgres16 psql -U root

List of databases
\l

To switch db inside psql
\c dbname


To list all table inside db
\dt or \d

To show details of a table use 
\d table_name

.


find kubernetes/ -type f -name "*.yaml" -exec kubectl apply -f {} \;

eval $(minikube docker-env)
eval $(minikube -p minikube docker-env)
eval $(minikube docker-env -u) - unset
eval $(minikube -p minikube docker-env -u)