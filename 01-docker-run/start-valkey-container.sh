podman run \
    -p 6379:7000 \
    -v valkey-data:/data \
    --network bridge \
    valkey/valkey:latest \
    valkey-server --port 7000