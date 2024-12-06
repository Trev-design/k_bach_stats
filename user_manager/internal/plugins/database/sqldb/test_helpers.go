package sqldb

func UUIDsFromBin(ids ...*string) error {
	return uuidsFrombin(ids...)
}

func ValidValues(valid ...bool) bool {
	return validValues(valid...)
}

func MakeBatchQuery(queryStatement, payload string, length int) string {
	return makeBatchQuery(queryStatement, payload, length)
}
