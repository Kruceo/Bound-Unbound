// https://vike.dev/data

import { PageContext } from "vike/types";
import { redirect, render } from 'vike/abort'
import { apiAxios, apiUrl } from "../../../../api/api.js";
import axios, { HttpStatusCode } from "axios";

export type Data = Awaited<ReturnType<typeof data>>;

export const data = async (pg: PageContext) => {
    console.log(pg.routeParams)
    return {routeId:pg.routeParams.routeId}
};