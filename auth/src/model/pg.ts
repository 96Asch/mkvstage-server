import { Pool } from 'pg'
import type { ClientConfig } from 'pg'

const config: ClientConfig = {
    user:       process.env.PG_USER,
    host:       process.env.PG_HOST,    
    database:   process.env.PG_DB,
    password:   process.env.PG_PASS,
    port:       Number(process.env.PG_PORT),
}
const pgPool = new Pool(config)

console.log("Trying Postgres connection with:", config)

pgPool.connect().then(() => {
    console.log("Connected to Postgres:", pgPool)
})

export default pgPool