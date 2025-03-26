import { useState } from "react"
import Ico from "../components/Ico"
import "./SideBar.less"
import { navigate } from "vike/client/router"
import logoImg from '../assets/logo.svg'
export default function SideBar() {
    const cookies = newMiniCookieParser()
    const username = cookies.get("user") ?? "User"
    return <nav id="sidebar">
        <img src={logoImg} alt="logo" />
        <div className="items">
            <a href="/"><Ico>home</Ico>Dashboard</a>
            <a href="/connections"><Ico>host</Ico>Nodes</a>
            {/* <a href="/docs"><Ico>host</Ico>Docs</a> */}
        </div>
        <UserOptionsButton user={username}></UserOptionsButton>
    </nav>
}

export function UserOptionsButton(props: { user: string }) {
    const [hide, setHide] = useState(true)
    const [hovering, setHovering] = useState(false)
    return <>
        {/* <img className="avatar" src="https://encrypted-tbn0.gstatic.com/images?q=tbn:ANd9GcSlBASaFPzRU2LryqgWoAkxM8nlhXpIsQcwKQ&s" alt="" /> */}
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

function newMiniCookieParser() {
    const cookies = document.cookie.split(/; */).reduce((acc, next) => { const [key, value] = next.split("="); acc[key] = value; return acc }, {} as Record<string, string>)
    return { get: (k: string) => { return cookies[k] as string | undefined }, ...cookies }
}