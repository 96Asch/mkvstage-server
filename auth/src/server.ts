import express from 'express';
import type { Request, Response, NextFunction } from 'express';
import { routes } from './routes';
import { AppError } from './model/error';

const app = express();

const logger = (req: Request, res: Response, next: NextFunction) => {
    res.on('finish', () => {
        console.log(req.method, decodeURI(req.url), res.statusCode, res.statusMessage);
    });
    next();
};

const errorHandler = (err: Error, _: Request, res: Response, next: NextFunction) => {
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

app.use(express.json());
app.use('/', logger, routes);
app.use(errorHandler);

app.listen(9080, () => {
    console.log('Listening on port:', 9080);
});
