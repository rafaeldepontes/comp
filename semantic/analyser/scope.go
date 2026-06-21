package analyser

import "github.com/rafaeldepontes/comp/ast"

type Symbol struct {
	Name       string
	Type       ast.Type
	IsConstant bool
	IsGlobal   bool
}

type Scope struct {
	Parent  *Scope
	Symbols map[string]*Symbol
}

func NewScope(parent *Scope) *Scope {
	return &Scope{
		Parent:  parent,
		Symbols: make(map[string]*Symbol),
	}
}

func (s *Scope) IsGlobal() bool {
	return s.Parent == nil
}

func (s *Scope) Define(name string, sym *Symbol) bool {
	if _, has := s.Symbols[name]; has {
		return false
	}
	s.Symbols[name] = sym
	return true
}

func (s *Scope) Lookup(name string) (*Symbol, bool) {
	if val, has := s.Symbols[name]; has {
		return val, true
	}

	if s.Parent != nil {
		return s.Parent.Lookup(name)
	}

	return nil, false
}
