package ast

type Type_ int

const (
	Invalid Type_ = iota
	Number
	String
	Boolean
	Null
	Void
	Array
	Struct
	Fn
)

// TODO: make it a real interface after the
// complex type implementation is completed.
type Type interface {
	GetType() Type_
	String() string
	Equals(other Type) bool
}

type PrimitiveType struct {
	Type Type_
}

func (p PrimitiveType) GetType() Type_ {
	return p.Type
}

func (p PrimitiveType) String() string {
	switch p.Type {
	case Number:
		return "number"
	case String:
		return "string"
	case Boolean:
		return "boolean"
	case Null:
		return "null"
	case Void:
		return "void"
	default:
		return "invalid"
	}
}

func (p PrimitiveType) Equals(other Type) bool {
	return p.Type == other.GetType()
}

type ArrayType struct {
	ElementType Type
}

func (a ArrayType) GetType() Type_ { return Array }
func (a ArrayType) String() string { return "[]" + a.ElementType.String() }
func (a ArrayType) Equals(other Type) bool {
	if other.GetType() != Array {
		return false
	}
	return a.ElementType.Equals(other.(ArrayType).ElementType)
}

type NamedType struct {
	Name string
}

func (n NamedType) GetType() Type_ { return Struct }
func (n NamedType) String() string { return n.Name }
func (n NamedType) Equals(other Type) bool {
	if other.GetType() != Struct {
		return false
	}
	if nt, ok := other.(NamedType); ok {
		return n.Name == nt.Name
	}
	return false
}
