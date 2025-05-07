import { AxiosError, AxiosInstance } from "axios"
import { ApiResponse, apiUrl } from "./api"

type Fetcher = { axios: AxiosInstance }

export async function onCreateRegisterRequest(this: Fetcher, role: string): Promise<CreateRegisterRequestResponse> {
    const url = apiUrl(`/auth/register/request`)
    try {
        const res = await this.axios.post(url, { roleId: role })
        return await res.data
    } catch (error) {
        return (error as AxiosError).response?.data as any
    }
}


export async function onDeleteRoles(this: Fetcher, roles: string[]): Promise<boolean> {
    const url = apiUrl(`/admin/roles`)
    try {
        const res = await this.axios.delete(url, {
            data: roles.map(e => ({ id: e }))
        })
        return await res.data
    } catch (error) {
        return (error as AxiosError).response?.data as any
    }
}

export async function onGetRoles(this: Fetcher): Promise<OnGetRolesResponse> {
    const url = apiUrl(`/admin/roles`)
    try {
        const res = await this.axios.get(url)
        return await res.data
    } catch (error) {
        return (error as AxiosError).response?.data as any
    }
}

export async function onPostRoles(this: Fetcher, roles: { name: string, permissions: string[] }[]): Promise<OnPostRolesResponse> {
    const url = apiUrl(`/admin/roles`)
    try {
        const res = await this.axios.post(url, roles)
        return await res.data
    } catch (error) {
        return (error as AxiosError).response?.data as any
    }
}

export async function onGetUsers(this: Fetcher): Promise<GetUsersResponse> {
    const url = apiUrl(`/admin/users`)
    try {
        const res = await this.axios.get(url)
        return await res.data
    } catch (error) {
        return (error as AxiosError).response?.data as any
    }
}

export async function onDeleteUsers(this: Fetcher, ids: string[]): Promise<GetUsersResponse> {
    const url = apiUrl(`/admin/users`)
    try {
        const res = await this.axios.delete(url, { data: ids.map(f => ({ id: f })) })
        return await res.data
    } catch (error) {
        return (error as AxiosError).response?.data as any
    }
}

export async function onAuthToken(this: Fetcher) {
    try {
        const url = apiUrl("auth/token")
        const res = await this.axios.get(url)
        const rawCookies = res.headers["set-cookie"]??[]

        const cookieMap = new Map();

        rawCookies.forEach(cookieStr => {
            const [cookiePair] = cookieStr.split(";"); // Pega só a parte "chave=valor"
            const [key, value] = cookiePair.split("=").map(s => s.trim());
            cookieMap.set(key, value);
        });

        return { ok: res.status == 200, cookies: cookieMap, rawCookies }
    }
    catch (err: any) {
        console.error("error: authTokenW")
    }
    return { ok: false }
}

export async function onUnblockAction(this: Fetcher, connectionId: string, domains: string[]): Promise<ApiResponse<any>> {
    const url = apiUrl(`v1/connections/${connectionId}/blocks`)
    try {
        const res = await this.axios.delete(url, {
            headers: { "Content-Type": "application/json" },
            data: { Names: domains },
        })
        return await res.data
    } catch (error: any) {
        return error.response?.data
    }
}

export async function onBlockAction(this: Fetcher, connectionId: string, domains: string[]): Promise<ApiResponse<any>> {
    const url = apiUrl(`v1/connections/${connectionId}/blocks`)
    try {
        const res = await this.axios.post(url, {
            Names: domains
        }, {
            headers: { "Content-Type": "application/json" }
        })
        return await res.data
    } catch (error: any) {
        return error.response?.data
    }
}

export async function onNewRedirectAction(this: Fetcher, connectionId: string, from: string, recordType: string, to: string, localZone: boolean): Promise<ApiResponse<any>> {
    const url = apiUrl(`v1/connections/${connectionId}/redirects`)
    try {
        const res = await this.axios.post(url, {
            From: from,
            To: to,
            RecordType: recordType,
            LocalZone: localZone,
        }, {
            headers: { "Content-Type": "application/json" },
        })
        return await res.data
    } catch (error: any) {
        return error.response?.data
    }
}

export async function onDeleteRedirectAction(this: Fetcher, connectionId: string, domain: string[]) {
    const url = apiUrl(`v1/connections/${connectionId}/redirects`)
    let responses
    for (let dom of domain) {
        try {
            const res = await this.axios.delete(url, {
                headers: { "Content-Type": "application/json" },
                data: { Domain: dom },
            })
            responses = res.data
        } catch (error) {
            return (error as AxiosError).response?.data
        }
    }
    return responses
}

