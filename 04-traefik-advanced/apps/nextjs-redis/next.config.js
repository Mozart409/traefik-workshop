module.exports = {
  output: "standalone",
  serverExternalPackages: [
    "@vercel/otel",
    "@opentelemetry/api",
    "@opentelemetry/instrumentation-ioredis",
  ],
};
