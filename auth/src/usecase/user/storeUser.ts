import { User, emptyUser } from "../../model/user";

export default function makeStoreUser({ userDb }) {
  return async function storeUser(user: User): Promise<User> {
    const createdUser = await userDb.create(user);

    return createdUser;
  };
}
