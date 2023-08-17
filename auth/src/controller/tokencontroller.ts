import { User } from '../model/user';
import { createToken } from '../usecase/token';

export default Object.freeze({
    createToken: (sender: string, user: User) => createToken(sender, user),
});
