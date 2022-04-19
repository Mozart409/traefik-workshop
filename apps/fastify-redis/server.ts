import Fastify, { FastifyInstance, RouteShorthandOptions } from 'fastify'
import { Server, IncomingMessage, ServerResponse } from 'http'
import fastifyRedis from 'fastify-redis'

const server: FastifyInstance = Fastify({ logger: true })

server.register(fastifyRedis, {
  host: 'redis.docker.localhost',
  // password: 'your strong password here',
  port: 6379, // Redis port
})

const opts: RouteShorthandOptions = {
  schema: {
    response: {
      200: {
        type: 'object',
        properties: {
          pong: {
            type: 'string',
          },
        },
      },
    },
  },
}

server.get('/ping', opts, async (request, reply) => {
  return { pong: 'it worked!' }
})

server.get('/foo', (req: any, reply) => {
  const { redis } = server
  redis.get(req.query.key, (err, val) => {
    reply.send(err || val)
  })
})

server.post('/foo', (req: any, reply) => {
  const { redis } = server
  redis.set(req.body.key, req.body.value, (err) => {
    reply.send(err || { status: 'ok' })
  })
})

const start = async () => {
  try {
    await server.listen(4040)

    const address = server.server.address()
    const port = typeof address === 'string' ? address : address?.port
  } catch (err) {
    server.log.error(err)
    process.exit(1)
  }
}
start()
