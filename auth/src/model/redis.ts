import { Redis } from 'ioredis';

const { REDIS_EXP_HOURS, REDIS_PORT, REDIS_HOST } = process.env;

const redisClient = new Redis({
    host: REDIS_HOST,
    port: parseInt(REDIS_PORT),
});

export default redisClient;

export const REDIS_EXP: number = parseInt(REDIS_EXP_HOURS) * 3600;
