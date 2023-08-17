import type { QueryResult, QueryResultRow } from 'pg';
import { makeDuplicateError, makeInternalError } from '../model/error';
import { User, emptyUser } from '../model/user';
import { query } from 'express';

export default function makeUserPg({ pgPool }) {
    async function create(user: User): Promise<User> {
        let createdUser: User = user;
        const createQuery =
            'INSERT INTO Users (email, password) VALUES ($1, $2) RETURNING *';

        try {
            const res = await pgPool.query(createQuery, [user.email, user.password]);
            createdUser.id = res.rows[0].id;
        } catch (error) {
            console.log(error.code);
            switch (error.code) {
                case '23505':
                    throw makeDuplicateError(['email'], [user.email]);

                default:
                    throw error;
            }
        }

        return createdUser;
    }

    async function read(ids: number[], emails: string[]): Promise<User[]> {
        let createQuery: string = 'SELECT * FROM users';
        let hasFilter: boolean = false;
        let params: string[] = [];
        let paramCount: number = 1;

        if (ids.length > 0) {
            createQuery += ` WHERE id in ($${paramCount})`;
            hasFilter = true;
            params.push(ids.join(', '));
            paramCount++;
        }

        if (emails.length > 0) {
            if (!hasFilter) {
                createQuery += ' WHERE ';
                hasFilter = true;
            } else {
                createQuery += ' AND ';
            }

            createQuery += `email in ($${paramCount})`;
            params.push(emails.join(', '));
            paramCount++;
        }

        console.log('Query:', createQuery, params);

        try {
            const res = await pgPool.query(createQuery, params);

            return res.rows.map((row: QueryResultRow) => {
                const user: User = {
                    id: row.id,
                    email: row.email,
                    password: row.password,
                };
                return user;
            });
        } catch (error) {
            throw makeInternalError();
        }
    }

    // async function readWithId(ids: number[]): Promise<User[]> {
    //     const query = 'SELECT * FROM users WHERE id in ($1)';

    //     try {
    //         const res = await pgPool.query(query, [ids.join(', ')]);

    //         return res.rows.map((row: QueryResultRow) => {
    //             const user: User = {
    //                 id: row.id,
    //                 email: row.id,
    //                 password: row.pasword,
    //             };

    //             return user;
    //         });
    //     } catch (error) {
    //         throw makeInternalError();
    //     }
    // }

    // async function readWithEmail(emails: string[]): Promise<User> {
    //     const query = 'SELECT * FROM users WHERE email in ($1)';

    //     try {
    //         const res = await pgPool.query(query, [emails.join(', ')]);

    //         return res.rows.map((row: QueryResultRow) => {
    //             const user: User = {
    //                 id: row.id,
    //                 email: row.id,
    //                 password: row.pasword,
    //             };

    //             return user;
    //         });
    //     } catch (error) {
    //         throw makeInternalError();
    //     }
    // }

    async function update(user: User): Promise<User> {
        const query = 'UPDATE users SET password = $1 WHERE id = $2';

        try {
            await pgPool.query(query, [user.password, user.id]);

            return user;
        } catch (error) {
            throw makeInternalError();
        }
    }

    return Object.freeze({
        create,
        read,
        update,
    });
}
