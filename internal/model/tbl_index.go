package model

import (
	"sort"
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

func sortIndexesByName(indexes []*Index) {
	sort.Slice(indexes, func(i, j int) bool {
		a, b := indexes[i], indexes[j]
		if a == nil && b == nil {
			return false
		}

		if a == nil {
			return false
		}

		if b == nil {
			return true
		}

		return strings.Compare(a.Name(), b.Name()) < 0
	})
}
