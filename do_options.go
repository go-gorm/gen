package gen

// DOOption gorm option interface
type DOOption interface {
	Apply(*DOConfig) error
	AfterInitialize(*DO) error
}

type DOConfig struct {
}

// Apply update config to new config
func (c *DOConfig) Apply(config *DOConfig) error {
	if config != c {
		*config = *c
	}
	return nil
}

// AfterInitialize initialize plugins after db connected
func (c *DOConfig) AfterInitialize(db *DO) error {
	return nil
}
