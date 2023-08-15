import { redisDb } from "../../repository";

import makeCreateToken from "./storetoken";

const createToken = makeCreateToken({redisDb})

const tokenService = Object.freeze({
    createToken
})

export default tokenService

export {
    createToken
}