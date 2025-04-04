// https://vike.dev/data

import { PageContext } from "vike/types";
import { redirect, render } from 'vike/abort'
import { apiAxios, ApiResponse, apiUrl } from "../../../api/api.js";
import axios, { AxiosResponse, HttpStatusCode } from "axios";

export type Data = Awaited<ReturnType<typeof data>>;

export const data = async (pg: PageContext) => {
  try {
    const url = apiUrl("/v1/connections")
    const res = await apiAxios.get(url)
    if (res.status == axios.HttpStatusCode.Unauthorized) {
      throw redirect("/auth/signin")
    }
    return res.data as Promise<ApiResponse<{ name: string, remoteAddress: string }[]>>
  }
  catch (error: any) {
    if (error.status == HttpStatusCode.Unauthorized)
      throw redirect("/auth/signin")
    throw render(error.status);

  }
};