import { useData } from "vike-react/useData";
import type { Data } from "./+data.js";
import "./Page.less"
import { useState } from "react";
import { navigate } from "vike/client/router";
import { onBlockAction, onGetConfigHash, onReloadActions } from '../../actions.telefunc.js'
import Ico from "../../../../components/Ico.jsx";
import ControlsReloadButton from "../../../../components/ControlsReloadButton.jsx";
import Input from "../../../../components/Input.jsx";
import { inputPatternFor } from "../../../utils.js";
import Table, { TableCell, TableRow } from "../../../../components/Table.jsx";
export default function Page() {
  const data = useData<Data>();
  const [selected, setSelected] = useState<(string)[]>([])
  const [DynamicComponent, setDynamicComponent] = useState(() => <></>)

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
            onSubmit={() => { setDynamicComponent(<></>); navigate("./blocks") }}
            onCancel={() => setDynamicComponent(<></>)} />)}>
            <Ico>add_box</Ico>
          </button>
          <ControlsReloadButton nodeId={data.nodeId} updateIfItChanges={data} />
        </div>
      </div>
      <Table select={{ selected, setSelected, uniqueKey: "domain" }} data={data.blockedNames.map(e => ({
        domain: e,
        buttons: () => <a target="_blank" href={`http://${e}`} className="table-button"><Ico>language</Ico></a>
      }))
      } headers={[{ acessor: "domain", name: "Domain", width: 5 }, { acessor: 'buttons', name: ""}]}>
      </Table>
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