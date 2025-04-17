// https://vike.dev/data

import type { PageContextServer } from "vike/types";
import { ApiResponse, apiUrl, useAPI } from "../../../../api/api";
import axios from "axios";
import { redirect, render } from "vike/abort";
export type Data = Awaited<ReturnType<typeof data>>;

export const data = async (pageContext: PageContextServer) => {
  const url = apiUrl(`/v1/connections/${pageContext.routeParams.id}/redirects`)
 const res = await useAPI().axios.get(url)
  if (res.status == axios.HttpStatusCode.Unauthorized) {
    throw redirect("/auth/signin")
  }
  else if (res.status != 200) {
    throw render(500, res.statusText)
  }
  const data = res.data as ApiResponse<{ from: string, to: string, recordType: string, localZone: boolean }[]>
  return { nodeId: pageContext.routeParams.id, redirects: data.data??[] };
};
