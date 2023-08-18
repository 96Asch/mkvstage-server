import bcrypt from 'bcrypt';

const saltRounds: number = 10;

const hash = async (password: string): Promise<string> => {
    return bcrypt.hash(password, saltRounds);
};

const validate = async (password: string, hash: string): Promise<boolean> => {
    return bcrypt.compare(password, hash);
};

export default Object.freeze({
    hash,
    validate,
});
