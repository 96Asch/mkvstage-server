import { redisDb } from '../../repository';
import makeRemoveTokensByEmail from './removetokensbyemail';

import makeCreateToken from './storetoken';

const createToken = makeCreateToken({ redisDb });
const removeTokensByEmail = makeRemoveTokensByEmail({ redisDb });

const tokenService = Object.freeze({
    createToken,
    removeTokensByEmail,
});

export default tokenService;

export { createToken, removeTokensByEmail };
