import { User } from '../../model/user';

export default function makeAuthorizeUser({ userDb }) {
    return async function authorizeUser(user: User): Promise<boolean> {
        return false;
    };
}
