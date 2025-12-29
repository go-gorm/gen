package model

import (
	"strings"

	"gorm.io/gorm"
)

// Index table index info
type Index struct {
	gorm.Index
	Priority int32 `gorm:"column:SEQ_IN_INDEX"`
}

// GroupByColumn group columns
func GroupByColumn(indexList []gorm.Index) map[string][]*Index {
	columnIndexMap := make(map[string][]*Index, len(indexList))
	if len(indexList) == 0 {
		return columnIndexMap
	}

	for _, idx := range indexList {
		if idx == nil {
			continue
		}
		for i, col := range idx.Columns() {
			columnIndexMap[col] = append(columnIndexMap[col], &Index{
				Index:    idx,
				Priority: int32(i + 1),
			})
		}
	}
	return columnIndexMap
}

func sortIndexesByName(a, b *Index) int {
	if a == nil && b == nil {
		return 0
	}

	if a == nil {
		return 1
	}

	if b == nil {
		return -1
	}

	return strings.Compare(a.Name(), b.Name())
}
