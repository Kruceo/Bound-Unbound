import { apiUrl } from '../../api/api'

export async function onBlockAction(connectionId: string, domains: string[], action: "DELETE" | "POST") {
  const url = apiUrl("v1/connections/" + connectionId + "/blocks")
  const res = await fetch(url, {
    method: action, headers: { "Content-Type": "application/json" }, body: JSON.stringify({
      Names: domains
    })
  })
}

export async function onNewRedirectAction(connectionId: string, from: string, recordType: string, to: string,localZone:boolean) {
  const url = apiUrl("v1/connections/" + connectionId + "/redirects")
  const res = await fetch(url, {
    method: "POST", headers: { "Content-Type": "application/json" }, body: JSON.stringify({
      From: from, To: to, RecordType: recordType,LocalZone:localZone
    })
  })
}
export async function onDeleteRedirectAction(connectionId: string, domain:string) {
  const url = apiUrl("v1/connections/" + connectionId + "/redirects")
  const res = await fetch(url, {
    method: "DELETE", headers: { "Content-Type": "application/json" }, body: JSON.stringify({
      Domain:domain
    })
  })
}

export async function onGetConfigHash(connectionId: string) {
  const url = apiUrl("v1/connections/" + connectionId + "/confighash")
  const res = await fetch(url, {
    method: "GET"
  })
  return await res.json() as {Data:{Hash:string}}
}


export async function onReloadActions(connectionId: string) {
  const res = await fetch(apiUrl(`v1/connections/${connectionId}/reload`), { method: "POST" })
}
