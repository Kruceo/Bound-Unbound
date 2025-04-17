// https://vike.dev/data

import { PageContext } from "vike/types";
import { redirect, render } from 'vike/abort'
import { useAPI, ApiResponse, apiUrl } from "../../../api/api.js";
import axios, { HttpStatusCode } from "axios";

export type Data = Awaited<ReturnType<typeof data>>;

export const data = async (pg: PageContext) => {
  try {
    const res = await useAPI().onGetNodes()
    console.log(res)
    return (res)
  }
  catch (error: any) {
    if (error.status == HttpStatusCode.Unauthorized)
      throw redirect("/auth/signin")
    throw render(error.status);

  }
};