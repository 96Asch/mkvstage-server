import jwt from 'jsonwebtoken';
import { makeBadRequestError } from '../model/error';
import { JWTPayload } from '../model/jwt';

const create = (payload: JWTPayload, secret: string): string => {
    const token: string = jwt.sign(payload, secret, {
        algorithm: 'HS256',
        issuer: process.env.JWT_ISS,
        expiresIn: process.env.JWT_EXP,
    });

    return token;
};

const validate = (token: string, secret: string): JWTPayload => {
    try {
        const decoded = jwt.verify(token, secret, {
            audience: process.env.JWTAUD,
        }) as JWTPayload;

        return decoded;
    } catch (error) {
        switch (error.name) {
            case 'NotBeforeError':
                throw makeBadRequestError('token is not active');
            case 'TokenExpiredError':
                throw makeBadRequestError('token is expired');
            default:
                throw makeBadRequestError(error.message);
        }
    }
};

export default Object.freeze({
    create,
    validate,
});
