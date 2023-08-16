import { makeDuplicateError } from "../model/error";
import { User, emptyUser } from "../model/user";
import type { Pool } from "pg";

export default function makeUserPg({ pgPool }) {
  async function create(user: User): Promise<User> {
    let createdUser: User = user;
    const createQuery =
      "INSERT INTO Users (email, password) VALUES ($1, $2) RETURNING *";
    try {
      const res = await pgPool.query(createQuery, [user.email, user.password]);
      createdUser.id = res.rows[0].id;
    } catch (error) {
      console.log(error.code);
      switch (error.code) {
        case "23505":
          throw makeDuplicateError(["email"], [user.email]);

        default:
          throw error;
      }
    }

    return createdUser;
  }

  async function readWithId(id: number): Promise<User[]> {
    return [];
  }

  async function readWithEmail(email: string): Promise<User> {
    const user: User = { id: 1, email: "", password: "" };
    return user;
  }

  async function update(user: User): Promise<User> {
    return { id: 0, email: "", password: "" };
  }

  async function deleteWithId(id: number) {}

  async function deleteWithEmail(email: string) {}

  return Object.freeze({
    create,
    readWithId,
    readWithEmail,
    update,
    deleteWithId,
    deleteWithEmail,
  });
}
