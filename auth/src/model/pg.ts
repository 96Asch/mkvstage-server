import { Pool } from 'pg'


const pgPool = new Pool({
    user: 'foo',
    host: 'bar',
    database: 'auth',
    password: 'pass',
    port: 6887,
})

export default pgPool