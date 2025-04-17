import { render } from "vike/abort";
import { useAPI } from "../../api/api.js";

export type Data = Awaited<ReturnType<typeof data>>;

export const data = async () => {
    const api = useAPI()
    
    const users = await api.onGetUsers()
    if (users.error || !users.data) {
        console.error("problem with users")
        throw render(500)
    }
    const roles = await api.onGetRoles()
    if (roles.error || !roles.data) {
        console.error("problem with roles")
        throw render(500)
    }

    console.log(roles)

    return { users: users.data, roles: roles.data }
};