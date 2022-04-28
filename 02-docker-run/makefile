dev:
	docker run \
    -e ALLOW_EMPTY_PASSWORD=yes \
    -e REDIS_PORT_NUMBER=7000 -p 6379:7000 \
    -v redis-data:/bitnami/redis/data \
    --network bridge \
    bitnami/redis:latest

start:
	docker run \
    -e ALLOW_EMPTY_PASSWORD=yes \
    -e REDIS_PORT_NUMBER=7000 -p 6379:7000 \
    -v redis-data:/bitnami/redis/data \
    --network bridge \
    bitnami/redis:latest