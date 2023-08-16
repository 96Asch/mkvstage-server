import { userDb } from "../../repository";
import makeAuthorizeUser from "./authorizeUser";
import makeGetUsers from "./readUser";

import makeStoreUser from "./storeUser";

const storeUser = makeStoreUser({ userDb });
const authorizeUser = makeAuthorizeUser({ userDb });
const getUsers = makeGetUsers({ userDb });

const userService = Object.freeze({
  storeUser,
  getUsers,
  authorizeUser,
});

export default userService;

export { storeUser, getUsers, authorizeUser };
