export default function Ico(props: React.HtmlHTMLAttributes<HTMLSpanElement>) {
    const {className,...rest } = props
    return <span className={"material-symbols-outlined " + (className ? className : "")} {...rest}>
        {props.children}
    </span>
}