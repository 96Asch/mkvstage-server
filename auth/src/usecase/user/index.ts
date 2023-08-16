import { userDb } from "../../repository";
import makeAuthorizeUser from "./authorizeUser";

import makeStoreUser from "./storeUser";

const storeUser = makeStoreUser({userDb})
const authorizeUser = makeAuthorizeUser({userDb})

const userService = Object.freeze({
    storeUser,
    authorizeUser,
})

export default userService

export {
    storeUser,
    authorizeUser
}