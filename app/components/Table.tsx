import React, { PropsWithChildren, useEffect, useRef, useState } from "react"
import "./Table.less";
interface TableProps {
    headers: { name: string, acessor: string, width?: number, customHandler?: (v: any) => any }[]
    data: Record<string, any>[]
    select?: { uniqueKey: string, setSelected: (id: string[]) => void, selected: string[] }
    // children: any
}

function defaultCellHandler(v: any) {
    const type = typeof v
    if (type == "function") return v()
    if (type == "boolean") return v ? "true" : "false"
    return v
}

export default function Table(props: TableProps) {

    const ref = useRef<HTMLUListElement>(null)
    const data = props.data
    // const [data, setData] = useState(props.data)
    const detailedHeaders = props.headers.map(e => ({
        ...e,
        customHandler: e.customHandler ?? defaultCellHandler,
        type: data.length > 0 ? typeof data[0][e.acessor] : "undefined",
        width: e.width ?? 1
    }))

    const [sortState, setSortState] = useState<{ order?: boolean, header: typeof detailedHeaders[0] }>({ header: detailedHeaders[0] })

    const sortedData = sort()

    const fractions = detailedHeaders.reduce((acc, next) => acc + next.width + "fr ", "")

    function nextSortState() {
        switch (sortState.order) {
            case true:
                return false
                break;
            case false:
                return undefined
                break;

            default:
                return true
                break;
        }
    }

    function sort() {
        if (data.length == 0) return data
        const n = sortState.header
        const sortStateString = (() => {
            if (sortState.order == undefined) return 'original'
            else if (sortState.order == true) return 'asc'
            return 'desc'
        })()
        switch (n.type + ":" + sortStateString) {
            case "string:asc":
                var d = data.toSorted((a, b) => a[n.acessor].localeCompare(b[n.acessor]))
                return (d)

            case "string:desc":
                var d = data.toSorted((a, b) => b[n.acessor].localeCompare(a[n.acessor]))
                return (d)

            case "number:asc":
                var d = data.toSorted((a, b) => a[n.acessor] - (b[n.acessor]))
                return (d)

            case "number:desc":
                var d = data.toSorted((a, b) => b[n.acessor] - (a[n.acessor]))
                return (d)

            default:
                return [...data]
        }
    }

    function sortHandler(n: typeof detailedHeaders[0]) {
        if (data.length == 0) return
        setSortState({ header: n, order: nextSortState() })
    }

    function onSelectHandler(isAdd: boolean, item: any) {
        if (!props.select) throw new Error("Select section not defined in component props");
        if (isAdd)
            props.select.setSelected([...props.select.selected, item[props.select.uniqueKey]])
        else
            props.select.setSelected([...props.select.selected].filter((f => f != item[props.select?.uniqueKey ?? -1])))
    }

    function onSelectAllHandler(isAdd: boolean) {
        if (!props.select) throw new Error("Select section not defined in component props");
        if (isAdd) {
            const allIds = data.map(e => {
                if (props.select)
                    return e[props.select.uniqueKey]
                return -1
            })
            props.select.setSelected([...allIds])
        }
        else
            props.select.setSelected([])
    }

    return <ul className="table" ref={ref}>
        <li className="t-header">
            <input type="checkbox" onChange={(e) => onSelectAllHandler(e.currentTarget.checked)} />
            <div className="row" style={{ gridTemplateColumns: fractions }}>
                {
                    detailedHeaders.map(n => <TableCell key={n.acessor}
                    >
                        <span onClick={() => { sortHandler(n) }} className="title">{n.name}</span>
                        <span className="order">
                            {sortState.order == true ? "▴" : ""}
                            {sortState.order == false ? "▾" : ""}
                            {sortState.order == undefined ? "-" : ""}
                        </span>
                    </TableCell>)
                }
            </div>
        </li>

        {
            sortedData.map((v, index) => {
                return <TableRow fractions={fractions} selected={props.select?.selected.includes(v[props.select.uniqueKey])} onSelect={(isAdd) => onSelectHandler(isAdd, v)} key={v[props.select?.uniqueKey ?? index]}>
                    {
                        detailedHeaders.map(h => {
                            return <TableCell className={h.type} key={h.acessor}>
                                {h.customHandler ? h.customHandler(v[h.acessor]) : v[h.acessor]}
                            </TableCell>
                        })
                    }
                </TableRow>
            })
        }



    </ul>
}

interface TableRowProps {
    fractions?: string,
    children?: any
    selected?: boolean,
    onSelect?: (selected: boolean) => void
}


export function TableRow(props: TableRowProps) {
    return <li className="t-row">
        <input type="checkbox" checked={props.selected} onChange={(e) => {
            if (props.onSelect) props.onSelect(e.currentTarget.checked)
        }} />
        <div className="row" style={{ gridTemplateColumns: props.fractions }}>
            {props.children}
        </div>
    </li>
}

export function TableCell(props: React.HTMLAttributes<HTMLDivElement>) {
    return <div {...props} className={`t-cell ${props.className}`}>
        {props.children}
    </div>
}


// {
//     data.blockedNames.map(each => <li className="domain" key={each}>

//         <input type="checkbox"
//         // checked={selected.includes(each)}
//         // onChange={(e) => !e.currentTarget.checked ? setSelected(selected.filter(f => f != each)) : setSelected([...selected, each])}
//         />
//         {/* <input className="apple-switch" type="checkbox" defaultChecked /> */}
//         <div className="table-row">
//             <p>{each}</p>
//             {/* <p>always_nxdomain</p> */}
//         </div>
//         <a target="_blank" href={"http://" + each}>
//             <span className="material-symbols-outlined">
//                 language
//             </span>
//         </a>
//     </li>)
// }