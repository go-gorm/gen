package postgres

const (
	columnQuery = `
SELECT
	c.table_name AS "TABLE_NAME",
	c.column_name AS "COLUMN_NAME",
	CASE WHEN pd.description IS NULL THEN '' ELSE pd.description END AS "COLUMN_COMMENT",
	c.data_type AS "DATA_TYPE",
	CASE WHEN tc.constraint_type = 'PRIMARY KEY' THEN 'PRI' ELSE '' END AS "COLUMN_KEY",
	c.udt_name AS "COLUMN_TYPE",
	c.is_nullable AS "IS_NULLABLE"
FROM
	pg_catalog.pg_statio_all_tables psat
INNER JOIN information_schema."columns" c ON
	c.table_schema = psat.schemaname
	AND c.table_name = psat.relname
LEFT JOIN information_schema.constraint_column_usage ccu ON
	ccu.table_schema = psat.schemaname
	AND ccu.table_name = psat.relname
	AND ccu.column_name = c.column_name
LEFT JOIN information_schema.table_constraints tc ON
	tc.constraint_schema = psat.schemaname
	AND tc.constraint_name = ccu.constraint_name
LEFT JOIN pg_catalog.pg_description pd ON
	pd.objoid = psat.relid
	AND pd.objsubid = c.ordinal_position
WHERE
	psat.relname = ?
ORDER BY c.ordinal_position`

	indexQuery = `
SELECT
	t.relname AS "TABLE_NAME",
	a.attname AS "COLUMN_NAME",
	CASE WHEN ix.indisprimary THEN 'PRIMARY' ELSE i.relname END AS "INDEX_NAME",
	CASE WHEN ix.indisunique THEN 1 ELSE 0 END AS "NON_UNIQUE"
FROM
	pg_class t,
	pg_class i,
	pg_index ix,
	pg_namespace ns,
	pg_attribute a
WHERE
	t.oid = ix.indrelid
	AND i.oid = ix.indexrelid
	AND a.attrelid = t.oid
	AND a.attnum = ANY(ix.indkey)
	AND ns."oid" = t.relnamespace
	AND t.relkind = 'r'
	AND t.relname = ?
ORDER BY
	t.relname,
	i.relname;`
)
