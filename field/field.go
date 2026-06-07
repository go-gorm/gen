package field

import "database/sql/driver"

// Field a standard field struct
type Field = genericsField[driver.Valuer]
