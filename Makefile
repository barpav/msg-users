up:
	sudo docker-compose up -d
	sudo docker ps
down:
	sudo docker-compose down
	sudo docker ps

run-s:
	sudo docker run -p5432:5432 --name msg-storage-users-v1 -e POSTGRES_PASSWORD=postgres -d barpav/msg-storage-users:v1
	sudo docker ps
stop-s:
	sudo docker stop msg-storage-users-v1
	sudo docker rm msg-storage-users-v1
	sudo docker ps

build-u:
	sudo docker image rm -f barpav/msg-users:v1
	sudo docker build -t barpav/msg-users:v1 -f docker/service/Dockerfile .
	sudo docker image ls
build-s:
	sudo docker image rm -f barpav/msg-storage-users:v1
	sudo docker build -t barpav/msg-storage-users:v1 -f docker/storage/Dockerfile ./docker/storage
	sudo docker image ls

clear-u:
	sudo docker image rm -f barpav/msg-users:v1
	sudo docker image ls
clear-s:
	sudo docker image rm -f barpav/msg-storage-users:v1
	sudo docker image ls
clear:
	sudo docker image rm -f barpav/msg-users:v1 barpav/msg-storage-users:v1
	sudo docker image ls

exec-u:
	sudo docker exec -it msg-users-v1 sh
exec-s:
	sudo docker exec -it -u postgres msg-storage-users-v1 bash

push-u:
	sudo docker push barpav/msg-users:v1
push-s:
	sudo docker push barpav/msg-storage-users:v1
push:
	sudo docker push barpav/msg-users:v1
	sudo docker push barpav/msg-storage-users:v1