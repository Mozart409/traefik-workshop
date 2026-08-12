# 04 Traefik Advanced

You have reached the final step of this tutorial. Congrats.

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
- Valkey
- Redis Commander
- Golang Webserver
- NextJS App to interact with Valkey
- Whoami Container with rate-limiting
- Jaeger Tracing

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
- [Redis Commander](http://rediscommander.docker.localhost)
- [Whoami](http://whoami.docker.localhost)
- [Fibonacci](http://fibo.docker.localhost)
- [NextJS](http://next.docker.localhost)
- [Jaeger](http://jaeger.docker.localhost)

## FAQ

#### How do I inspect the redis image with redis commander

Open the [Redis Commander](http://rediscommander.docker.localhost) website. It's pre-configured via the `REDIS_HOST`/`REDIS_PORT` environment variables to connect straight to `valkey:6379`, so the keyspace is visible immediately, no setup required.

#### How do I hit the rate-limit

Open the [Whoami](http://whoami.docker.localhost) website and press reload as fast as you can. You should start seeing `429 Too Many Requests`.

#### Can I change the rate limit of the whoami container

Yes. Just change the limit from

```yaml
traefik.http.middlewares.http-ratelimit.ratelimit.average=1
traefik.http.middlewares.http-ratelimit.ratelimit.burst=2
```

to

```yaml
traefik.http.middlewares.http-ratelimit.ratelimit.average=5
traefik.http.middlewares.http-ratelimit.ratelimit.burst=20
```

#### How do I view traces from traefik

Hit one of the other sites first so there is something to trace, then open
[Jaeger](http://jaeger.docker.localhost) and search for the `traefik` service.

#### Why Jaeger and not Zipkin

Traefik v3 removed the per-vendor tracing backends (zipkin, jaeger, datadog,
instana, haystack) and speaks OpenTelemetry only. Leaving a `tracing.zipkin`
block in `traefik.yml` makes Traefik refuse to start:

```
command traefik error: field not found, node: spanNameLimit
```

Jaeger all-in-one ingests OTLP natively on `:4318`, so Traefik exports straight
to it with no collector in between:

```yaml
# traefik.yml
tracing:
  serviceName: traefik
  otlp:
    http:
      endpoint: http://jaeger:4318/v1/traces
```
