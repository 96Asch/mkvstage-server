import { User } from "../../model/user";

export default function makeGetUsers({ userDb }) {
  return async function getUsers(): Promise<User[]> {
    const users = await userDb.read();
    return users;
  };
}
