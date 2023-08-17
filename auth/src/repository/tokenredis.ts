import { makeInternalError } from '../model/error';
import { REDIS_EXP } from '../model/redis';

export default function makeRedisTokenRepo({ redisClient }) {
    async function create(sender: string, email: string, refreshtoken: string) {
        const key = `${email}:${sender}`;

        try {
            await redisClient.set(key, refreshtoken, 'EX', REDIS_EXP);
        } catch (error) {
            console.error(error);
            throw makeInternalError();
        }
    }

    async function get(sender: string, email: string): Promise<string> {
        const key = `${email}:${sender}`;

        try {
            const value = await redisClient.get(key);
            console.log(value);

            return value;
        } catch (error) {
            console.error(error);
            throw makeInternalError();
        }
    }

    return Object.freeze({
        create,
        get,
    });
}
