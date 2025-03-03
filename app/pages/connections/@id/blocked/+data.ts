// https://vike.dev/data

import type { PageContextServer } from "vike/types";
export type Data = Awaited<ReturnType<typeof data>>;

export const data = async (pageContext: PageContextServer) => {
  const res = await fetch("http://localhost:8080/v1/connections/" + pageContext.routeParams.id + "/blocked")
  const data = await res.json() as {
    Message: string,
    Data: { Names: string[] }
  }
  console.log(data)
  return { nodeId: pageContext.routeParams.id, blockedNames: data.Data.Names.sort() as string[] };
};
