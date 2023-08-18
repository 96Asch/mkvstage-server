import { Pool } from 'pg';
import type { ClientConfig } from 'pg';

const { PG_USER, PG_HOST, PG_PASS, PG_DB, PG_PORT } = process.env;

const config: ClientConfig = {
    user: PG_USER,
    host: PG_HOST,
    database: PG_DB,
    password: PG_PASS,
    port: Number(PG_PORT),
};
const pgPool = new Pool(config);

console.log('Trying Postgres connection with:', config);

pgPool.connect().then(() => {
    console.log('Connected to Postgres:', pgPool);
});

export default pgPool;
