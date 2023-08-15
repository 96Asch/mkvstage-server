import { Router } from "express";
import usercontroller from "../controller/usercontroller";


const userRoute = Router()

userRoute.post('/', (req, res) => {
    const body = req.params
    try {
        const user = {id: 1, email:"foobar@gmail.com", password:"password"}
        usercontroller.storeUser(user)        

    } catch (error) {
    }
    res.send("User Create")
    res.status(201)
})

export default userRoute