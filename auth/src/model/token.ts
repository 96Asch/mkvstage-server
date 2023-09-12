export interface TokenPair {
    access: string;
    refresh: string;
}

const { JWT_ACCESS, JWT_REFRESH, JWT_ACCESS_EXP, JWT_REFRESH_EXP } = process.env;

export interface JWTAccessPayload {
    email: string;
}

export interface JWTRefreshPayload {
    id: number;
}

export type JWTPayload = JWTAccessPayload | JWTRefreshPayload;

if (!JWT_ACCESS || !JWT_REFRESH || !JWT_ACCESS_EXP || !JWT_REFRESH_EXP) {
    process.exit(-1);
}

export default Object.freeze({
    JWT_ACCESS,
    JWT_ACCESS_EXP,
    JWT_REFRESH,
    JWT_REFRESH_EXP,
});
