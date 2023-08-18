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
import { TokenPair } from '../model/token';

const tokenRoute = Router();

tokenRoute.post('/login', async (req: Request, res: Response, next: NextFunction) => {
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

tokenRoute.post('/refresh', async (req: Request, res: Response, next: NextFunction) => {
    const { refresh } = req.body;

    if (!refresh) {
        next(makeBadRequestError('refresh field must not be empty'));

        return;
    }

    tokencontroller
        .renewAccess(refresh)
        .then((accessToken) => {
            const tokenPair: TokenPair = { access: accessToken, refresh: refresh };
            res.status(202).json({ tokens: tokenPair });
        })
        .catch(next);
});

export default tokenRoute;
