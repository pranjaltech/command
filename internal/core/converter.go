package core

// Converter translates natural language into shell commands.
type Converter struct{}

// NewConverter returns a new Converter.
func NewConverter() *Converter {
	return &Converter{}
}

// ToCommand converts a phrase into a shell command.
func (c *Converter) ToCommand(phrase string) string {
	switch phrase {
	case "list all directories":
		return "ls -d */"
	default:
		return "# unknown command"
	}
}
