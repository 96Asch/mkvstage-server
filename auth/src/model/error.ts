class AppError extends Error {
  httpCode: number;

  constructor(message: string, httpCode: number) {
    super(message);
    this.httpCode = httpCode;
  }
}

const makeDuplicateError = (fields: string[], values: string[]) => {
  const message = `fields (${fields.join(", ")}) with value(s) (${values.join(
    ", "
  )}) already exists`;
  return new AppError(message, 401);
};

const makeEmailFormatError = (email: string) => {
  const message = `${email} is not a valid email`;
  return new AppError(message, 401);
};

const makeInternalError = () => {
  return new AppError("an error has occured on the server", 500);
};

export {
  AppError,
  makeDuplicateError,
  makeEmailFormatError,
  makeInternalError,
};
