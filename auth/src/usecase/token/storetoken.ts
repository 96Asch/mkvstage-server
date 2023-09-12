import { TokenPair } from '../../model/token';
import secrets from '../../model/token';
import { User } from '../../model/user';
import jwt from '../../util/jwt';

export default function makeCreateTokens({ redisDb }: { redisDb: any }) {
    return async function createTokens(sender: string, user: User): Promise<TokenPair> {
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

        await redisDb.create(sender, user.email, refreshToken);

        return { access: accesToken, refresh: refreshToken };
    };
}
