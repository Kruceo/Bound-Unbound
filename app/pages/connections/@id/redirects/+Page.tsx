import { useData } from "vike-react/useData";
import type { Data } from "./+data.js";
import "./Page.less"
import { useEffect, useState } from "react";
import { navigate } from "vike/client/router";
import { onDeleteRedirectAction, onNewRedirectAction, onReloadActions } from '../../actions.telefunc.js'
import Ico from "../../../../components/Ico.jsx";
import Input, { Select } from "../../../../components/Input.jsx";
export default function Page() {
  const data = useData<Data>();
  const [selected, setSelected] = useState<(string)[]>([])
  const [DynamicComponent, setDynamicComponent] = useState(() => <></>)
  const [haveChanges, setHaveChanges] = useState(false)

  useEffect(() => {
    function handler(event: BeforeUnloadEvent) {
      event.preventDefault();
      event.returnValue = ""; // Necessário para exibir o alerta no Chrome
    }
    window.addEventListener("beforeunload", handler);
    return () => window.removeEventListener("beforeunload", handler)
  }, [])

  return (
    <main id="redirects-page">
      {DynamicComponent}
      <div>
        <h1 className="page-title">Redirects</h1>
        <div className="controls">
          {selected.length > 0 ?
            <button aria-label="Delete" data-balloon-pos="down" className="delete" onClick={async () => { await onDeleteRedirectAction(data.nodeId, selected[0]); setSelected([]); setHaveChanges(true); navigate("./redirects") }}>
              <Ico>delete</Ico>
            </button> : null}
          <button aria-label="Add" data-balloon-pos="down" className="add" onClick={() => setDynamicComponent(<AddAddressForm
            onSubmit={() => { setHaveChanges(true); setDynamicComponent(<></>); navigate("./blocked") }}
            onCancel={() => setDynamicComponent(<></>)} />)}>
            <Ico>add_box</Ico>
          </button>
          <button aria-label="Reload Server" data-balloon-pos="down" className={"reload "+ (haveChanges ? "has-changes" : "")} onClick={() => onReloadActions(data.nodeId)}><Ico>sync</Ico></button>
        </div>
      </div>
      <ul className="domains">
        <li className="header">
          {/* <header> */}
          <input type="checkbox" onChange={(e) => e.currentTarget.checked ? setSelected(data.redirects.map(e => e.From)) : setSelected([])} />
          <div className="table-row">
            <p>From</p>
            <p>Type</p>
            <p>To</p>
          </div>
          <div className="end"></div>
        </li>
        {/* {data.blockedNames.length} */}
        {data.redirects.map(each => <li className="domain" key={each.From}>

          <input type="checkbox"
            checked={selected.includes(each.From)}
            onChange={(e) => !e.currentTarget.checked ? setSelected(selected.filter(f => f != each.From)) : setSelected([...selected, each.From])}
          />
          {/* <input className="apple-switch" type="checkbox" defaultChecked /> */}
          <div className="table-row">
            <p>{each.From}</p>
            <p>{each.RecordType}</p>
            <p>{each.To}</p>
            {/* <p>always_nxdomain</p> */}
          </div>
          <a target="_blank" className="end" href={"http://" + each.From}>
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

    const from = formData.get('from')
    const to = formData.get('to')
    const type = formData.get('record-type')
    if (!from || !to || !type) return alert("no domain");
    await onNewRedirectAction(data.nodeId, from.toString(), type.toString(), to.toString())
    props.onSubmit()

  }}>
    <Input title="From" required type="text" name="from" placeholder="www.domain.com" />
    <div className="dock">
      <Select title="Type" required name="record-type" id="record-type">
        <option value="A">A</option>
        <option value="AAAA">AAAA</option>
        <option value="CNAME">CNAME</option>
        <option value="TXT">TEXT</option>
        <option value="MX">MX</option>
      </Select>

      <Input title="To" required type="text" placeholder="domain.com..." name="to" pattern="^(?>(\d|[1-9]\d{2}|1\d\d|2[0-4]\d|25[0-5])\.){3}(?1)$" />
    </div>
    <div className="b-dock">
      <button onClick={props.onCancel} type="reset" >Cancel</button>
      <button type="submit">Block</button>
    </div>
  </form>
}