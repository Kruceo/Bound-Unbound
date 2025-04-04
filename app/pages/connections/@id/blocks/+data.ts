// https://vike.dev/data

import type { PageContextServer } from "vike/types";
import { apiAxios, ApiResponse, apiUrl } from "../../../../api/api";
import { redirect, render } from "vike/abort";
import axios from "axios";
export type Data = Awaited<ReturnType<typeof data>>;

export const data = async (pageContext: PageContextServer) => {
  const url = apiUrl(`v1/connections/${pageContext.routeParams.id}/blocks`)
  const res = await apiAxios.get(url)

  if (res.status == axios.HttpStatusCode.Unauthorized) {
    throw redirect("/auth/signin")
  }
  else if (res.status != 200) {
    throw render(500, res.statusText)
  }

  const data = res.data as ApiResponse<{ names: string[] }>
  return { nodeId: pageContext.routeParams.id, blockedNames: data.data?.names ?? [] as string[] };
};
