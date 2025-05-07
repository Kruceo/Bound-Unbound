import { useData } from "vike-react/useData";
import type { Data } from "./+data.js";
import "./Page.less"
import { useState } from "react";
import { navigate, reload } from "vike/client/router";
import { onBlockAction, onUnblockAction } from '../../actions.telefunc.js'
import Ico from "../../../../components/Ico.jsx";
import ControlsReloadButton from "../../../../components/ControlsReloadButton.jsx";
import Input from "../../../../components/Input.jsx";
import { inputPatternFor } from "../../../utils.js";
import Table from "../../../../components/Table.jsx";
import Form, { FormBlock } from "../../../../components/Form.jsx";
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
            <button aria-label="Delete" data-balloon-pos="down" className="delete" onClick={async () => { await onUnblockAction(data.nodeId, selected); setSelected([]); reload() }}>
              <Ico>delete</Ico>
            </button> : null}
          <button aria-label="Add" data-balloon-pos="down" className="add" onClick={() => setDynamicComponent(<AddAddressForm
            onSubmit={() => { setDynamicComponent(<></>) ;reload()}}
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
      } headers={[{ acessor: "domain", name: "Domain", width: 5 }, { acessor: 'buttons', name: "" }]}>
      </Table>
    </main>
  );
}


function AddAddressForm(props: { onCancel: () => void, onSubmit: () => void }) {
  const data = useData<Data>()
  return <Form 
    title="Blocking"
    desc="This will block all requests to the selected domain." onCancel={props.onCancel} onSubmit={async (formData) => {
    
    const domain = formData.get('domain')
    if (!domain) return alert("no domain");
    await onBlockAction(data.nodeId, domain.toString().split(","))
    props.onSubmit()

  }}>
    <FormBlock columns={1}>
      <Input title="Domain" required type="text" placeholder="domain.com" pattern={inputPatternFor("CNAME")} name="domain" />
    </FormBlock>
  </Form>
}