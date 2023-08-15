import { User, emptyUser } from "../../repository/userpg";

export default function makeStoreUser( {userDb} ) {

    return async function storeUser(user: User): Promise<User>{

        console.log("Use-Case: storeUser:", user)
        userDb.create(user)
        return emptyUser
    }
}