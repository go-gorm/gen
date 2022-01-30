package postgres

const (
	columnQuery = `
SELECT
	c.table_name AS "TABLE_NAME",
	c.column_name AS "COLUMN_NAME",
	c.is_nullable AS "IS_NULLABLE",
	c.udt_name AS "COLUMN_TYPE",
	pd.description AS "COLUMN_COMMENT"
FROM
	pg_catalog.pg_statio_all_tables psat
INNER JOIN information_schema."columns" c ON
	c.table_schema = psat.schemaname
	AND c.table_name = psat.relname
LEFT JOIN pg_catalog.pg_description pd ON
	pd.objoid = psat.relid
	AND pd.objsubid = c.ordinal_position
WHERE
	psat.schemaname = ? AND psat.relname = ?
ORDER BY c.ordinal_position`
)
