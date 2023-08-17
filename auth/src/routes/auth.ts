import { Router } from 'express';
import type { Response, Request, NextFunction } from 'express';

import { makeBadRequestError, makeEmailFormatError } from '../model/error';
import tokencontroller from '../controller/tokencontroller';
import validateEmail from '../util/validateemail';

const authRoute = Router();

authRoute.post('/login', async (req: Request, res: Response, next: NextFunction) => {
    const { email, password } = req.body;

    if (!(email && password)) {
        next(makeBadRequestError('email and password cannot be empty'));
        return;
    }

    if (!validateEmail(email)) {
        next(makeEmailFormatError(email));

        return;
    }

    console.log(email, password);
    try {
        const tokenPair = await tokencontroller.authorizeUser(email, password);

        res.status(200).json({ tokens: tokenPair });
    } catch (error) {
        next(error);
    }
});

export default authRoute;
