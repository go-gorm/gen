package diagnostic

const (
	// CodeSQLIncomplete indicates that SQL or a template block ended prematurely.
	CodeSQLIncomplete = "SQL_INCOMPLETE"
	// CodeSQLVar indicates that a template SQL variable could not be parsed.
	CodeSQLVar = "SQL_VAR"
	// CodeTemplateParse indicates invalid template syntax.
	CodeTemplateParse = "TEMPLATE_PARSE"
	// CodeSQLBuild indicates that parsed template sections could not produce SQL.
	CodeSQLBuild = "SQL_BUILD"
)
