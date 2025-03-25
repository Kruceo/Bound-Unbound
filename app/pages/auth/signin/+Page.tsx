import Input from "../../../components/Input";
import { onLoginAction } from "./Page.telefunc";
import "./Page.less";
import { navigate } from "vike/client/router";
import { useState } from "react";
import logoImg from '../../../assets/logo.svg'
export default function Page() {

  const [problem, setProblem] = useState<string | undefined>(undefined)

  async function loginHandler(e: React.FormEvent<HTMLFormElement>) {
    e.preventDefault()
    setSubmitEnabled(false);
    setTimeout(() => {
      setSubmitEnabled(true)
    }, 1000)

    const formData = new FormData(e.currentTarget)

    const user = formData.get('user')
    const password = formData.get('password')

    if (!user || !password)
      return

    const re = await onLoginAction(user.toString(), password.toString())

    if (!re || re.Error || !re.Data) {
      setProblem(re?.ErrorCode)
      return
    }

    document.cookie = `session=${re.Data.Token}; SameSite=none; Secure; Path=/`
    // window.localStorage.setItem("session",re.Data.Token)
    navigate("/")
  }

  function inputHandler() {
    if (problem == "AUTH") {
      setProblem(undefined)
    }
  }
const [submitEnabled,setSubmitEnabled] = useState(true)

  return (
    <>
      {/* {problem} */}
      <main className="login-frame">
        <form onSubmit={loginHandler}>
       
          <h2>Account</h2>
          <img src={logoImg} width={48} alt="logo" />
          <p>Sign in to control your DNS servers.</p>
          <div className="inputs">
            <Input onInput={inputHandler} name="user" placeholder="Type your username" type="text" title="User" required />
            <Input onInput={inputHandler} name="password" placeholder="Type your password" type="password" title="Password" required />
          </div>
          <a tabIndex={0} className="forgot-password-btn" href="/auth/account-help">Can't access your account?</a>
          <div className="bottom-container">
          <button type="submit" disabled={!submitEnabled} >
            {"Sign In"}
           
          </button>
          {
              problem == "AUTH" ?
                <p>Wrong username/password</p>
                : null
            }
          </div>
        </form>
      </main>
    </>
  );
}
