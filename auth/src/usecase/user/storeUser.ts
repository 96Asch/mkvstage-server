import { User } from '../../model/user';
import authpass from '../../util/password';

export default function makeStoreUser({ userDb }) {
    return async function storeUser(user: User): Promise<User> {
        user.password = await authpass.hash(user.password);
        const createdUser = await userDb.create(user);

        return createdUser;
    };
}
