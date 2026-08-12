import { registerOTel } from "@vercel/otel";
import { IORedisInstrumentation } from "@opentelemetry/instrumentation-ioredis";

export function register() {
  registerOTel({
    serviceName: process.env.OTEL_SERVICE_NAME ?? "nextjs",
    instrumentations: [
      new IORedisInstrumentation(),
      // IORedis instrumentation is critical for this workshop:
      // without it, only HTTP route spans are exported and Valkey/Redis
      // calls remain invisible. Including it traces all Redis operations.
    ],
  });
}
