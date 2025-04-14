interface APIResponse<T = any> {
    error?: false, errorCode?: string, message: string, data?: T
}

interface ConfigHashResponse extends APIResponse {
   Data?: { Hash: string }
}

interface CreateRegisterRequestResponse extends APIResponse {
    Data?: { routeId: string }
 }