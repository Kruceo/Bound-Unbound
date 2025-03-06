import { useData } from "vike-react/useData";
import type { Data } from "./+data.js";
import "./Page.less"
import { useEffect, useState } from "react";
import { navigate } from "vike/client/router";
import { onDeleteRedirectAction, onGetConfigHash, onNewRedirectAction, onReloadActions } from '../../actions.telefunc.js'
import Ico from "../../../../components/Ico.jsx";
import Input, { Select } from "../../../../components/Input.jsx";
import { inputPatternFor, RecordTypes } from "../../../utils.js";
import ControlsReloadButton from "../../../../components/ControlsReloadButton.jsx";
export default function Page() {
  const data = useData<Data>();
  const [selected, setSelected] = useState<(string)[]>([])
  const [DynamicComponent, setDynamicComponent] = useState(() => <></>)




  return (
    <main id="redirects-page">
      {DynamicComponent}
      <div>
        <h1 className="page-title">Redirects</h1>
        <div className="controls">
          {selected.length > 0 ?
            <button aria-label="Delete" data-balloon-pos="down" className="delete" onClick={async () => { await onDeleteRedirectAction(data.nodeId, selected[0]); setSelected([]); navigate("./redirects") }}>
              <Ico>delete</Ico>
            </button> : null}
          <button aria-label="Add" data-balloon-pos="down" className="add" onClick={() => setDynamicComponent(<AddAddressForm
            onSubmit={() => { setDynamicComponent(<></>); navigate(location.pathname) }}
            onCancel={() => setDynamicComponent(<></>)} />)}>
            <Ico>add_box</Ico>
          </button>
          <ControlsReloadButton nodeId={data.nodeId} updateIfItChanges={data} />
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
            <p>Local Zone</p>
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
            <p>{each.LocalZone ? "Yes" : "No"}</p>
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
  const [subtype, setSubtype] = useState<RecordTypes>("A")
  return <form className="add-form" action="" onSubmit={async (e) => {
    e.preventDefault()
    const formData = new FormData(e.currentTarget)

    const from = formData.get('from')
    const to = formData.get('to')
    const type = formData.get('record-type')
    if (!from || !to || !type) return alert("no domain");
    await onNewRedirectAction(data.nodeId, from.toString(), type.toString(), to.toString(), true)
    props.onSubmit()

  }}>
    <Input title="From" required type="text" name="from" pattern={inputPatternFor("CNAME")} placeholder="www.domain.com" />
    <div className="dock">
      <Select onChange={(e) => setSubtype(e.currentTarget.value as RecordTypes)} title="Type" required name="record-type" id="record-type">
        <option value="A">A</option>
        <option value="AAAA">AAAA</option>
        <option value="CNAME">CNAME</option>
        <option value="TXT">TEXT</option>
        <option value="MX">MX</option>
      </Select>

      <Input title="To" required type="text" placeholder="domain.com..." name="to" pattern={inputPatternFor(subtype)} />
    </div>
    <div className="b-dock">
      <button onClick={props.onCancel} type="reset" >Cancel</button>
      <button type="submit">Block</button>
    </div>
  </form>
}

