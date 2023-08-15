import { Router } from 'express'

import controller from "../controller/tokencontroller"

const authRoute = Router()

authRoute.post("/", async (req, res) => {
    const email = req.body.email
    let token = ""
    try {
        token = await controller.createToken(email)
        
    } catch (error) {
        
    }

    res.send(token)
    res.status(200)
})

export default authRoute