import type { Request, Response, NextFunction } from 'express';
import { AppError, makeBadRequestError, makeNotAuthenticatedError } from '../model/error';
import jwt from '../util/jwt';
import secrets, { JWTAccessPayload } from '../model/token';

const verifyAndExtractEmail = (req: Request, res: Response, next: NextFunction) => {
    const bearerHeader = req.headers.authorization;

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
};

const logger = async (req: Request, res: Response, next: NextFunction) => {
    res.on('finish', () => {
        console.log(req.method, decodeURI(req.url), res.statusCode, res.statusMessage);
    });
    next();
};

const errorHandler = async (
    err: Error,
    _: Request,
    res: Response,
    next: NextFunction
) => {
    if (err instanceof AppError) {
        const appError = err as AppError;
        res.status(appError.httpCode);
    } else {
        res.status(500);
    }
    console.error(err.stack);
    res.json({ error: err.message });
    next();
};

export default Object.freeze({
    verifyAndExtractEmail,
    logger,
    errorHandler,
});
