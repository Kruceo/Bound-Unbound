import React from "react";
import "./Input.less";
export default function Input(props: React.DetailedHTMLProps<React.InputHTMLAttributes<HTMLInputElement>, HTMLInputElement>) {
    const { className, title, ...rest } = props
    return <div className="input-group">
        <input {...rest} className={"input " + className} />
        <label className="user-label">{props.title ?? "No title"}</label>
    </div>
}

export function Select(props: React.DetailedHTMLProps<React.SelectHTMLAttributes<HTMLSelectElement>, HTMLSelectElement>) {
    const { className, title, ...rest } = props
    return <div className="input-group">
        <select {...rest} className={"input " + className}>
            {props.children}
        </select>
        <label className="user-label">{props.title ?? "No title"}</label>
    </div>
}