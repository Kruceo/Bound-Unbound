interface APIResponse<T = any> {
    error?: false, errorCode?: string, message: string, data?: T
}

interface ConfigHashResponse extends APIResponse {
    data?: { Hash: string }
}

interface CreateRegisterRequestResponse extends APIResponse {
    data?: { routeId: string }
}

interface GetUsersResponse extends APIResponse {
    data?: { name: string, role: { id: string, name: string, permissions: string[] } }[]
}

interface OnGetRolesResponse extends APIResponse {
    data?: { name: string, permissions: string[],id:string }[]
}

interface OnPostRolesResponse extends APIResponse {
    data?: { name: string, permissions: string[], id: string }[]
}