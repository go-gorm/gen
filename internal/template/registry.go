package template

// ModelRegistry template for generating model registry file
const ModelRegistry = NotEditMark + `
package {{.Package}}

import (
	"reflect"
	"sync"
)

// ModelRegistry holds information about registered models in this package
type ModelRegistry struct {
	mu     sync.RWMutex
	models []interface{}
	names  []string
}

// registry is the package-level model registry
var registry = &ModelRegistry{
	models: make([]interface{}, 0),
	names:  make([]string, 0),
}

// RegisterModel registers a model to the package registry
func RegisterModel(model interface{}, name string) {
	registry.RegisterModel(model, name)
}

// RegisterModel registers a model to this registry
func (r *ModelRegistry) RegisterModel(model interface{}, name string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	
	// Check if already registered
	for i, existing := range r.models {
		if reflect.TypeOf(existing) == reflect.TypeOf(model) {
			r.names[i] = name // Update name if type already exists
			return
		}
	}
	
	r.models = append(r.models, model)
	r.names = append(r.names, name)
}

// GetAllModels returns all registered models in this package
func GetAllModels() []interface{} {
	return registry.GetAllModels()
}

// GetAllModels returns all registered models
func (r *ModelRegistry) GetAllModels() []interface{} {
	r.mu.RLock()
	defer r.mu.RUnlock()
	
	result := make([]interface{}, len(r.models))
	copy(result, r.models)
	return result
}

// GetAllModelNames returns all registered model names in this package
func GetAllModelNames() []string {
	return registry.GetAllModelNames()
}

// GetAllModelNames returns all registered model names
func (r *ModelRegistry) GetAllModelNames() []string {
	r.mu.RLock()
	defer r.mu.RUnlock()
	
	result := make([]string, len(r.names))
	copy(result, r.names)
	return result
}


`
