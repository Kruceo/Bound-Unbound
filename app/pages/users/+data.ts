import { PageContext } from "vike/types";
import { redirect, render } from 'vike/abort'
import { apiAxios, ApiResponse, apiUrl } from "../../api/api.js";
import axios, { AxiosResponse, HttpStatusCode } from "axios";
import { onGetRoles, onGetUsers } from "./+Page.telefunc.js";

export type Data = Awaited<ReturnType<typeof data>>;

export const data = async (pg: PageContext) => {

    const users = await onGetUsers()
    if (users.error || !users.data) throw render(500)
    const roles = await onGetRoles()
    if (roles.error || !roles.data) throw render(500)


    return { users: users.data, roles: roles.data }
};