export async function onGetConfigHash(this: Fetcher, connectionId: string): Promise<ConfigHashResponse> {
    const url = apiUrl(`v1/connections/${connectionId}/confighash`)
    try {
        const res = await this.axios.get(url)
        return await res.data
    } catch (error) {
        return (error as AxiosError).response?.data as any
    }
}

export async function onReloadActions(this: Fetcher, connectionId: string): Promise<ApiResponse<any>> {
    const url = apiUrl(`v1/connections/${connectionId}/reload`)
    try {
        const res = await this.axios.post(url)
        return await res.data
    } catch (error: any) {
        return error.response?.data
    }
}

export async function onGetNodes(this: Fetcher): Promise<ApiResponse<{ name: string, remoteAddress: string }[]>> {
    try {
        const url = apiUrl("/v1/connections")
        const res = await this.axios.get(url)
        return res.data as ApiResponse<{ name: string, remoteAddress: string }[]>
    }
    catch (error: any) {

    }
    return { message: "", data: [] }
}



export async function onLoginAction(this: Fetcher, username: string, password: string): Promise<ApiResponse<{ token: string }>> {

    const encoder = new TextEncoder();
    const data = encoder.encode(password);
    const hashBuffer = await crypto.subtle.digest('SHA-256', data);
    const hashArray = Array.from(new Uint8Array(hashBuffer));
    const hashedPassword = hashArray.map(byte => byte.toString(16).padStart(2, '0')).join('');

    const url = apiUrl("/auth/login")
    try {
        const res = await this.axios.post(url, {
            User: username,
            Password: hashedPassword
        })
        return res.data
    } catch (error: any) {
        return error.response.data
    }
}



export async function onResetPassword(this: Fetcher, routeId: string, password: string): Promise<ApiResponse<{ secretCode: string }>> {

    const encoder = new TextEncoder();
    const url = apiUrl("/auth/reset/pwd")

    const data = encoder.encode(password);
    const hashBuffer = await crypto.subtle.digest('SHA-256', data);
    const hashArray = Array.from(new Uint8Array(hashBuffer));
    const hashedPassword = hashArray.map(byte => byte.toString(16).padStart(2, '0')).join('');

    try {
        const res = await this.axios.post(url, {
            routeId,
            password: hashedPassword
        })

        return res.data
    } catch (error: any) {
        console.log(error)
        return error.response.data
    }


}



export async function onRegisterAction(this: Fetcher, username: string, password: string, routeId?: string): Promise<ApiResponse<{ secretCode: string }>> {
    const encoder = new TextEncoder();
    const url = apiUrl("/auth/register")

    const data = encoder.encode(password);
    const hashBuffer = await crypto.subtle.digest('SHA-256', data);
    const hashArray = Array.from(new Uint8Array(hashBuffer));
    const hashedPassword = hashArray.map(byte => byte.toString(16).padStart(2, '0')).join('');

    try {
        const res = await this.axios.post(url, {
            user: username,
            password: hashedPassword,
            routeId
        })

        return res.data
    } catch (error: any) {
        console.log(error)
        return error.response.data
    }
}



export async function onResetAccount(this: Fetcher, user: string, secretCode: string): Promise<ApiResponse<{ routeId: string }>> {
    const url = apiUrl("/auth/reset")

    try {
        const res = await this.axios.post(url, {
            secretCode,
            user
        })

        return res.data
    } catch (error: any) {
        console.log(error)
        return error.response.data
    }
}


export async function onPostBind(this: Fetcher, binds: { nodeId: string, roleId: string }[]): Promise<ApiResponse<{ id: string }[]>> {
    const url = apiUrl("/admin/roles/bind/nodes")

    try {
        const res = await this.axios.post(url, binds)

        return res.data
    } catch (error: any) {
        console.log(error)
        return error.response.data
    }
}

export async function onGetBinds(this: Fetcher): Promise<ApiResponse<{ id: string, node: { id: string, name: string }, role: { id: string, name: string, permissions: string[] } }[]>> {
    const url = apiUrl("/admin/roles/bind/nodes")

    try {
        const res = await this.axios.get(url)
        return res.data
    } catch (error: any) {
        console.log("error:",error.message)
        return error.response.data
    }
}

export async function onDeleteBinds(this: Fetcher, ids: string[]): Promise<ApiResponse<boolean | undefined>> {
    const url = apiUrl("/admin/roles/bind/nodes")

    try {
        const res = await this.axios.delete(url, {
            data: ids.map(f => ({ id: f }))
        })
        return res.data
    } catch (error: any) {
        console.error("error:",error.message)
        return error.response.data
    }
}