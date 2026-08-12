docker run \
    -e ALLOW_EMPTY_PASSWORD=yes \
    -e VALKEY_PORT_NUMBER=7000 -p 6379:7000 \
    -v valkey-data:/bitnami/valkey/data \
    --network bridge \
    bitnami/valkey:latest