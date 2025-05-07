import { useState } from "react";
import "./CheckBox.less";

export default function ({className,defaultChecked,label,name,reverseMode,value}: { label?: string, defaultChecked?: boolean, className?: string, name?: string, value?: number | string, reverseMode?: boolean }) {
    const [checked, setChecked] = useState(defaultChecked?true:false)
    return <div className={`component-checkbox ${className}`}>
        <input type="checkbox" name={"visual-" + name} defaultChecked={defaultChecked} onInput={(e) => setChecked(e.currentTarget.checked)} />
        <input style={{display:"none"}} type="checkbox" name={name} value={value} checked={reverseMode ? !checked : checked} />
        {label ? <span>{label}</span> : null}
    </div>
}