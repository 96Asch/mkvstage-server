import { User } from "../model/user";
import { storeUser } from "../usecase/user";

export default Object.freeze({
  storeUser: (user: User) => storeUser(user),
});
