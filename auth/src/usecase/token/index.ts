import { redisDb, userDb } from '../../repository';
import makeRenewAccess from './refresh';
import makeRemoveTokensByEmail from './removetokensbyemail';

import makeCreateToken from './storetoken';

const createToken = makeCreateToken({ redisDb });
const removeTokensByEmail = makeRemoveTokensByEmail({ redisDb });
const renewAccess = makeRenewAccess({ userDb, redisDb });

const tokenService = Object.freeze({
    createToken,
    removeTokensByEmail,
    renewAccess,
});

export default tokenService;

export { createToken, removeTokensByEmail, renewAccess };
