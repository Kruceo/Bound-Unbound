// https://vike.dev/data

import type { PageContextServer } from "vike/types";
export type Data = Awaited<ReturnType<typeof data>>;

export const data = async (pageContext: PageContextServer) => {
  return {routeId:pageContext.routeParams.routeId as string}
}