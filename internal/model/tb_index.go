package model

// Index table index info
type Index struct {
	TableName  string `gorm:"column:TABLE_NAME"`
	ColumnName string `gorm:"column:COLUMN_NAME"`
	IndexName  string `gorm:"column:INDEX_NAME"`
	SeqInIndex int32  `gorm:"column:SEQ_IN_INDEX"`
	NonUnique  int32  `gorm:"column:NON_UNIQUE"`
}

func (c *Index) IsPrimaryKey() bool {
	return c != nil && c.IndexName == "PRIMARY"
}

// not primary key but unique key
func (c *Index) IsUnique() bool {
	return c != nil && !c.IsPrimaryKey() && c.NonUnique == 0
}

func GroupByColumn(indexList []*Index) map[string][]*Index {
	columnIndexMap := make(map[string][]*Index, len(indexList))
	if len(indexList) == 0 {
		return columnIndexMap
	}
	for _, idx := range indexList {
		if idx == nil {
			continue
		}
		columnIndexMap[idx.ColumnName] = append(columnIndexMap[idx.ColumnName], idx)
	}
	return columnIndexMap
}
