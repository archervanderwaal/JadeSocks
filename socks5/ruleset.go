package socks5

type RuleSet interface {
	Allow(req *Request) bool
}

func PermitAll() RuleSet {
	return &PermitCommand{true, true, true}
}

func PermitNone() RuleSet {
	return &PermitCommand{false, false, false}
}

type PermitCommand struct {
	EnableConnect bool
	EnableBind bool
	EnableAssociate bool
}

func (p *PermitCommand) Allow(req *Request) bool {
	switch req.Command {
	case connectCommand:
		return p.EnableConnect
	case bindCommand:
		return p.EnableBind
	case associateCommand:
		return p.EnableAssociate
	default:
		return false
	}
}