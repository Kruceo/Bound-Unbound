// CreateTodo.telefunc.ts
// Environment: server

// Telefunc makes onNewTodo() remotely callable
// from the browser.
export { onBlockAction,onReloadActions }
import { apiUrl } from '../../../../api/api'
// Telefunction arguments are automatically validated
// at runtime, so `text` is guaranteed to be a string.
async function onBlockAction(connectionId: string, domains: string[], action: "DELETE" | "POST") {
  const url = apiUrl("/connections/" + connectionId + "/blocked")
  const res = await fetch(url, {
    method: action, headers: { "Content-Type": "application/json" }, body: JSON.stringify({
      Names: domains
    })
  })
}

async function onReloadActions(connectionId: string) {
  const res = await fetch(apiUrl(`/connections/${connectionId}/reload`), { method: "POST" })
}