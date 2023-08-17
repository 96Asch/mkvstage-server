import { makeInternalError } from '../model/error';
import { REDIS_EXP } from '../model/redis';

export default function makeRedisTokenRepo({ redisClient }) {
    async function create(email: string, refreshtoken: string) {
        try {
            await redisClient.set(email, refreshtoken, 'EX', REDIS_EXP);
        } catch (error) {
            console.error(error);
            throw makeInternalError();
        }
    }

    async function get(email: string): Promise<string> {
        try {
            const value = await redisClient.get(email);
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
