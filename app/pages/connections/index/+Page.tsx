import { useData } from "vike-react/useData";
import { data, type Data } from "./+data.js";
import { Link } from "../../../components/Link.jsx";
import "./Page.less"
import Ico from "../../../components/Ico.jsx";
import ControlsReloadButton from "../../../components/ControlsReloadButton.jsx";
import { useState } from "react";
import { onDeleteBinds, onPostBind } from "./+Page.telefunc.js"
import Input from "../../../components/Input.jsx";
import Form, { FormBlock } from "../../../components/Form.jsx";
import CheckBox from "../../../components/CheckBox.jsx";
import { reload } from "vike/client/router";
export default function Page() {
  const data = useData<Data>()
  const [DynamicComponent, setDynamicComponent] = useState(() => <></>)
  const clearDynamicComponent = () => setDynamicComponent(<></>)
  return (
    <>
      {DynamicComponent}
      <h1 className="page-title">DNS Servers</h1>
      {
        !data.nodes.data || data.nodes.data.length == 0 ?
          <p className="help-message">You don't have connected nodes. Visit <a href="/help">/help</a>.</p>
          : null
      }
      <div className="grid">

        {
          data.nodes.data?.map(each => {
            return <div key={each.name} className="connection-box">
              <h3>{each.name}</h3>
              <p>{each.remoteAddress}</p>
              <p>Unbound</p>
              <div className="bottom-bar">
                {
                  data.user.permissions.includes("manage_users") ?
                    <button aria-label="Bind Node" data-balloon-pos="down" onClick={() => setDynamicComponent(<BindForm nodeId={each.remoteAddress} onCancel={clearDynamicComponent} onSubmit={clearDynamicComponent} ></BindForm>)} >
                      <Ico>join</Ico>
                    </button>
                    : null
                }
                {
                  // data.permissions.includes("manage_nodes") ?
                    <>
                      <ControlsReloadButton nodeId={each.remoteAddress} updateIfItChanges={data} />
                      <Link aria-label="Domain Blocks" data-balloon-pos="down" className="button" href={`/connections/${each.remoteAddress}/blocks`}>
                        <Ico>block</Ico>
                      </Link>
                      <Link aria-label="Domain Redirects" data-balloon-pos="down" className="button" href={`/connections/${each.remoteAddress}/redirects`}>
                        <Ico>airline_stops</Ico>
                      </Link>
                    </>
                    // : null
                }
              </div>
            </div>
          })
        }
      </div>
    </>
  );
}



function BindForm(props: { onCancel: () => void, onSubmit: () => void, nodeId: string }) {
  const da = useData() as Data
  const oldBinds = da.binds.data?.filter(f => f.node.id == props.nodeId)

  return <Form desc="Create binds between roles and nodes." title="Bind Node" onCancel={props.onCancel} onSubmit={async (e) => {
    const newBindWithRoles = e.getAll('role')
    const toRemoveBinds = e.getAll('remove-bind')
    if (toRemoveBinds.length > 0) {
      await onDeleteBinds(toRemoveBinds.map(t => t.toString()))
    }
    if (newBindWithRoles.length > 0)
      await onPostBind(
        newBindWithRoles.map(f => ({ nodeId: props.nodeId, roleId: f.toString() }))
      )
    props.onSubmit()
    reload()
  }}>
    <FormBlock columns={1}>
      {
        oldBinds?.map(f => <CheckBox key={f.id} value={f.id} label={f.role.name} name="remove-bind" defaultChecked reverseMode />)
      }
      {
        da.roles.data?.filter((g) => !oldBinds?.map(f => f.role.id).includes(g.id))?.map(e => e.id != "0" ?
          <CheckBox key={e.id} name="role" value={e.id} label={e.name} />
          : null)
      }
    </FormBlock>
  </Form>
}