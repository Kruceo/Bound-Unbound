import { AxiosError } from "axios"
import { apiAxios, apiUrl } from "../../api/api"

export async function onCreateRegisterRequest(role: string): Promise<CreateRegisterRequestResponse> {
    const url = apiUrl(`/auth/register/request`)
    try {
        console.log({ roleId:role })
        const res = await apiAxios.post(url, { roleId:role })
        return await res.data
    } catch (error) {
        return (error as AxiosError).response?.data as any
    }
}


export async function onDeleteRoles(roles: string[]): Promise<boolean> {
    const url = apiUrl(`/admin/roles`)
    try {
        const res = await apiAxios.delete(url, {
            data: roles.map(e => ({ id: e }))
        })
        return await res.data
    } catch (error) {
        return (error as AxiosError).response?.data as any
    }
}

export async function onGetRoles(): Promise<OnGetRolesResponse> {
    const url = apiUrl(`/admin/roles`)
    try {
        const res = await apiAxios.get(url)
        return await res.data
    } catch (error) {
        return (error as AxiosError).response?.data as any
    }
}

export async function onPostRoles(roles: { name: string, permissions: string[] }[]): Promise<OnPostRolesResponse> {
    const url = apiUrl(`/admin/roles`)
    try {
        const res = await apiAxios.post(url, roles)
        return await res.data
    } catch (error) {
        return (error as AxiosError).response?.data as any
    }
}


export async function onGetUsers(): Promise<GetUsersResponse> {
    const url = apiUrl(`/admin/users`)
    try {
        const res = await apiAxios.get(url)
        return await res.data
    } catch (error) {
        return (error as AxiosError).response?.data as any
    }
}



