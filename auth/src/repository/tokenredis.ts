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

    async function read(email: string): Promise<string[]> {
        const key = `${email}:*`;

        return new Promise((resolve, _) => {
            var stream = redisClient.scanStream({
                match: key,
            });

            var keys = [];
            stream.on('data', function (resultKeys) {
                for (var i = 0; i < resultKeys.length; i++) {
                    keys.push(resultKeys[i]);
                }
            });

            stream.on('end', function () {
                resolve(keys);
            });
        });
    }

    async function del(email: string): Promise<void> {
        const stream = redisClient.scanStream({
            match: `${email}:*`,
        });

        stream.on('data', function (keys) {
            if (keys.length) {
                redisClient.unlink(keys);
            }
        });

        stream.on('end', function () {
            console.log('Finished deleting keys');
        });
    }

    return Object.freeze({
        create,
        read,
        del,
    });
}
