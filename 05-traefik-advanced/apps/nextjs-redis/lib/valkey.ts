import Redis from "ioredis";

const valkey = new Redis(process.env.REDIS_URL!);

export default valkey;
