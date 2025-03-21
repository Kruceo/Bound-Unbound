import { apiAxios, apiUrl } from "../../../api/api"

export async function onResetAccount(sc: string) {
  const url = apiUrl("/auth/reset")
  
  try {
    await apiAxios.post(url, {
     SecretCode:sc
    })

    return true
  } catch (error: any) {
    console.log(error)
    return false
  }


}
