import dotenv from 'dotenv'

dotenv.config()

const port = process.env.API_PORT ?? "8080"
const address = process.env.API_ADDRESS ?? "127.0.0.1"
const proto = process.env.API_PROTO ?? "http"


export function apiUrl(pathRoute: string) {
    return `${proto}://${address}:${port}/v1${(pathRoute.startsWith("/") ? "" : "/")}${pathRoute}`
}