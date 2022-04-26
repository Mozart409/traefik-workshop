import Fastify from 'fastify'
import fetch from 'node-fetch'
import fastifyZipkin from 'fastify-zipkin'

const server = Fastify({})

server.register(fastifyZipkin, {
  serviceName: 'fastify-zipkin-example',
  servicePort: 3334,
  httpReporterUrl: 'http://zipkin.docker.localhost:9411/api/v2/spans',
})

server.get('/', (req, reply) => {
  reply.send({ hello: 'world' })
})

const start = async () => {
  try {
    await server.listen(3334)
    console.log('Server listenting on localhost:', server.server.address().port)
  } catch (err) {
    server.log.error(err)
    process.exit(1)
  }
}
start()
