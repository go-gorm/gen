package diagnostic

func DefaultMessage(code string) string {
	switch code {
	case CodeSQLIncomplete:
		return "incomplete SQL"
	case CodeSQLVar:
		return "variable parse error"
	case CodeTemplateParse:
		return "template parse error"
	case CodeSQLBuild:
		return "build SQL error"
	default:
		return ""
	}
}

func DefaultHint(code string) string {
	switch code {
	case CodeSQLIncomplete:
		return "Check for unclosed quotes or unclosed template blocks like {{...}}."
	case CodeSQLVar:
		return "Check @var / @@var usage and ensure the expression is valid."
	case CodeTemplateParse:
		return "Check template syntax inside {{...}} blocks."
	case CodeSQLBuild:
		return "Check template variables and ensure generated SQL is valid."
	default:
		return ""
	}
}
