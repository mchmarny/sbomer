package target

const (
	ImportTargetTypeUndefined ImportTargetType = iota
	ImportTargetTypeBigQuery
	ImportTargetTypePostgres
	ImportTargetTypeMySQL

	ImportTargetUndefined = "undefined"
	ImportTargetBigQuery  = "bigquery"
	ImportTargetPostgres  = "postgres"
	ImportTargetMySQL     = "mysql"
)

// ImportTargetType is the type of target to import into.
type ImportTargetType int64

// String returns a string representation of the ImportTargetType.
func (t ImportTargetType) String() string {
	switch t {
	case ImportTargetTypeBigQuery:
		return ImportTargetBigQuery
	case ImportTargetTypePostgres:
		return ImportTargetPostgres
	case ImportTargetTypeMySQL:
		return ImportTargetMySQL
	default:
		return ImportTargetUndefined
	}
}

// ParseImportTargetType parses a string into an ImportTargetType.
func ParseImportTargetType(t string) ImportTargetType {
	switch t {
	case ImportTargetBigQuery:
		return ImportTargetTypeBigQuery
	case ImportTargetPostgres:
		return ImportTargetTypePostgres
	case ImportTargetMySQL:
		return ImportTargetTypeMySQL
	default:
		return ImportTargetTypeUndefined
	}
}
