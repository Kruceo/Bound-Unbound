import { useData } from "vike-react/useData";
import type { Data } from "./+data.js";
import "./Page.less"
import { useEffect, useState } from "react";
import { navigate, reload } from "vike/client/router";
import { onDeleteRedirectAction, onGetConfigHash, onNewRedirectAction, onReloadActions } from '../../actions.telefunc.js'
import Ico from "../../../../components/Ico.jsx";
import Input, { Select } from "../../../../components/Input.jsx";
import { inputPatternFor, RecordTypes } from "../../../utils.js";
import ControlsReloadButton from "../../../../components/ControlsReloadButton.jsx";
import Table from "../../../../components/Table.jsx";
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
            <button aria-label="Delete" data-balloon-pos="down" className="delete" onClick={async () => { await onDeleteRedirectAction(data.nodeId, selected); setSelected([]); reload() }}>
              <Ico>delete</Ico>
            </button> : null}
          <button aria-label="Add" data-balloon-pos="down" className="add" onClick={() => setDynamicComponent(<AddAddressForm
            onSubmit={() => { setDynamicComponent(<></>); reload() }}
            onCancel={() => setDynamicComponent(<></>)} />)}>
            <Ico>add_box</Ico>
          </button>
          <ControlsReloadButton nodeId={data.nodeId} updateIfItChanges={data} />
        </div>
      </div>

      {/* {data.blockedNames.length} */}
      <Table
        select={{ selected, setSelected, uniqueKey: "From" }}
        data={data.redirects.map(e => ({
          ...e,
          buttons: () => <a target="_blank" href={`http://${e.From}`} className="table-button"><Ico>language</Ico></a>
        }))}
        headers={[
          { acessor: "From", name: "From", width: 3 },
          { acessor: "RecordType", name: "Type", width: 3 },
          { acessor: "To", name: "To", width: 3 },
          { acessor: 'buttons', name: "" }]}>
      </Table>
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
    <h2>Redirect</h2>
    <p>This will redirect all requests from address to other address.</p>
    <div className="inputs">
      <Input title="From" required type="text" name="from" pattern={inputPatternFor("CNAME")} placeholder="www.domain.com" />
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
      <button type="submit">Finish</button>
    </div>
  </form>
}

