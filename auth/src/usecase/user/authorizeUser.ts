import { User } from "../../model/user";

export default function makeAuthorizeUser({ userDb }) {
  return async function authorizeUser(
    callback: (e?: Error) => void,
    user: User
  ): Promise<boolean> {
    return false;
  };
}
