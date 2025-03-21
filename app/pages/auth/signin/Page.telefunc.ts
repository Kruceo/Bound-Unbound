import { apiAxios, apiUrl } from "../../../api/api"

export async function onLoginAction(username: string, password: string) {
  const url = apiUrl("/auth/login")
  console.log(url)
  try {
    const res = await apiAxios.post(url, {
      User: username,
      Password: password
    })

    return res.data as { Message: string, Data: { Token: string }, Error?: boolean,ErrorCode:undefined }
  } catch (error: any) {
    return error.response.data as { Message: string, ErrorCode: "AUTH",Error:boolean, Data:undefined }
    return null
  }


}
