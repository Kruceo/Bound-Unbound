import axios from 'axios'
import dotenv from 'dotenv'
import * as endpoints from './endpoints'
import { getPageContext } from 'vike/getPageContext'
import { getContext } from 'telefunc'
import { PageContext } from 'vike/types'
dotenv.config()

const port = process.env.API_PORT ?? "8080"
const address = process.env.API_ADDRESS ?? "127.0.0.1"
const proto = process.env.API_PROTO ?? "http"

type BoundEndpoints<T> = {
    [K in keyof T]: T[K] extends (this: infer Ctx, ...args: infer A) => infer R
    ? (...args: A) => R
    : never;
};

export function apiUrl(pathRoute: string) {
    const f = `${proto}://${address}:${port}${(pathRoute.startsWith("/") ? "" : "/")}${pathRoute}`
    console.log("fetching",f)
    return f
}

export function useAPI(pg?: PageContext) {
    let type: "telefunc" | "vike" | "none" = "none"
    let getFetcher = () => axios.create()

    try {
        let pageContext = pg??getPageContext()
        if (pageContext) {
            getFetcher = () => {
                return axios.create({ headers: { "Authorization": pageContext.headers?.authorization } })
            }
            type = "vike"
        }
    } catch (err) {
        // console.error(err)
    }
    // get credentials using context, from telefunc
    if (type == "none") {
        getFetcher = () => {
            const context = getContext() as { sessionToken: string }
            return axios.create({ headers: { "Authorization": `Bearer ${context.sessionToken}` } })
        }
    }

    const base = {
        get axios() {
            try {
                return getFetcher()
            } catch (error) {
                return axios.create()
            }
        }
    }
    const functions = Object.fromEntries(
        Object.entries(endpoints).map(([key, func]) => [key, func.bind(base)])
    ) as BoundEndpoints<typeof endpoints>

    // console.log(functions)
    return { ...functions, ...base }
}

export interface ApiResponse<T> {
    data?: T
    error?: boolean
    errorCode?: string
    message: string
}