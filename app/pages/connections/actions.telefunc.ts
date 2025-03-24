import { AxiosError } from 'axios'
import { apiAxios, apiUrl } from '../../api/api'

export async function onUnblockAction(connectionId: string, domains: string[]) {
  const url = apiUrl(`v1/connections/${connectionId}/blocks`)
  try {
    const res = await apiAxios.delete(url, {
      headers: { "Content-Type": "application/json" },
      data: { Names: domains },
    })
    return await res.data
  } catch (error) {
    return (error as AxiosError).response?.data
  }
}

export async function onBlockAction(connectionId: string, domains: string[]) {
  const url = apiUrl(`v1/connections/${connectionId}/blocks`)
  try {
    const res = await apiAxios.post(url, {
      Names: domains
    }, {
      headers: { "Content-Type": "application/json" }
    })
    return await res.data
  } catch (error) {
    return (error as AxiosError).response?.data
  }
}

export async function onNewRedirectAction(connectionId: string, from: string, recordType: string, to: string, localZone: boolean) {
  const url = apiUrl(`v1/connections/${connectionId}/redirects`)
  try {
    const res = await apiAxios.post(url, {
      From: from,
      To: to,
      RecordType: recordType,
      LocalZone: localZone,
    }, {
      headers: { "Content-Type": "application/json" },
    })
    return await res.data
  } catch (error) {
    return (error as AxiosError).response?.data
  }
}

export async function onDeleteRedirectAction(connectionId: string, domain: string[]) {
  const url = apiUrl(`v1/connections/${connectionId}/redirects`)
  let responses
  for (let dom of domain) {
    try {
      const res = await apiAxios.delete(url, {
        headers: { "Content-Type": "application/json" },
        data: { Domain: dom },
      })
      responses = res.data
    } catch (error) {
      return (error as AxiosError).response?.data
    }
  }
  return responses
}

export async function onGetConfigHash(connectionId: string) {
  const url = apiUrl(`v1/connections/${connectionId}/confighash`)
  try {
    const res = await apiAxios.get(url)
    return await res.data
  } catch (error) {
    return (error as AxiosError).response?.data
  }
}

export async function onReloadActions(connectionId: string) {
  const url = apiUrl(`v1/connections/${connectionId}/reload`)
  try {
    const res = await apiAxios.post(url)
    return await res.data
  } catch (error) {
    return (error as AxiosError).response?.data
  }
}
