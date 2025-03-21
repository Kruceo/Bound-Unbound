import { apiAxios, apiUrl } from "../../../api/api"

export async function onRegisterAction(username: string, password: string) {
  const url = apiUrl("/auth/register")
  
  try {
    const res = await apiAxios.post(url, {
      User: username,
      Password: password
    })

    return res.data as {Data:{SecretCode:string},Error:boolean,Message:string,ErrorCode:string}
  } catch (error: any) {
    console.log(error)
    return error.response.data as { Message: string, ErrorCode: "AUTH",Error:boolean, Data:undefined }
  }


}
