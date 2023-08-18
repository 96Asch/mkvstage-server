import { User } from '../model/user';
import { authenticateUser, getUsers, storeUser } from '../usecase/user';

export default Object.freeze({
    storeUser: (user: User) => storeUser(user),
    getUsers: (ids: number[], emails: string[]) => getUsers(ids, emails),
    authenticateUser: (email: string, password: string) =>
        authenticateUser(email, password),
});
