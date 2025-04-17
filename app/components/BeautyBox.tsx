import "./BeautyBox.less";

export default function (props: React.HTMLAttributes<HTMLDivElement>) {
    const { className, children, ...restProps } = props
    return <div className={`beauty-box ${className}`} {...restProps} >
        {children}
    </div>
}