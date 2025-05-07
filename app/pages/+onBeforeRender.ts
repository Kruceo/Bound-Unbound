import { redirect } from "vike/abort";
import { useAPI } from "../api/api";

const f = async (f: any) => {
    console.log("using 'onBeforeRender.ts'")
    try {
        // const pg = usePageContext()
        const api = useAPI()
        const res = await api.onAuthToken()
        if (!res.ok) {
            throw new Error("problems with authentication")
        }
        const rawPerms = res.cookies?.get("permissions") ?? ""
        const permissions: string[] = rawPerms.split(",")
        const username: string | undefined = res.cookies?.get("username")
        return {
            pageContext: {
                data: {
                    ...(f.data as Record<string, string>),
                    user: {
                        username,
                        permissions,
                        authenticated: true
                    }
                }
            }
        }
    } catch (error: any) {
        console.error("error:", error.message)
    }
    if (f.urlOriginal != "/auth/signin")
        throw redirect("/auth/signin")

    // this is just to help in typescript types
    return {
        pageContext: {
            data: {
                ...f.data as Record<string, any>,
                user: {
                    permissions: [] as string[],
                    username: undefined,
                    authenticated: false
                }
            }
        }
    }
}

export default f