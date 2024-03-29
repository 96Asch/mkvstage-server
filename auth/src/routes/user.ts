import Router from 'express';
import type { NextFunction, Request, Response } from 'express';
import usercontroller from '../controller/usercontroller';
import validateEmail from '../util/validateemail';
import {
    makeBadRequestError,
    makeEmailFormatError,
    makeInternalError,
} from '../model/error';
import { User } from '../model/user';
import middleware from './middleware';
import tokencontroller from '../controller/tokencontroller';

const userRoute = Router();

userRoute.post('/', async (req: Request, res: Response, next: NextFunction) => {
    const { email, password } = req.body;

    if (!email && password) {
        next(makeBadRequestError('email and password must not be empty'));

        return;
    }

    if (!validateEmail(email)) {
        next(makeEmailFormatError(email));

        return;
    }

    const user: User = { id: 0, email: email, password: password };

    usercontroller
        .storeUser(user)
        .then((createdUser) =>
            res.status(201).json({
                user: {
                    id: createdUser.id,
                    email: createdUser.email,
                },
            })
        )
        .catch(next);
});

userRoute.get('/', async (req: Request, res: Response, next: NextFunction) => {
    const emailsParam: string =
        req.query.emails != null ? (req.query.emails as string) : '';
    const idsParam = req.query.ids != null ? (req.query.ids as string) : '';

    const emails: string[] = emailsParam.split(',').filter((email) => {
        return email != '';
    });

    const ids: number[] = idsParam.split(',').reduce((result: number[], el: string) => {
        const id = parseInt(el);

        if (!Number.isNaN(id)) {
            result.push(id);
        }

        return result;
    }, []);

    console.log(ids);

    usercontroller
        .getUsers(ids, emails)
        .then((retrievedUsers) =>
            res.status(200).json({
                users: retrievedUsers.map((user: User) => {
                    return { id: user.id, email: user.email };
                }),
            })
        )
        .catch(next);
});

userRoute.get(
    '/me',
    middleware.verifyAndExtractEmail,
    async (req: Request, res: Response, next: NextFunction) => {
        const email = res.locals.email;

        if (!email) {
            next(makeInternalError());

            return;
        }

        usercontroller
            .getUsers([], [email])
            .then((retrievedUsers) => {
                res.status(200).json({ users: retrievedUsers });
            })
            .catch(next);
    }
);

userRoute.get(
    '/logout',
    middleware.verifyAndExtractEmail,
    async (req: Request, res: Response, next: NextFunction) => {
        const email = res.locals.email;

        if (!email) {
            next(makeInternalError());

            return;
        }

        tokencontroller
            .removeTokensByEmail(email)
            .then(() => {
                res.sendStatus(202);
            })
            .catch(next);
    }
);

userRoute.get('/:id', async (req: Request, res: Response, next: NextFunction) => {
    const id: number = parseInt(req.params.id);

    if (Number.isNaN(id)) {
        next(makeBadRequestError('given id was not a number'));

        return;
    }

    usercontroller
        .getUsers([id], [])
        .then((retrievedUsers) => {
            res.status(200).json({ users: retrievedUsers });
        })
        .catch(next);
});

export default userRoute;
