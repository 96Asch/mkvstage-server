import type { RedisClientType} from 'redis'
import { redisClient } from '../model'

export default function makeRedisTokenRepo({redisClient: RedisClientType}) {

    async function create(email: string, refreshtoken: string) {
        redisClient.set("email", email)
    }

    async function get(email: string): Promise<string> {
        return ""
    }

    return Object.freeze({
        create,
        get
    })

}