import React from "react";

export default function ButtonWPlaceholder(props:React.HtmlHTMLAttributes<HTMLButtonElement>){
    return <button {...props}>
        {props.children}
        <div className="placeholder"></div>
    </button>
}