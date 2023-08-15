

export default function makeCreateToken({redisDb}) {

   return async function createToken(email: string): Promise<string> {
        const refreshToken = "foobar"
        redisDb.create(email, refreshToken)
        return refreshToken
   }
}