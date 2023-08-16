import { Router } from "express";
import usercontroller from "../controller/usercontroller";
import validateEmail from "../util/validateemail";
import { makeEmailFormatError } from "../model/error";

const userRoute = Router();

userRoute.post("/", async (req, res, next) => {
  const { email, password } = req.body;

  if (!validateEmail(email)) {
    next(makeEmailFormatError(email));
    return;
  }

  const user = { id: 0, email: email, password: password };

  usercontroller
    .storeUser(user)
    .then((createdUser) =>
      res.status(201).json({ id: createdUser.id, email: createdUser.email })
    )
    .catch(next);
});

export default userRoute;
