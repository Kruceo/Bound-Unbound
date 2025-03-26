import { getContext } from "telefunc";
import { apiAxios, apiUrl } from "../api/api";

export async function onAuthToken() {
    try {
        const { sessionToken }: { sessionToken: string } = getContext()
        apiAxios.defaults.headers["Authorization"] = `Bearer ${sessionToken}`
        const url = apiUrl("auth/token")

        const res = await apiAxios.get(url)
        return { ok: res.status == 200, cookies: res.headers["set-cookie"] }
    }
    catch (err: any) {
        console.log(err.response.statusText)
    }
    return { ok: false }
}

export async function onAuthStatus() {
    try {
        const url = apiUrl("auth/status")
        const res = await apiAxios.get(url)
        return res.data as { Data: { AlreadyRegistered: boolean } }
    }
    catch (err: any) {
        console.log(err)
    }
    return { Data: { AlreadyRegistered: true } }
}