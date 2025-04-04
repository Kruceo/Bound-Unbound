import { apiAxios, ApiResponse, apiUrl } from "../../../api/api"

const encoder = new TextEncoder();

export async function onLoginAction(username: string, password: string): Promise<ApiResponse<{ token: string }>> {
  const data = encoder.encode(password);
  const hashBuffer = await crypto.subtle.digest('SHA-256', data);
  const hashArray = Array.from(new Uint8Array(hashBuffer));
  const hashedPassword = hashArray.map(byte => byte.toString(16).padStart(2, '0')).join('');

  const url = apiUrl("/auth/login")
  try {
    const res = await apiAxios.post(url, {
      User: username,
      Password: hashedPassword
    })
    console.log(res.data)
    return res.data
  } catch (error: any) {
    return error.response.data
  }


}
