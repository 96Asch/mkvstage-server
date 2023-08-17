import { createToken } from '../usecase/token';
import { authorizeUser } from '../usecase/user';

export default Object.freeze({
    createToken: (email: string) => createToken(email),
    authorizeUser: (email: string, password: string) => authorizeUser(email, password),
});
