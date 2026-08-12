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
- Jaeger Tracing, end to end from Traefik into both apps

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

#### How do I view traces

Hit one of the other sites first so there is something to trace, then open
[Jaeger](http://jaeger.docker.localhost) and search for the `traefik`,
`fiboservice` or `nextjs` service.

Traefik propagates the W3C `traceparent` header to its backends and both apps
read it, so one request is one connected trace instead of three unrelated ones.
A `/fibonacci/:n` trace looks like this:

```
traefik      EntryPoint
traefik      Router
traefik      Service
traefik      ReverseProxy
fiboservice  /fibonacci/:n
fiboservice  Fibonacci       <- the computation itself
```

The NextJS app is instrumented the same way, and because it also loads the
ioredis instrumentation the Valkey command shows up as its own span:

```
traefik  ...
nextjs   GET /api/features
nextjs   executing api route (pages) /api/features
nextjs   hvals            <- the Valkey call
```

#### How are the apps instrumented

`fiboService` builds an OTLP exporter in `main.go` and registers the `otelfiber`
middleware. The NextJS app uses Next's `instrumentation.ts` hook with
`@vercel/otel`. Neither hardcodes the collector address — both read
`OTEL_EXPORTER_OTLP_ENDPOINT` and `OTEL_SERVICE_NAME` from `docker-compose.yml`.

The non-obvious part is the propagator. Without registering the W3C TraceContext
propagator, the app spans still reach Jaeger but start their own trace, and you
lose the link back to Traefik.

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
