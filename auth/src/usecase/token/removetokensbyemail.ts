export default function makeRemoveTokensByEmail({ redisDb }) {
    return async function removeTokensByEmail(email: string): Promise<void> {
        redisDb.del(email);
    };
}
