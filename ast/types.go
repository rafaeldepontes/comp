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
	// Class
	Fn
)

type Type interface {
	GetType() Type_
	GetValue() any
	String() string
	Equals(other Type) bool
}

type PrimitiveType struct {
	Val  any
	Type Type_
}

func (p PrimitiveType) GetType() Type_ {
	return p.Type
}

func (p PrimitiveType) GetValue() any {
	return p.Val
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
func (a ArrayType) GetValue() any  { return nil }
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
func (n NamedType) GetValue() any  { return nil }
func (n NamedType) String() string { return n.Name }
func (n NamedType) Equals(other Type) bool {
	if other.GetType() != Struct { return false }
	if nt, ok := other.(NamedType); ok { return n.Name == nt.Name }
	if st, ok := other.(StructType); ok { return n.Name == st.Name }
	return false
}

func (n NamedType) Equals_OLD(other Type) bool {
	if other.GetType() != Struct {
		return false
	}
	if nt, ok := other.(NamedType); ok {
		return n.Name == nt.Name
	}
	return false
}

type FunctionType struct {
	Name       string
	Params     []Type
	ReturnType Type
}

func (f FunctionType) GetType() Type_ { return Fn }
func (f FunctionType) GetValue() any  { return nil }
func (f FunctionType) String() string { return f.Name }
func (f FunctionType) Equals(other Type) bool {
	fn, ok := other.(FunctionType)
	if !ok {
		return false
	}

	if len(f.Params) != len(fn.Params) {
		return false
	}

	for i := range f.Params {
		if !f.Params[i].Equals(fn.Params[i]) {
			return false
		}
	}

	return f.ReturnType.Equals(fn.ReturnType)
}

type ParamType struct {
	Name string
	Type Type_
}

func (p ParamType) GetType() Type_ {
	return p.Type
}

func (p ParamType) GetValue() any {
	return nil
}

func (p ParamType) String() string {
	switch p.Type {
	case Number:
		return "number"
	case String:
		return "string"
	case Boolean:
		return "boolean"
	case Void:
		return "void"
	default:
		return "invalid"
	}
}

func (p ParamType) Equals(other Type) bool {
	return p.Type == other.GetType()
}

type StructType struct {
	Name    string
	Fields  map[string]Type
	Methods map[string]FunctionType
}

func (n StructType) GetType() Type_ { return Struct }
func (s StructType) GetValue() any  { return nil }
func (n StructType) String() string { return n.Name }
func (n StructType) Equals(other Type) bool {
	if other.GetType() != Struct { return false }
	if nt, ok := other.(NamedType); ok { return n.Name == nt.Name }
	if st, ok := other.(StructType); ok { return n.Name == st.Name }
	return false
}

func (n StructType) Equals_OLD(other Type) bool {
	if other.GetType() != Struct {
		return false
	}
	if nt, ok := other.(NamedType); ok {
		return n.Name == nt.Name
	}
	return false
}
