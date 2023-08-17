import { makeBadRequestError, makeNotAuthenticatedError } from '../../model/error';
import { TokenPair } from '../../model/token';
import secrets from '../../model/token';
import { User } from '../../model/user';
import validator from '../../util/password';
import jwt from '../../util/jwt';

export default function makeAuthorizeUser({ userDb, redisDb }) {
    return async function authorizeUser(
        email: string,
        password: string
    ): Promise<TokenPair> {
        const retrievedUsers: User[] = await userDb.read([], [email]);

        if (retrievedUsers.length <= 0) {
            throw makeBadRequestError(`no users found with email ${email}`);
        }

        const user = retrievedUsers[0];

        if (user.email !== email || !validator.validate(password, user.password)) {
            throw makeNotAuthenticatedError();
        }

        const accesToken = jwt.create(
            { email: user.email },
            secrets.JWT_ACCESS,
            secrets.JWT_ACCESS_EXP
        );
        const refreshToken = jwt.create(
            { id: user.id },
            secrets.JWT_REFRESH,
            secrets.JWT_REFRESH_EXP
        );

        await redisDb.create(email, refreshToken);

        return { access: accesToken, refresh: refreshToken };
    };
}
