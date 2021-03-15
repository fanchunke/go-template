CWD = $(shell pwd)
IMAGE = go-dev
CONTAINER = go-dev

# 本机开发测试
dev: stop
	docker build -f Dockerfile -t ${IMAGE} .
	docker run --detach --publish=8773:8000 \
		--volume=${CWD}/logs:/home/works/program/logs \
		--restart=always --memory=1GB --name=${CONTAINER} \
		${IMAGE}

# 停止本机容器并删除
stop:
	docker ps -aq --filter name=${CONTAINER} | xargs docker stop; true
	docker ps -aq --filter name=${CONTAINER} | xargs docker rm; true

# 删除容器和镜像
clean: stop
	docker image rmi ${IMAGE}; true
