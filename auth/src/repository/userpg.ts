export interface User {
    id: number
    email: string
    password: string
}

export const emptyUser = {
    id: 0,
    email: "",
    password: "",
}

export default function makeUserPg({ pgPool }) {

    async function create(user: User) {
        console.log('Repository: Create:', user)
    }

    async function readWithId(id: number): Promise<User[]> {
        return []
    }

    async function readWithEmail(email: string): Promise<User> {
        const user: User = {id: 1, email: "", password: ""}
        return user
    }

    async function update(user: User): Promise<User> {
        return {id: 0, email:"", password:""}
    }

    async function deleteWithId(id: number) {

    }

    async function deleteWithEmail(email: string) {

    }

    return Object.freeze({
        create,
        readWithId,
        readWithEmail,
        update,
        deleteWithId,
        deleteWithEmail,
    })
}

