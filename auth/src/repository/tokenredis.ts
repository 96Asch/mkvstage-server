import { makeInternalError } from '../model/error';
import { REDIS_EXP } from '../model/redis';

export default function makeRedisTokenRepo({ redisClient }) {
    async function create(email: string, refreshtoken: string) {
        try {
            await redisClient.setEx(email, REDIS_EXP, refreshtoken);
        } catch (error) {
            throw error;
        }
    }

    async function get(email: string): Promise<string> {
        try {
            const value = await redisClient.get(email);
            console.log(value);
        } catch (error) {
            throw error;
        }
        return '';
    }

    return Object.freeze({
        create,
        get,
    });
}
