# 03 Traefik

You have reached the traefik part.

## Warning

> Traefik discovers your containers by talking to the container API. It speaks
> the docker-compatible API, so the podman socket works as-is — there is no
> docker daemon and no socket-proxy sidecar involved. Depending on your setup
> you may need to point the mount at a different socket.

```yaml
# docker-compose.yml
services:
  traefik:
    security_opt:
      - label=disable
    volumes:
      # Rootful podman machine (the default on macOS / Podman Desktop)
      - /run/podman/podman.sock:/var/run/docker.sock:ro
      # Rootless podman instead? Use your own UID
      #- /run/user/1000/podman/podman.sock:/var/run/docker.sock:ro
```

The podman machine runs SELinux in enforcing mode, which otherwise stops Traefik
from reading that socket. `label=disable` turns the label check off for the
traefik container only, which is nicer than relabelling the host's socket.

## Features

- Traefik
- Whoami Container Loadbalanced

## Start the project

See all logs from podman compose

```bash
just dev
```

Hide all logs from podman compose

```bash
just start
```

## Websites

- [Traefik](http://localhost:8080/dashboard/) — the trailing slash is required
- [Whoami](http://whoami.docker.localhost) — reload a few times to see the
  hostname change

## FAQ

#### Why are there no ports exposed in the whoami container

Because Traefik takes care of that. Never define ports again. Scale to your
desired amount and traefik automagically loadbalances them.

#### How do I check what Traefik discovered

```bash
curl -s http://localhost:8080/api/http/services
```

The `whoami` service should list one server per replica, all `UP`.
