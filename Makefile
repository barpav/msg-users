up:
	sudo docker-compose up -d --wait
	sudo docker ps
down:
	sudo docker-compose down
	sudo docker ps

run-s:
	sudo docker run -p5432:5432 --name msg-storage-users-v1 -e POSTGRES_PASSWORD=postgres -d ghcr.io/barpav/msg-storage-users:v1
	sudo docker ps
stop-s:
	sudo docker stop msg-storage-users-v1
	sudo docker rm msg-storage-users-v1
	sudo docker ps

up-debug:
	sudo docker-compose -f compose-debug.yaml up -d --wait
down-debug:
	sudo docker-compose -f compose-debug.yaml down

new-user:
	curl -v -X POST	-H "Content-Type: application/vnd.newUser.v1+json" \
	-d '{"id": "jane", "name": "Jane Doe", "password": "My1stGoodPassword"}' \
	localhost:8080
new-session:
	curl -v -X POST -H "Authorization: Basic amFuZTpNeTFzdEdvb2RQYXNzd29yZA==" localhost:8081
# make get-info KEY=session-key 
get-info:
	curl -v -H "Authorization: Bearer $(KEY)" \
	-H "Accept: application/vnd.userInfo.v1+json" \
	localhost:8080
# make edit-info KEY=session-key NAME="New name"
edit-info:
	curl -v -X PATCH -H "Authorization: Bearer $(KEY)" \
	-H "Content-Type: application/vnd.userProfileCommon.v1+json" \
	-d '{"name": "$(NAME)"}' \
	localhost:8080

build-u:
	sudo docker image rm -f ghcr.io/barpav/msg-users:v1
	sudo docker build -t ghcr.io/barpav/msg-users:v1 -f docker/service/Dockerfile .
	sudo docker image ls
build-s:
	sudo docker image rm -f ghcr.io/barpav/msg-storage-users:v1
	sudo docker build -t ghcr.io/barpav/msg-storage-users:v1 -f docker/storage/Dockerfile ./docker/storage
	sudo docker image ls

clear-u:
	sudo docker image rm -f ghcr.io/barpav/msg-users:v1
	sudo docker image ls
clear-s:
	sudo docker image rm -f ghcr.io/barpav/msg-storage-users:v1
	sudo docker image ls
clear:
	sudo docker image rm -f ghcr.io/barpav/msg-users:v1 ghcr.io/barpav/msg-storage-users:v1
	sudo docker image ls

exec-u:
	sudo docker exec -it msg-users-v1 sh
exec-s:
	sudo docker exec -it -u postgres msg-storage-users-v1 bash

push-u:
	sudo docker push ghcr.io/barpav/msg-users:v1
push-s:
	sudo docker push ghcr.io/barpav/msg-storage-users:v1
push:
	sudo docker push ghcr.io/barpav/msg-users:v1
	sudo docker push ghcr.io/barpav/msg-storage-users:v1

proto:
	protoc --go_out=. --go_opt=paths=source_relative \
    --go-grpc_out=. --go-grpc_opt=paths=source_relative \
    users_service_go_grpc/users_service.proto