// https://vike.dev/data

import { PageContext } from "vike/types";

export type Data = Awaited<ReturnType<typeof data>>;

export const data = async (pg: PageContext) => {
    return {routeId:pg.routeParams.routeId}
};