import { createToken } from "../usecase/token";

export default Object.freeze({
    createToken: (email: string) => createToken(email)
})