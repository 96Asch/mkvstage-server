import { userDb, redisDb } from '../../repository';
import makeAuthorizeUser from './authorizeuser';
import makeGetUsers from './getusers';
import makeStoreUser from './storeuser';

const storeUser = makeStoreUser({ userDb });
const authorizeUser = makeAuthorizeUser({ userDb, redisDb });
const getUsers = makeGetUsers({ userDb });

const userService = Object.freeze({
    storeUser,
    getUsers,
    authorizeUser,
});

export default userService;

export { storeUser, getUsers, authorizeUser };
