// https://vike.dev/data

import { redirect } from "vike/abort";
import type { PageContextServer } from "vike/types";
import { apiUrl } from "../../../../api/api";
export type Data = Awaited<ReturnType<typeof data>>;

export const data = async (pageContext: PageContextServer) => {
  const res = await fetch(apiUrl("/connections/" + pageContext.routeParams.id + "/redirects"))
  const data = await res.json() as {
    Message: string,
    Data: { From: string, To: string, RecordType: string, LocalZone: boolean }[]
  }
  return { nodeId: pageContext.routeParams.id, redirects: data.Data };
};
