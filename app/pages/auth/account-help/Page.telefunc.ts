import { apiAxios, ApiResponse, apiUrl } from "../../../api/api"

export async function onResetAccount(user: string, secretCode: string): Promise<ApiResponse<{ routeId: string }>> {
  const url = apiUrl("/auth/reset")

  try {
    const res = await apiAxios.post(url, {
      secretCode,
      user
    })

    return res.data
  } catch (error: any) {
    console.log(error)
    return error.response.data
  }


}
