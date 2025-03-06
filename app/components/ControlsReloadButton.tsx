import { useEffect, useState } from "react"
import { onGetConfigHash, onReloadActions } from "../pages/connections/actions.telefunc"
import { navigate } from "vike/client/router"
import Ico from "./Ico"
import "./ControlsReloadButton.less";
export default function ControlsReloadButton(props: { className?: string, nodeId: string, updateIfItChanges: any }) {
    const key = props.nodeId + "-last-confighash"
    const [changed, setChanged] = useState(false)
    useEffect(() => {
        onGetConfigHash(props.nodeId).then(d => {
            const storedHash = window.localStorage.getItem(key)
            if (!storedHash) return window.localStorage.setItem(key, d.Data.Hash)
            if (storedHash != d.Data.Hash) {
                setChanged(true)
            }
        })
    }, [props.updateIfItChanges])
    return <button className={"reload-button " + (props.className ? props.className : "") + (changed ? "changed" : "")} aria-label="Reload Server" data-balloon-pos="down" onClick={() => { onReloadActions(props.nodeId); window.localStorage.removeItem(key); setChanged(false); navigate(window.location.pathname) }}><Ico>sync</Ico></button>
}