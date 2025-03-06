import { useData } from "vike-react/useData";
import type { Data } from "./+data.js";
import "./Page.less"
import { useEffect, useState } from "react";
import { navigate } from "vike/client/router";
import { onBlockAction, onGetConfigHash, onReloadActions } from '../../actions.telefunc.js'
import Ico from "../../../../components/Ico.jsx";
import ControlsReloadButton from "../../../../components/ControlsReloadButton.jsx";
import Input from "../../../../components/Input.jsx";
import { inputPatternFor } from "../../../utils.js";
export default function Page() {
  const data = useData<Data>();
  const [selected, setSelected] = useState<(string)[]>([])
  const [DynamicComponent, setDynamicComponent] = useState(() => <></>)
  const [changed, setChanged] = useState(false)

  useEffect(() => {
    onGetConfigHash(data.nodeId).then(d => {
      const key = "last-config-save-hash"
      const storedHash = window.localStorage.getItem(key)
      if (!storedHash) return window.localStorage.setItem(key, d.Data.Hash)
      if (storedHash != d.Data.Hash) {
        setChanged(true)
      }
    })
  }, [data])

  return (
    <main id="blocks-page">
      {DynamicComponent}
      <div>
        <h1 className="page-title">Blocked Domains</h1>
        <div className="controls">
          {selected.length > 0 ?
            <button aria-label="Delete" data-balloon-pos="down" className="delete" onClick={async () => { await onBlockAction(data.nodeId, selected, "DELETE"); setSelected([]); navigate("./blocked") }}>
              <Ico>delete</Ico>
            </button> : null}
          <button aria-label="Add" data-balloon-pos="down" className="add" onClick={() => setDynamicComponent(<AddAddressForm
            onSubmit={() => { setDynamicComponent(<></>); navigate("./blocked") }}
            onCancel={() => setDynamicComponent(<></>)} />)}>
            <Ico>add_box</Ico>
          </button>
          <ControlsReloadButton nodeId={data.nodeId} updateIfItChanges={data}/>
        </div>
      </div>
      <ul className="domains">
        <li className="header">
          {/* <header> */}
          <input type="checkbox" onChange={(e) => e.currentTarget.checked ? setSelected(data.blockedNames) : setSelected([])} />
          <p>Domain</p>
        </li>
        {/* {data.blockedNames.length} */}
        {data.blockedNames.map(each => <li className="domain" key={each}>

          <input type="checkbox"
            checked={selected.includes(each)}
            onChange={(e) => !e.currentTarget.checked ? setSelected(selected.filter(f => f != each)) : setSelected([...selected, each])}
          />
          {/* <input className="apple-switch" type="checkbox" defaultChecked /> */}
          <div className="table-row">
            <p>{each}</p>
            {/* <p>always_nxdomain</p> */}
          </div>
          <a target="_blank" href={"http://" + each}>
            <span className="material-symbols-outlined">
              language
            </span>
          </a>
        </li>)}


      </ul>
    </main>
  );
}


function AddAddressForm(props: { onCancel: () => void, onSubmit: () => void }) {
  const data = useData<Data>()
  return <form className="add-form" action="" onSubmit={async (e) => {
    e.preventDefault()
    const formData = new FormData(e.currentTarget)

    const domain = formData.get('domain')
    if (!domain) return alert("no domain");
    await onBlockAction(data.nodeId, domain.toString().split(","), "POST")
    props.onSubmit()

  }}>
    <Input title="Domain" required type="text" placeholder="domain.com" pattern={inputPatternFor("CNAME")} name="domain" />
    <div className="b-dock">
      <button onClick={props.onCancel} type="reset" >Cancel</button>
      <button type="submit">Block</button>
    </div>
  </form>
}