import express from 'express'
import type { Request, Response, NextFunction } from 'express'
import { routes } from './routes'

const app = express()

const logger = (req: Request, res: Response, next: NextFunction) => {
    res.on('finish', () => {
        console.log(req.method, decodeURI(req.url), res.statusCode, res.statusMessage)
    })
    next()
}

app.use(express.json())
app.use("/", logger, routes)

app.listen(9080, ()=>{
   console.log("Listening on port:", 9080) 
})

