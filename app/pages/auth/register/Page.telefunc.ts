import { apiAxios, apiUrl } from "../../../api/api"

const encoder = new TextEncoder();

export async function onRegisterAction(username: string, password: string) {
  const url = apiUrl("/auth/register")
  
  const data = encoder.encode(password);
  const hashBuffer = await crypto.subtle.digest('SHA-256', data);
  const hashArray = Array.from(new Uint8Array(hashBuffer));
  const hashedPassword = hashArray.map(byte => byte.toString(16).padStart(2, '0')).join('');
  
  try {
    const res = await apiAxios.post(url, {
      User: username,
      Password: hashedPassword
    })

    return res.data as { Data: { SecretCode: string }, Error: boolean, Message: string, ErrorCode: string }
  } catch (error: any) {
    console.log(error)
    return error.response.data as { Message: string, ErrorCode: "AUTH", Error: boolean, Data: undefined }
  }


}
