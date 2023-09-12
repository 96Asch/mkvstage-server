import { User } from '../../model/user';

export default function makeGetUsers({ userDb }: { userDb: any }) {
    return async function getUsers(ids: number[], emails: string[]): Promise<User[]> {
        const users = await userDb.read(ids, emails);

        return users;
    };
}
