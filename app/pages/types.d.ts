// types/vike.d.ts ou global.d.ts
import {PageContextClient,PageContextServer} from "vike/types";

import onBeforeRender from './+onBeforeRender';

type OnBeforeRenderResult = Awaited<ReturnType<typeof onBeforeRender>>['pageContext']['data'];

declare module 'vike/types' {
  type PageContext<Data = OnBeforeRenderResult> = PageContextClient<Data> | PageContextServer<Data>;
}

declare module 'vike-react/useData' {
  export function useData<Data>(): Data & OnBeforeRenderResult; 
}