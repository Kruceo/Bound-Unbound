// CreateTodo.telefunc.ts
// Environment: server

// Telefunc makes onNewTodo() remotely callable
// from the browser.

import { apiUrl } from '../../api/api'
// Telefunction arguments are automatically validated
// at runtime, so `text` is guaranteed to be a string.
export async function onBlockAction(connectionId: string, domains: string[], action: "DELETE" | "POST") {
  const url = apiUrl("/connections/" + connectionId + "/blocked")
  const res = await fetch(url, {
    method: action, headers: { "Content-Type": "application/json" }, body: JSON.stringify({
      Names: domains
    })
  })
}

export async function onNewRedirectAction(connectionId: string, from: string, recordType: string, to: string,localZone:boolean) {
  const url = apiUrl("/connections/" + connectionId + "/redirects")
  const res = await fetch(url, {
    method: "POST", headers: { "Content-Type": "application/json" }, body: JSON.stringify({
      From: from, To: to, RecordType: recordType,LocalZone:localZone
    })
  })
}
export async function onDeleteRedirectAction(connectionId: string, domain:string) {
  const url = apiUrl("/connections/" + connectionId + "/redirects")
  const res = await fetch(url, {
    method: "DELETE", headers: { "Content-Type": "application/json" }, body: JSON.stringify({
      Domain:domain
    })
  })
}

export async function onGetConfigHash(connectionId: string) {
  const url = apiUrl("/connections/" + connectionId + "/confighash")
  const res = await fetch(url, {
    method: "GET"
  })
  return await res.json() as {Data:{Hash:string}}
}


export async function onReloadActions(connectionId: string) {
  const res = await fetch(apiUrl(`/connections/${connectionId}/reload`), { method: "POST" })
}
