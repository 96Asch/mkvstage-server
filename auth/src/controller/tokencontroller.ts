import { User } from '../model/user';
import { createToken, removeTokensByEmail } from '../usecase/token';

export default Object.freeze({
    createToken: (sender: string, user: User) => createToken(sender, user),
    removeTokensByEmail: (email: string) => removeTokensByEmail(email),
});
