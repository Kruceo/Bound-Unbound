import { navigate, reload } from "vike/client/router";
import { onCreateRegisterRequest, onDeleteRoles, onPostRoles } from "./+Page.telefunc";
import { useContext, useState } from "react";
import Ico from "../../components/Ico";
import "./Page.less";
import Input, { Select } from "../../components/Input";
import Form, { FormBlock } from "../../components/Form";
import { useData } from "vike-react/useData";
import { Data } from "./+data";
import BeautyBox from "../../components/BeautyBox";
import Table from "../../components/Table";
import { NotificationContext } from "../../layouts/NotificationContext";
import { CopyButton } from "../../components/CopyButton";

export default function () {
    const { spawnNotification } = useContext(NotificationContext)
    const [selectedRoles, setSelectedRoles] = useState<(string)[]>([])
    const [seletedUsers, setSelectedUsers] = useState<string[]>([])
    const [DynamicComponent, setDynamicComponent] = useState(() => <></>)
    const data = useData() as Data

    async function deleteRolesHandler() {
        let res: any
        res = await onDeleteRoles(selectedRoles);

        if (res.error && res.errorCode) {
            spawnNotification(res.errorCode)
        }

        setSelectedRoles([]);
        reload()

    }

    return <>
        <main id="users-page">
            {DynamicComponent}

            <div>
                <h1 className="page-title">User Manager</h1>
                <div className="controls">
                    {
                        // selected.length > 0 ?
                        // <button aria-label="Delete" data-balloon-pos="down" className="delete" onClick={async () => { await onUnblockAction(data.nodeId, selected); setSelected([]); reload() }}>
                        // <Ico>delete</Ico>
                        // </button> : null
                    }

                    <button aria-label="Add User" data-balloon-pos="down" className="add" onClick={() => setDynamicComponent(
                        <AddUserForm
                            onSubmit={() => { setDynamicComponent(<></>); reload() }}
                            onCancel={() => setDynamicComponent(<></>)}
                        />
                    )}>
                        <Ico>add_box</Ico>
                    </button>
                </div>
            </div>
            <Table data={data.users} headers={[{ acessor: "name", name: "Name", width: 3 }, { acessor: "role", name: "Role", customHandler(v) { return v.name } }]}></Table>
            <br></br>
            <div>
                <h1 className="page-title">User Roles</h1>
                <div className="controls">
                    {
                        selectedRoles.length > 0 ?
                            <button aria-label="Delete" data-balloon-pos="down" className="delete" onClick={deleteRolesHandler}>
                                <Ico>delete</Ico>
                            </button> : null
                    }

                    <button aria-label="Add User" data-balloon-pos="down" className="add" onClick={() => setDynamicComponent(
                        <AddRoleForm
                            onSubmit={() => { setDynamicComponent(<></>); reload() }}
                            onCancel={() => setDynamicComponent(<></>)}
                        />
                    )}>
                        <Ico>add_box</Ico>
                    </button>
                </div>
            </div>
            <Table select={{ setSelected: setSelectedRoles, selected: selectedRoles, uniqueKey: "id" }} data={data.roles} headers={[
                { acessor: "name", name: "Name", width: 1 },
                { acessor: "permissions", name: "Permission", customHandler(v) { return JSON.stringify(v) } }]
            }/>
        </main>
    </>
}



function AddUserForm(props: { onCancel: () => void, onSubmit: () => void }) {

    const { spawnNotification } = useContext(NotificationContext)

    const data = useData() as Data
    const [route, setRoute] = useState<string>()
    async function handler(data: FormData) {
        const role = data.get('role')
        if (!role) return alert("no role");
        //   await onBlockAction(data.nodeId, domain.toString().split(","))
        const res = await onCreateRegisterRequest(role.toString())
        console.log(res)
        if (res.error && res.errorCode) {
            spawnNotification(res.errorCode)
            return
        }
        setRoute(res.data?.routeId)
    }

    return <>{
        !route ?
            <Form title="New User" desc="Adding new user." onCancel={props.onCancel} onSubmit={handler} >
                <FormBlock columns={1}>
                    <Select title="Role" required name="role" >
                        {/* {data.data?.map(e => <option key={e.id} value={e.id}>{e.name} {e.permissions}</option>)} */}
                    </Select>
                </FormBlock>
            </Form>
            :
            <Form title="URL" desc="Send this link to person that you want register, or open it." onCancel={props.onCancel} onSubmit={handler} >
                <FormBlock columns={1}>
                    <CopyButton textToCopy={`${window.location.protocol}//${window.location.host}/auth/register/${route}`}></CopyButton>
                    <a href={"/auth/register/" + route}>Open Register</a>
                </FormBlock>
            </Form>
    }</>
}



function AddRoleForm(props: { onCancel: () => void, onSubmit: () => void }) {

    const { spawnNotification } = useContext(NotificationContext)

    async function handler(data: FormData) {
        const name = data.get('name')
        const permssions = data.getAll('permission')
        if (!name) return alert("no name");
        const role = { name: name.toString(), permissions: permssions.map(e => e.toString()) }
        console.log(role)
        const res = await onPostRoles([role])
        console.log("dsda", res)
        if (res.error) {
            spawnNotification(res.errorCode ?? "Unknown Error")
            props.onCancel();
            return
        }
        props.onSubmit()
    }

    return <Form title="New Role" desc="Adding new role." onCancel={props.onCancel} onSubmit={handler} >
        <FormBlock columns={1}>
            <Input title="Name" name="name" required placeholder="Type a name"></Input>
        </FormBlock>
        <br />
        <FormBlock columns={1}>

            <h3>Permissions</h3>
            <div className="perms">
                <div className="permission">
                    <input type="checkbox" name="permission" value={"admin"} />
                    <span>Admin</span>
                </div>
                <div className="permission">
                    <input type="checkbox" name="permission" value={"manage_users"} />
                    <span>Manage users</span>
                </div>
                <div className="permission">
                    <input type="checkbox" name="permission" value={"manage_nodes"} />
                    <span>Manage nodes</span>
                </div>
                <div className="permission">
                    <input type="checkbox" name="permission" value={"view_all_nodes"} />
                    <span>View all nodes</span>
                </div>
                <div className="permission">
                    <input type="checkbox" name="permission" value={"t"} />
                    <span>Test</span>
                </div>
            </div>
        </FormBlock>
    </Form>
}

