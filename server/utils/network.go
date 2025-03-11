package utils

func ValidateRecordType(value string) (bool, string) {
	domainRecordTypes := []string{"CNAME", "MX"}
	// txtRecordTypes := []string{"TXT"}
	ip4RecordTypes := []string{"A"}
	ip6RecordTypes := []string{"AAAA"}
	typeExists := false
	subType := "txt"

	for _, v := range domainRecordTypes {
		if v == value && subType == "txt" {
			typeExists = true
			subType = "domain"
		} else {
			break
		}
	}
	for _, v := range ip4RecordTypes {
		if v == value && subType == "txt" {
			typeExists = true
			subType = "ip4"
		} else {
			break
		}
	}
	for _, v := range ip6RecordTypes {
		if v == value && subType == "txt" {
			typeExists = true
			subType = "ip6"
		} else {
			break
		}
	}
	return typeExists, subType
}
