import axios from 'axios'
import dotenv from 'dotenv'
import { PageContext } from 'vike/types'
import * as endpoints from './endpoints'
import { getPageContext } from 'vike/getPageContext'
import { getContext } from 'telefunc'
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
    console.log(pathRoute)
    return `${proto}://${address}:${port}${(pathRoute.startsWith("/") ? "" : "/")}${pathRoute}`
}

export function useAPI() {
    let type: "telefunc" | "vike" | "" = ""
    let getFetcher = () => axios.create()

    try {
        let pageContext = getPageContext()
        if (pageContext) {
            getFetcher = () => {
                return axios.create({ headers: { "Authorization": pageContext.headers?.authorization } })
            }
            type = "vike"
        }
    } catch (err) {
        console.log(err, "deu ruim")
    }
    // get credentials using context, from telefunc
    if (type == "") {
        getFetcher = () => {
            const context = getContext() as { sessionToken: string }
            return axios.create({ headers: { "Authorization": `Bearer ${context.sessionToken}` } })
        }
    }

    const base = {
        get axios() {
            console.log("getting fetcher")
            try {
                return getFetcher()
            } catch (error) {
                console.log("getting fallback fetcher")
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