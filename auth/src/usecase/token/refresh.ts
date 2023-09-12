import { makeNotAuthenticatedError, makeNotAuthorizedError } from '../../model/error';
import secrets, { JWTRefreshPayload } from '../../model/token';
import jwt from '../../util/jwt';

export default function makeRenewAccess({
    userDb,
    redisDb,
}: {
    userDb: any;
    redisDb: any;
}) {
    return async function renewAccess(refresh: string): Promise<string> {
        const refreshPayload = jwt.validate(
            refresh,
            secrets.JWT_REFRESH
        ) as JWTRefreshPayload;

        const user = await userDb.read([refreshPayload.id], []);
        const retrievedRefreshs = await redisDb.read(user.email);

        if (refresh! in retrievedRefreshs) {
            throw makeNotAuthenticatedError();
        }

        const accesToken = jwt.create(
            { email: user.email },
            secrets.JWT_ACCESS,
            secrets.JWT_ACCESS_EXP
        );

        return accesToken;
    };
}
