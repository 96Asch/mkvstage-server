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

    async function del(email: string): Promise<void> {
        console.log('delete');
        const stream = redisClient.scanStream({
            match: `${email}:*`,
        });
        console.log('stream');

        stream.on('data', function (keys) {
            console.log(keys);
            if (keys.length) {
                redisClient.unlink(keys);
            }
        });

        stream.on('end', function () {
            console.log('done');
        });
    }

    return Object.freeze({
        create,
        get,
        del,
    });
}
