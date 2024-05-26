postgres:
	@echo Creating a new container for postgres
	docker run --name postgres16 -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d postgres:16-alpine

createdb:
	@echo Creating gestapo db
	docker exec -it postgres16 createdb --username=root --owner=root gestapo

dropdb:
	docker exec -it postgres16 dropdb gestapo

redis:
	@echo Creating a new container for postgres
	docker run --name redis7.2 -p 6379:6379 -d redis:7.2-alpine

authentication_server:
	@echo Running authentication service
	go run cmd/authentication_service/main.go 

admin_server:
	@echo Running admin service
	go run cmd/admin_service/main.go

user_server:
	@echo Running user service
	go run cmd/user_service/main.go	

merchant_server:
	@echo Running merchant service
	go run cmd/merchant_service/main.go

product_server:
	@echo Running product service
	go run cmd/product_service/main.go

order_server:
	@echo Running order service
	go run cmd/order_service/main.go


grpc_gateway:
	@echo Running grpc gateway
	go run cmd/grpc_gateway/main.go

proto:
	@echo deleting generated files if exist..
	rm -f pkg/api/proto/*.go
	@echo Generating all proto pb files..
	protoc -I . \
	--go_out pkg/ --go_opt=paths=source_relative \
	--go-grpc_out pkg/ --go-grpc_opt=paths=source_relative \
	--grpc-gateway_out pkg/ --grpc-gateway_opt=paths=source_relative \
	api/proto/*.proto
	@echo done..

evans:
	@echo Starting evans gRPC client..
	evans --host localhost --port 9002 -r repl      

compose_down: 
	@echo Stopping docker containers
	cd deploy && sudo docker compose down --remove-orphans
	@echo done

compose_up:
	@echo Start docker compose
	cd deploy && sudo docker compose up --build -d 
	@echo done

# prune_images:
# 	@echo prune all images 
# 	cd deploy && sudo docker image prune -a
# 	@echo done 

AUTH_BINARY=authenticationServiceApp
ADMIN_BINARY=adminServiceApp
USER_BINARY=userServiceApp
MERCHANT_BINARY=merchantServiceApp
PRODUCT_BINARY=productServiceApp
ORDER_BINARY=orderServiceApp
GATEWAY_BINARY=gatewayApp

build_authentication:
	@echo Building authentication binary...
	cd cmd/authentication_service && GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o ${AUTH_BINARY} .
	@echo Moving file..
	mv cmd/authentication_service/${AUTH_BINARY} deploy/build
	@echo Done!

build_admin:
	@echo Building admin binary...
	cd cmd/admin_service && GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o ${ADMIN_BINARY} .
	@echo Moving file..
	mv cmd/admin_service/${ADMIN_BINARY} deploy/build
	@echo Done!	

build_user:
	@echo Building user binary...
	cd cmd/user_service && GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o ${USER_BINARY} .
	@echo Moving file..
	mv cmd/user_service/${USER_BINARY} deploy/build
	@echo Done!	

build_merchant:
	@echo Building merchant binary...
	cd cmd/merchant_service && GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o ${MERCHANT_BINARY} .
	@echo Moving file..
	mv cmd/merchant_service/${MERCHANT_BINARY} deploy/build
	@echo Done!		

build_product:
	@echo Building product binary...
	cd cmd/product_service && GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o ${PRODUCT_BINARY} .
	@echo Moving file..
	mv cmd/product_service/${PRODUCT_BINARY} deploy/build
	@echo Done!		

build_order:
	@echo Building order binary...
	cd cmd/order_service && GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o ${ORDER_BINARY} .
	@echo Moving file..
	mv cmd/order_service/${ORDER_BINARY} deploy/build
	@echo Done!				

build_gateway:
	@echo Building gateway binary...
	cd cmd/grpc_gateway && GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o ${GATEWAY_BINARY} .
	@echo Moving file..
	mv cmd/grpc_gateway/${GATEWAY_BINARY} deploy/build
	@echo Done!


run: build_authentication build_admin build_user build_merchant build_product build_order build_gateway 
	@echo Stopping docker images if running...
	cd deploy && docker compose down --remove-orphans
	@echo Building when required and starting docker images...
	cd deploy && sudo docker compose up --build -d
	@echo Docker images built and started!


build_binary: build_authentication build_admin build_user build_merchant build_product build_order build_gateway 

build_image: 
	@echo Building docker images...
	cd deploy && docker build -t deploy-authentication-service -f authentication-service.dockerfile . 
	cd deploy && docker build -t deploy-admin-service -f admin-service.dockerfile . 
	cd deploy && docker build -t deploy-user-service -f user-service.dockerfile . 
	cd deploy && docker build -t deploy-merchant-service -f merchant-service.dockerfile . 
	cd deploy && docker build -t deploy-product-service -f product-service.dockerfile . 
	cd deploy && docker build -t deploy-order-service -f order-service.dockerfile . 
	cd deploy && docker build -t deploy-grpc-gateway -f grpc-gateway.dockerfile . 
	@echo Docker images built!

tag_build_image:
	@echo Tagging services
	docker tag deploy-authentication-service gestapo/authentication:1.0.0
	docker tag deploy-admin-service gestapo/admin:1.0.0
	docker tag deploy-user-service gestapo/user:1.0.0
	docker tag deploy-merchant-service gestapo/merchant:1.0.0
	docker tag deploy-product-service gestapo/product:1.0.0
	docker tag deploy-order-service gestapo/order:1.0.0
	docker tag deploy-grpc-gateway gestapo/gateway:1.0.0
	@echo services tagged successfully 

prune_tags:
	@echo removing tagged images..
	sudo docker rmi gestapo/gateway:1.0.0
	sudo docker rmi gestapo/authentication:1.0.0
	sudo docker rmi gestapo/admin:1.0.0
	sudo docker rmi gestapo/user:1.0.0
	sudo docker rmi gestapo/merchant:1.0.0
	sudo docker rmi gestapo/product:1.0.0
	sudo docker rmi gestapo/order:1.0.0
	@echo done !!

prune:
	@echo removing images..
	sudo docker rmi deploy-authentication-service:latest
	sudo docker rmi deploy-admin-service:latest
	sudo docker rmi deploy-user-service:latest
	sudo docker rmi deploy-merchant-service:latest
	sudo docker rmi deploy-product-service:latest
	sudo docker rmi deploy-order-service:latest
	sudo docker rmi deploy-grpc-gateway:latest
	@echo done !!

prune_minikube_old_images:
	@echo removing images..
	minikube image rm deploy-authentication-service:latest
	minikube image rm deploy-admin-service:latest
	minikube image rm deploy-user-service:latest
	minikube image rm deploy-merchant-service:latest
	minikube image rm deploy-product-service:latest
	minikube image rm deploy-order-service:latest
	minikube image rm deploy-grpc-gateway:latest
	@echo done !!

prune_minikube_images:
	@echo removing tagged images..
	minikube image rm gestapo/gateway:1.0.0
	minikube image rm gestapo/authentication:1.0.0
	minikube image rm gestapo/admin:1.0.0
	minikube image rm gestapo/user:1.0.0
	minikube image rm gestapo/merchant:1.0.0
	minikube image rm gestapo/product:1.0.0
	minikube image rm gestapo/order:1.0.0
	@echo done !!	


# Befor calling minikube call `eval $(minikube -p minikube docker-env)` after call `eval $(minikube -p minikube docker-env -u)`
minikube: build_binary build_image tag_build_image prune_minikube_old_images
	@echo =============Running in minikube=============
	find kubernetes/ -type f -name "*.yaml" -exec kubectl apply -f {} \;
	


.PHONY: postgres createdb dropdb server proto build_authentication run
