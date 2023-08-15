import { User } from "../repository/userpg";
import { storeUser } from "../usecase/user";

export default Object.freeze({
    storeUser: (user: User) => storeUser(user)
})