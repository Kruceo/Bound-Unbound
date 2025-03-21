import { getContext } from "telefunc";
import { apiAxios, apiUrl } from "../api/api";

export async function onAuthToken() {
    try {
        const { sessionToken }: { sessionToken: string } = getContext()
        apiAxios.defaults.headers["Authorization"] = `Bearer ${sessionToken}`
        const url = apiUrl("auth/token")

        const res = await apiAxios.get(url)
        return res.status == 200
    }
    catch (err: any) {
        console.log(err.response.statusText)
    }
    return false
}

export async function onAuthStatus() {
    try {
        const url = apiUrl("auth/status")
        const res = await apiAxios.get(url)
        return res.data as { Data: { AlreadyRegistered: boolean } }
    }
    catch (err: any) {
        console.log(err.response.statusText)
    }
    return { Data: { AlreadyRegistered: true } }
}