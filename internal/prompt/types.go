package prompt

// Info represents prompt information
type Info struct {
	Name        string
	Path        string
	Category    string
	Description string
	HasMetadata bool
	Metadata    *Metadata
}

// Variable represents a prompt variable
type Variable struct {
	Name        string
	Description string
}

// Metadata represents the YAML frontmatter of a prompt file
type Metadata struct {
	Name        string     `yaml:"name"`
	Description string     `yaml:"description"`
	Author      string     `yaml:"author"`
	Version     string     `yaml:"version"`
	Category    string     `yaml:"category"`
	Tags        []string   `yaml:"tags"`
	Variables   []Variable `yaml:"variables"`
	Extends     string     `yaml:"extends"` // Name of the template this one extends
	Path        string     `yaml:"-"`       // Path is not part of the YAML but added for reference
}
