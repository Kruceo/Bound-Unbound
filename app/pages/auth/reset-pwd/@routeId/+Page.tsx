import Input from "../../../../components/Input";
import "./Page.less";
import { useState } from "react";
import { onResetPassword } from "./Page.telefunc";
import { navigate } from "vike/client/router";
import { useData } from "vike-react/useData";
import { Data } from "./+data";
export default function Page() {
  const data = useData() as Data
  const [problem, setProblem] = useState<string | undefined>(undefined)

  function inputHandler() {
    if (problem) {
      setProblem(undefined)
    }
  }

  const [submitEnabled, setSubmitEnabled] = useState(true)
  const [secretCode, setSecretCode] = useState<string | undefined>()

  async function resetHandler(e: React.FormEvent<HTMLFormElement>) {
    e.preventDefault()
    setSubmitEnabled(false)
    const formData = new FormData(e.currentTarget)

    const p1 = formData.get('password')?.toString()
    const p2 = formData.get('password2')?.toString()

    if (!p1 || !p2)
      return

    if (p1 != p2) {
      setProblem("Passwords don't match")
      setSubmitEnabled(true)
      return
    }

    const re = await onResetPassword(data.routeId, p1)
    if (!re || re.error) {
      setTimeout(() => {
        switch (re.errorCode) {
          case "OVERWRITING_REGISTER":
            setProblem("Isn't possible to register now")
            setSubmitEnabled(true)
            break;

          default:
            setProblem("Unknown problem")
            setSubmitEnabled(true)
            break;
        }
      }, 2000)
      return
    }
    setSecretCode(re.data?.secretCode)
  }

  return (
    <>
      {/* {problem} */}
      <main className="login-frame">
        {
          !secretCode ?
            <form onSubmit={resetHandler}>
              <h2>Password Reset</h2>
              <p>Reseting your password.</p>
              <div className="inputs">
                <Input onInput={inputHandler} name="password" placeholder="Type your password" type="password" title="Password" required pattern=".{7}.+" />
                <Input onInput={inputHandler} name="password2" placeholder="Type your password" type="password" title="Repeat Password" required pattern=".{7}.+" />
              </div>
              <div className="bottom-container">
                <button type="submit" disabled={!submitEnabled} >
                  {"Update"}
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
