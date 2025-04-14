import Input from "../../../../components/Input";
import { onRegisterAction } from "./Page.telefunc";
import "./Page.less";
import { navigate } from "vike/client/router";
import { useState } from "react";
import { useData } from "vike-react/useData";
import { Data } from "./+data";
export default function Page() {

  const [problem, setProblem] = useState<string | undefined>()
  const [secretCode, setSecretCode] = useState<string | undefined>()
  const data = useData() as Data
  async function registerHandler(e: React.FormEvent<HTMLFormElement>) {
    e.preventDefault()
    setSubmitEnabled(false);
    setTimeout(() => {
      setSubmitEnabled(true)
    }, 1000)

    const formData = new FormData(e.currentTarget)

    const user = formData.get('user')?.toString()
    const password = formData.get('password')?.toString()
    const password2 = formData.get('password2')?.toString()

    if (!user || !password || !password2)
      return

    if (password != password2) {
      setProblem("Passwords don't match")
      return
    }

    const re = await onRegisterAction(user.toString(), password.toString(),data.routeId)

    if (!re || re.error) {
      switch (re.errorCode) {
        case "OVERWRITING_REGISTER":
          setProblem("Isn't possible to register now")
          break;

        default:
          console.log(console.error(re))
          setProblem("Unknown problem")
          break;
      }
      return
    }
    setSecretCode(re.data?.secretCode)
    // window.localStorage.setItem("session",re.Data.Token)

  }

  function inputHandler() {
    if (problem) {
      setProblem(undefined)
    }
  }
  const [submitEnabled, setSubmitEnabled] = useState(true)

  return (
    <>
      {/* {problem} */}

      <main className="login-frame">
        {!secretCode ?
          <form onSubmit={registerHandler}>
            <h2>Create Account</h2>
            <p>Security is the first step in this journey. Please fill the form.</p>
            <div className="inputs">
              <Input onInput={inputHandler} name="user" placeholder="Type your username" type="text" title="User" required />
              <Input onInput={inputHandler} name="password" placeholder="Type your password" type="password" title="Password" required pattern=".{7}.+" />
              <Input onInput={inputHandler} name="password2" placeholder="Type your password" type="password" title="Repeat Password" required pattern=".{7}.+" />
            </div>
            {/* {problem} */}
            <div className="bottom-container">
              <button type="submit" disabled={!submitEnabled} >
                {"Sign Up"}
              </button>
              {
                problem ?
                  <p>{problem}</p>
                  : null
              }
            </div>
          </form>
          :
          <form onSubmit={(e) => { e.preventDefault(); navigate("/auth/signin") }}>
            <h2>Recovery Code</h2>
            <p>Here is your recovery code. Maybe you will need this in future, store it in a secure local.</p>
            <div className="inputs">
              <Input type="text" title="Recovery Code" placeholder="" readOnly value={secretCode} />
            </div>
            <div className="bottom-container">
              <button>Finish</button>
            </div>
          </form>
        }
      </main>
    </>
  );
}
