import axios from 'axios'
import dotenv from 'dotenv'
import { getContext } from 'telefunc'
import { getPageContext } from 'vike/getPageContext'

dotenv.config()

const port = process.env.API_PORT ?? "8080"
const address = process.env.API_ADDRESS ?? "127.0.0.1"
const proto = process.env.API_PROTO ?? "http"


export function apiUrl(pathRoute: string) {
    return `${proto}://${address}:${port}${(pathRoute.startsWith("/") ? "" : "/")}${pathRoute}`
}

export const apiAxios = () => {
    let sessionToken = ""
    try {
        sessionToken = "Bearer " +(getContext() as { sessionToken: string }).sessionToken
    } catch (err) {
        const pc = getPageContext()
        if (!pc || !pc.headers) throw new Error("No authorization found")
        sessionToken = pc.headers["Authorization"]
    } finally {
        // do nothing
    }

    return axios.create({ headers: { Authorization: sessionToken } })
}

export interface ApiResponse<T> {
    data?: T
    error?: boolean
    errorCode?: string
    message: string
}