import "./Form.less";

export default function (props: { title: string, desc: string, children: React.ReactNode, onCancel: () => void, onSubmit: (data: FormData) => void }) {
    return <form id="styled-form" action="" onSubmit={async (e) => {
        e.preventDefault()
        const formData = new FormData(e.currentTarget)
        props.onSubmit(formData)
    }}>
        <h2>{props.title}</h2>
        <p>{props.desc}</p>
        {props.children}
        <div className="b-dock">
            <button onClick={props.onCancel} type="reset" >Cancel</button>
            <button type="submit">Finish</button>
        </div>
    </form>
}

export function FormBlock(props: { columns: number, children?: React.ReactNode }) {
    return <div className="inputs" style={{ gridTemplateColumns: `repeat(${props.columns},1fr)` }}>
        {props.children}
    </div>
}

// export function FormBottom(props: {}) {
//     return
// }