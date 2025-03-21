import Input from "../../../components/Input";
// import { onRegisterAction } from "./Page.telefunc";
import "./Page.less";
// import { navigate } from "vike/client/router";
import { useState } from "react";
import { onResetAccount } from "./Page.telefunc";
import { navigate } from "vike/client/router";
export default function Page() {

  const [problem, setProblem] = useState<string | undefined>(undefined)

  function inputHandler() {
    if (problem) {
      setProblem(undefined)
    }
  }
  const [submitEnabled, setSubmitEnabled] = useState(true)

  async function resetHandler(e: React.FormEvent<HTMLFormElement>){
    e.preventDefault()
    setSubmitEnabled(false)
    const formData = new FormData(e.currentTarget)

    const sc = formData.get('secret-code')?.toString()

    if (!sc){
      setSubmitEnabled(true)
      setProblem("This code is not valid")
      return 
    
    }
    if(await onResetAccount(sc)){
      navigate('/auth/register')
    }else{
      setSubmitEnabled(true)
      setProblem("This code is not valid")
    }
  }

  return (
    <>
      {/* {problem} */}
      <main className="login-frame">
        <form onSubmit={resetHandler}>
          <h2>Account Recovery</h2>
          <p>When you sign up, is displayed a recovery code, use it to reset your account.</p>
          <div className="inputs">
            <Input onInput={inputHandler} name="secret-code" placeholder="Type your code here" type="text" title="Recovery Code" required />
          </div>
          {/* {problem} */}
          <a tabIndex={0} className="forgot-password-btn" href="/auth/account-help/lost-code">Have you lost it?</a>
          <div className="bottom-container">
            <button type="submit" disabled={!submitEnabled} >
              {"Check"}
            </button>
            {
              problem ?
                <p>{problem}</p>
                : null
            }
          </div>
        </form>
      </main>
    </>
  );
}
