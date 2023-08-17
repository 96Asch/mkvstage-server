import { createClient } from 'redis';

const { REDIS_EXP_HOURS, REDIS_PORT, REDIS_HOST } = process.env;

const redisClient = createClient({
    url: `redis://${REDIS_HOST}:${REDIS_PORT}`,
});

redisClient.connect().then(() => {
    console.log('Redis connected');
});

export default redisClient;

export const REDIS_EXP: number = parseInt(REDIS_EXP_HOURS) * 3600;
