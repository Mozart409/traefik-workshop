# 04 Traefik

You have reached the traefik part.

## Warning

> Depending on your os you need to change the docker socket to fit your system.

```yaml
# docker compose-yml
version: '3.7'
services:
  socket-proxy:
    image: tecnativa/docker-socket-proxy
    container_name: socket-proxy
    volumes:
      # Docker Desktop (Windows or Mac OSX)
      - /var/run/docker.sock:/var/run/docker.sock
      # Linux only
      #- /var/run/user/1000/docker.sock:/var/run/docker.sock:ro
```

## Features

- Traefik
- Whoami Container Loadbalanced

## Start the project

See all logs from docker compose

```bash
  make dev
```

Hide all logs from docker compose

```bash
  make start
```

## Websites

- [Traefik](http://localhost:8080/dashboard)
- [Whoami](http://whoami.docker.localhost)

## FAQ

#### Why are the no ports in the whoami contaier exposed

Because Traefik takes care of that. Never define ports again. Scale to your desired amount and traefik automagically loadbalances them.
