export default function makeRemoveTokensByEmail({ redisDb }: { redisDb: any }) {
    return async function removeTokensByEmail(email: string): Promise<void> {
        redisDb.del(email);
    };
}
