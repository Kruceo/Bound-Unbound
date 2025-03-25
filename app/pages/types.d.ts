interface APIResponse<T = any> {
    Error?: false, ErrorCode?: string, Message: string, Data?: T
}

interface ConfigHashResponse extends APIResponse {
   Data?: { Hash: string }
}