import { getContext } from "telefunc";
import { apiAxios, apiUrl, ApiResponse } from "../api/api";

export async function onAuthToken() {
    try {
        const { sessionToken }: { sessionToken: string } = getContext()
        apiAxios.defaults.headers["Authorization"] = `Bearer ${sessionToken}`
        const url = apiUrl("auth/token")
        console.log('fetching')

        
        const res = await apiAxios.get(url)

        return { ok: res.status == 200, cookies: res.headers["set-cookie"] }
    }
    catch (err: any) {
        console.log(err.response.data)
    }
    return { ok: false }
}

export async function onAuthStatus(): Promise<ApiResponse<{ alreadyRegistered: boolean }>> {
    try {
        const url = apiUrl("auth/status")
        const res = await apiAxios.get(url)
        return res.data
    }
    catch (err: any) {
        // console.log(err)
    }
    return { data: { alreadyRegistered: true }, message: "" }
}