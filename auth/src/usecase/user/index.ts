import { userDb } from "../../repository";

import makeStoreUser from "./storeUser";

const storeUser = makeStoreUser({userDb})

const userService = Object.freeze({
    storeUser
})

export default userService

export {
    storeUser
}