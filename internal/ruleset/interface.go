// Package ruleset provides core business logic for managing AI editor rulesets.
package ruleset

// ServiceInterface defines the interface for ruleset operations
type ServiceInterface interface {
	Create(rs *Ruleset) error
	Get(name string) (*Ruleset, error)
	Update(name string, updates *Update) error
	Upsert(rs *Ruleset, updates *Update) error
	Delete(name string) error
	List() ([]*Ruleset, error)
	Search(pattern string) ([]*Ruleset, error)
	Exists(name string) (bool, error)
	ListNames() ([]string, error)
}
