import { useEffect, useState } from "react"
import { onGetConfigHash, onReloadActions } from "../pages/connections/actions.telefunc"
import { navigate } from "vike/client/router"
import Ico from "./Ico"
import "./ControlsReloadButton.less";
export default function ControlsReloadButton(props: { className?: string, nodeId: string, updateIfItChanges: any }) {
    const key = props.nodeId + "-last-confighash"
    const [changed, setChanged] = useState(false)
    const [problemWithNode, setProblemWithNode] = useState(false)
    // useEffect(() => {
    //     onGetConfigHash(props.nodeId).then(d => {
    //         const storedHash = window.localStorage.getItem(key)
    //         if (!storedHash && d.data) return window.localStorage.setItem(key, d.data.Hash)
    //         if (d.data && storedHash != d.data.Hash) return setChanged(true)
    //         if (d.error)
    //             setProblemWithNode(true)
    //     })
    // }, [props.updateIfItChanges])
    return <button disabled={problemWithNode} className={`reload-button ${(props.className ? props.className : "")} ${(changed ? "changed" : "")} ${problemWithNode ? "problem" : ""}`} aria-label="Reload Server" data-balloon-pos="down" onClick={() => { onReloadActions(props.nodeId); window.localStorage.removeItem(key); setChanged(false); navigate(window.location.pathname) }}><Ico>sync</Ico></button>
}