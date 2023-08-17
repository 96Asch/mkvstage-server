class AppError extends Error {
    httpCode: number;

    constructor(message: string, httpCode: number) {
        super(message);
        this.httpCode = httpCode;
    }
}

const makeDuplicateError = (fields: string[], values: string[]) => {
    const message = `fields (${fields.join(', ')}) with value(s) (${values.join(
        ', '
    )}) already exists`;
    return new AppError(message, 400);
};

const makeEmailFormatError = (email: string) => {
    const message = `${email} is not a valid email`;
    return new AppError(message, 400);
};

const makeInternalError = () => {
    return new AppError('an error has occured on the server', 500);
};

const makeNotAuthorizedError = () => {
    return new AppError('not authorized to make changes', 403);
};

const makeBadRequestError = (message: string) => {
    return new AppError(message, 401);
};

export {
    AppError,
    makeDuplicateError,
    makeEmailFormatError,
    makeInternalError,
    makeNotAuthorizedError,
    makeBadRequestError,
};
