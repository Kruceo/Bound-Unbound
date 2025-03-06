// https://vike.dev/data

import { apiUrl } from "../../../api/api.js";
import { onGetConfigHash } from "../actions.telefunc.js";
import type { Movie, MovieDetails } from "../types.js";
import { useConfig } from "vike-react/useConfig";

export type Data = Awaited<ReturnType<typeof data>>;

export const data = async () => {

  const res = await fetch(apiUrl("/connections"))
  // console.log(await res.text())
  const data = await res.json() as { Message: string, Data: { Name: string,RemoteAddress:string }[] }

  console.log(data)
  return data;
};