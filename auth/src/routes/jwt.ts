import type { Request, Response, NextFunction } from 'express';
import { makeBadRequestError, makeNotAuthenticatedError } from '../model/error';
import jwt from '../util/jwt';
import secrets, { JWTAccessPayload } from '../model/token';

export default function verifyAndExtractEmail(
    req: Request,
    res: Response,
    next: NextFunction
) {
    const bearerHeader: string = req.headers.authorization;

    if (!bearerHeader) {
        next(makeBadRequestError('authorization header is missing'));

        return;
    }

    const bearer = bearerHeader.split(' ');

    if (bearer.length != 2) {
        next(makeBadRequestError("authorization header must contain: 'Bearer [TOKEN]'"));

        return;
    }

    try {
        const payload = jwt.validate(bearer[1], secrets.JWT_ACCESS) as JWTAccessPayload;
        res.locals.email = payload.email;
    } catch (error) {
        next(error);
    }

    next();
}
