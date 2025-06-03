package template_test

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"gorm.io/gen/internal/template"
)

func TestCanGenerateCRUDMethodTest(t *testing.T) {
	actual := template.GenerateCRUDMethodTest()
	delimiter := "\n"
	actualSplit := strings.Split(actual, delimiter)
	expectedSplit := strings.Split(CRUDMethodTest, delimiter)
	assert.EqualValues(t, len(expectedSplit), len(actualSplit), "length of generated code does not match expected length")
	for i, actualLine := range actualSplit {
		assert.Equal(t, expectedSplit[i], actualLine, "line %d does not match expected line\nACT:[%s]\nEXP:[%s]", i+1, actualLine, expectedSplit[i])
	}
	assert.Equal(t, CRUDMethodTest, actual)
}
