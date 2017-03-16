package outgoing

type Command struct {
	Name     string
	Triggers []Trigger
	Father   *Command
	Children []*Command
}

func (p *Command) Commands() []string {
	node := p
	var commands []string
	for node != nil && len(node.Name) > 0 {
		commands = append(commands, node.Name)
		node = node.Father
	}

	l := len(commands)

	if l == 0 {
		return nil
	}

	var reversed = make([]string, l)
	for i := 0; i < l; i++ {
		reversed[i] = commands[l-1-i]
	}

	return reversed
}

func (p *Command) AddChild(child *Command) error {
	child.Father = p
	p.Children = append(p.Children, child)
	return nil
}

func (p *Command) Match(commands ...string) *Command {

	if len(commands) == 0 {
		return p
	}

	node := p
	level := 0

nextFor:
	for i := 0; i < len(node.Children); i++ {
		if node.Children[i].Name == commands[level] {
			node = node.Children[i]
			level++
			if level >= len(commands) {
				return node
			}
			goto nextFor
		}
	}

	return node
}
