

export type RecordTypes = "A" | "AAAA" | "TXT" | "CNAME" | "MX"

export function inputPatternFor(type: RecordTypes) {
    let result = undefined
    switch (type) {
        case "A":
            result = /^\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3}$/.source
            break
        case "AAAA":
            result = /^(([0-9a-fA-F]{1,4}:){7,7}[0-9a-fA-F]{1,4}|([0-9a-fA-F]{1,4}:){1,7}:|([0-9a-fA-F]{1,4}:){1,6}:[0-9a-fA-F]{1,4}|([0-9a-fA-F]{1,4}:){1,5}(:[0-9a-fA-F]{1,4}){1,2}|([0-9a-fA-F]{1,4}:){1,4}(:[0-9a-fA-F]{1,4}){1,3}|([0-9a-fA-F]{1,4}:){1,3}(:[0-9a-fA-F]{1,4}){1,4}|([0-9a-fA-F]{1,4}:){1,2}(:[0-9a-fA-F]{1,4}){1,5}|[0-9a-fA-F]{1,4}:((:[0-9a-fA-F]{1,4}){1,6})|:((:[0-9a-fA-F]{1,4}){1,7}|:)|fe80:(:[0-9a-fA-F]{0,4}){0,4}%[0-9a-zA-Z]{1,}|::(ffff(:0{1,4}){0,1}:){0,1}((25[0-5]|2[0-4][0-9]|1[0-9]{2}|[1-9]?[0-9])\.){3,3}(25[0-5]|2[0-4][0-9]|1[0-9]{2}|[1-9]?[0-9])|([0-9a-fA-F]{1,4}:){1,4}:((25[0-5]|2[0-4][0-9]|1[0-9]{2}|[1-9]?[0-9])\.){3,3}(25[0-5]|2[0-4][0-9]|1[0-9]{2}|[1-9]?[0-9]))$/.source
            break
        case "CNAME":
            result = /^\w+(\.\w+?)+$/.source
            break
        case "MX":
            result = /^\s*(\d+)\s+([a-zA-Z0-9.-]+\.[a-zA-Z]{2,})\s*$/.source
            break
        case "TXT":
            result = undefined
            break
        default:
            result = undefined
            break
    }

    return result//?.replace(/\\/g,"\\\\")
}
