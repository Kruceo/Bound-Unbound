import { useState } from "react"
import Ico from "./Ico"
import "./SideBar.less"
import { navigate } from "vike/client/router"
import logoImg from '../assets/logo.svg'
import { usePageContext } from "vike-react/usePageContext"
import { PageContext } from "vike/types"

export default function SideBar() {
    const pg:PageContext = usePageContext()
    return <nav id="sidebar">
        <img src={logoImg} alt="logo" />

        <div className="items">
            <a href="/"><Ico>home</Ico>Dashboard</a>
            <a href="/connections"><Ico>host</Ico>Nodes</a>
            {
                pg.data.user.permissions.includes("manage_users") ?
                    <a href="/users"><Ico>user</Ico>Users</a>
                    : null
            }
            {/* <a href="/docs"><Ico>host</Ico>Docs</a> */}
        </div>
        <UserOptionsButton user={pg.data.user.username ?? "..."}></UserOptionsButton>
    </nav>
}

export function UserOptionsButton(props: { user: string }) {
    const [hide, setHide] = useState(true)
    const [hovering, setHovering] = useState(false)
    return <>
        <nav className={"options " + (hide ? "hide" : "")}
            onMouseEnter={() => setHovering(true)}
            onMouseLeave={() => setHovering(false)}
            onFocus={() => setHide(false)}
            onBlur={() => {
                if (hovering == false)
                    setHide(true)
            }}
        >
            <div className="main-label">
                <button >
                    <img className="avatar" src="https://avatars.githubusercontent.com/u/98596719?v=4" alt="" />
                    <span className="username">
                        <strong>{props.user}</strong>
                    </span>
                </button>
            </div>
            <div className={"dropdown " + (hide ? "hide" : "")}>
                <button onClick={() => navigate("/help")}><span>Help</span></button>
                <button onClick={() => { document.cookie = `session=none; Path=/; SameSite=none; Secure`; window.location.reload() }}><span>Log Out</span></button>
            </div>
        </nav>
    </>
}
