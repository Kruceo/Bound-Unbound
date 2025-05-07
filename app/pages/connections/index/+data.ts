// https://vike.dev/data

import { PageContext } from "vike/types";
import { useAPI } from "../../../api/api.js";

export type Data = Awaited<ReturnType<typeof data>>;

export const data = async () => {
  const api = useAPI()
  const nodes = await api.onGetNodes()
  const binds = await api.onGetBinds()
  const roles = await api.onGetRoles()
  return {
    nodes, binds, roles
  }

};