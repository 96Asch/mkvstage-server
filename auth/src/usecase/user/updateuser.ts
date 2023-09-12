import { makeNotAuthorizedError } from '../../model/error';
import { User } from '../../model/user';

export default function makeUpdateUser({ userDb }: { userDb: any }) {
    return async function updateUser(update: User, principal: User) {
        if (update.id != principal.id) {
            throw makeNotAuthorizedError();
        }

        const updatedUser = await userDb.update(update);

        return updatedUser;
    };
}
