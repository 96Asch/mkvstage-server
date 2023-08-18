import { userDb, redisDb } from '../../repository';
import makeAuthenticateUser from './authenticateUser';
import makeGetUsers from './getusers';
import makeStoreUser from './storeuser';

const storeUser = makeStoreUser({ userDb });
const authenticateUser = makeAuthenticateUser({ userDb });
const getUsers = makeGetUsers({ userDb });

const userService = Object.freeze({
    storeUser,
    getUsers,
    authenticateUser,
});

export default userService;

export { storeUser, getUsers, authenticateUser };
