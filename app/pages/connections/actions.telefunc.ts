import { ApiResponse, apiUrl, useAPI } from '../../api/api'

export const {onGetConfigHash,onBlockAction,onNewRedirectAction,onDeleteRedirectAction,onReloadActions,onUnblockAction} = useAPI()
export async function onGetBlocks(nodeId:string){
    const url = apiUrl(`/v1/connections/${nodeId}/redirects`)
  const res = await useAPI().axios.get(url)
  const data = res.data as ApiResponse<{ from: string, to: string, recordType: string, localZone: boolean }[]>
  return { redirects: data.data??[] };
}