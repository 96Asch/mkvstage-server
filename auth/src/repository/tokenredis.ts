import type { RedisClientType, createClient } from 'redis'

export default function makeRedisTokenRepo({redisClient}) {

    async function create(email: string, refreshtoken: string) {

    }

    async function get(email: string): Promise<string> {
        return ""
    }

    return Object.freeze({
        create,
        get
    })

}