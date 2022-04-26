import Fastify from 'fastify'
const fetch = require('node-fetch')
import fastifyZipkin from 'fastify-zipkin'
const server = Fastify({})

server.register(fastifyZipkin, {
  serviceName: 'fastify-zipkin-example',
  servicePort: 3334,
  httpReporterUrl: 'http://localhost:9411/api/v2/spans',
})

server.get('/', (req, reply) => {
  reply.send({ hello: 'world' })
})

server.get('/ping', async (request, reply) => {
  return { pong: 'it worked!' }
})

const start = async () => {
  try {
    await server.listen(3334)

    const address = server.server.address()
    const port = typeof address === 'string' ? address : address?.port
  } catch (err) {
    server.log.error(err)
    process.exit(1)
  }
}
start()
