import models from '../model';

import makeRedisTokenRepo from './tokenredis';
import makeUserPg from './userpg';

const redisDb = makeRedisTokenRepo(models);
const userDb = makeUserPg(models.pgPool);

export { redisDb, userDb };
