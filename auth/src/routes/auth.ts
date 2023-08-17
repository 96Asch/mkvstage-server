import { Router } from 'express';
import type { Response, Request, NextFunction } from 'express';

import {
    makeBadRequestError,
    makeEmailFormatError,
    makeNotAuthenticatedError,
} from '../model/error';
import tokencontroller from '../controller/tokencontroller';
import validateEmail from '../util/validateemail';
import usercontroller from '../controller/usercontroller';

const authRoute = Router();

authRoute.post('/login', async (req: Request, res: Response, next: NextFunction) => {
    const { senderId, email, password } = req.body;

    if (!(email && password && senderId)) {
        next(makeBadRequestError('email, password and senderId cannot be empty'));

        return;
    }

    if (!validateEmail(email)) {
        next(makeEmailFormatError(email));

        return;
    }

    console.log(email, password);
    try {
        const user = await usercontroller.authenticateUser(email, password);

        if (!user) {
            next(makeNotAuthenticatedError());

            return;
        }

        const tokenPair = await tokencontroller.createToken(senderId, user);
        res.status(200).json({ tokens: tokenPair });
    } catch (error) {
        next(error);
    }
});

export default authRoute;
