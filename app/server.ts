import express from 'express'
import { renderPage } from 'vike/server'
import { telefunc } from 'telefunc'
import cookieParser from 'cookie-parser'
import cors from 'cors'
import { PageContext } from 'vike/types'
const isProduction = process.env.NODE_ENV === 'production'
const root = `${__dirname}/`

startServer()

async function startServer() {
    const app = express()
    app.use(cors({credentials:true}))
    app.use(cookieParser())
    if (isProduction) {
        app.use(express.static(`${root}/dist/client`))
    } else {
        const vite = require('vite')
        const viteDevMiddleware = (
            await vite.createServer({
                root,
                server: { middlewareMode: true },
            })
        ).middlewares
        app.use(viteDevMiddleware)
    }

    app.use(express.text()) // Parse & make HTTP request body available at `req.body`
    app.all('/_telefunc', async (req, res) => {
        const context = {sessionToken:req.cookies["session"]}
        const httpResponse = await telefunc({ url: req.originalUrl, method: req.method, body: req.body, context })
        const { body, statusCode, contentType } = httpResponse
        res.status(statusCode).type(contentType).send(body)
    })

    app.get('*', async (req, res, next) => {
        const pageContextInit = {
            urlOriginal: req.originalUrl,
            headersOriginal: {...req.headers,Authorization:"Bearer " +req.cookies["session"]}

        } as PageContext
        const pageContext = await renderPage(pageContextInit)
        const { httpResponse } = pageContext
        if (!httpResponse) return next()
        const { statusCode, headers } = httpResponse
        res.status(statusCode)
        headers.forEach(([name, value]) => res.setHeader(name, value))
        httpResponse.pipe(res)
    })

    const port = process.env.PORT || 3000
    app.listen(port)
    console.log(`Server running at http://localhost:${port}`)
}