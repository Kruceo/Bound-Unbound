// https://vike.dev/data

import type { Movie, MovieDetails } from "../types.js";
import { useConfig } from "vike-react/useConfig";

export type Data = Awaited<ReturnType<typeof data>>;

export const data = async () => {

  const res = await fetch("http://localhost:8080/v1/connections")
  // console.log(await res.text())
  const data = await res.json() as { Message: string, Data: { Name: string,RemoteAddress:string }[] }
  console.log(data)
  return data;
};