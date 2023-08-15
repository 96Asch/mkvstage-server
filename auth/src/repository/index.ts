import * as models from "../model"

import makeRedisTokenRepo from "./tokenredis"

const redisDb = makeRedisTokenRepo(models)

export {
    redisDb
}