const accessSecret = process.env.JWT_ACCESS;
const refreshSecret = process.env.JWT_REFRESH;

export interface JWTAccessPayload {
    email: string;
}

export interface JWTRefreshPayload {
    id: number;
}

export type JWTPayload = JWTAccessPayload | JWTRefreshPayload;

export default Object.freeze({
    accessSecret,
    refreshSecret,
});
