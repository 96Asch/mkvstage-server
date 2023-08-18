import { makeBadRequestError } from '../../model/error';
import { User } from '../../model/user';
import validator from '../../util/password';

export default function makeAuthenticateUser({ userDb }) {
    return async function authenticateUser(
        email: string,
        password: string
    ): Promise<User> {
        const retrievedUsers: User[] = await userDb.read([], [email]);

        if (retrievedUsers.length <= 0) {
            throw makeBadRequestError(`no users found with email ${email}`);
        }

        const user = retrievedUsers[0];

        if (user.email !== email || !validator.validate(password, user.password)) {
            return null;
        }

        return user;
    };
}
