// https://vike.dev/data

import type { PageContextServer } from "vike/types";
import { apiUrl } from "../../../../api/api";
export type Data = Awaited<ReturnType<typeof data>>;

export const data = async (pageContext: PageContextServer) => {
  const res = await fetch(apiUrl("/connections/" + pageContext.routeParams.id + "/blocked"))
  console.log(res)
  const data = await res.json() as {
    Message: string,
    Data: { Names: string[] }
  }
  console.log(data)
  return { nodeId: pageContext.routeParams.id, blockedNames: data.Data.Names.sort() as string[] };
};
