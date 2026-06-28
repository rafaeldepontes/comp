package builder

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/rafaeldepontes/comp/ast"
)

type Builder struct {
	Headers    strings.Builder
	Structs    strings.Builder
	Functions  strings.Builder
	Main       strings.Builder
	InMain     bool
	StructName string
}

func NewBuilder() *Builder {
	b := &Builder{}
	b.Headers.WriteString("#include <stdio.h>\n#include <stdlib.h>\n#include <string.h>\n#include <stdbool.h>\n\n")
	return b
}

func (b *Builder) Build(program ast.BlockStmt, filename string) error {
	for _, stmt := range program.Body {
		switch stmt.(type) {
		case ast.StructStmt, ast.ImplStmt, ast.FuncStmt:
			b.genStmt(stmt)
		default:
			b.InMain = true
			b.genStmt(stmt)
			b.InMain = false
		}
	}

	cCode := b.Headers.String() + "\n" + b.Structs.String() + "\n" + b.Functions.String() + "\n" +
		"int main() {\n" + b.Main.String() + "\nreturn 0;\n}\n"

	binDir := "./bin"
	if err := os.MkdirAll(binDir, os.ModePerm); err != nil {
		return err
	}

	baseName := filepath.Base(filename)
	nameWithoutExt := strings.TrimSuffix(baseName, filepath.Ext(baseName))

	cFile := filepath.Join(binDir, nameWithoutExt+".c")
	err := os.WriteFile(cFile, []byte(cCode), 0644)
	if err != nil {
		return err
	}

	binFile := filepath.Join(binDir, nameWithoutExt)
	cmd := exec.Command("gcc", cFile, "-o", binFile)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("gcc compilation failed: %s\n%s", err, output)
	}

	fmt.Printf("[INFO] Successfully compiled to %s\n", binFile)
	return nil
}

func (b *Builder) write(s string) {
	if b.InMain {
		b.Main.WriteString(s)
	} else {
		b.Functions.WriteString(s)
	}
}

func (b *Builder) genType(t ast.Type) string {
	if t == nil {
		return "void"
	}
	switch t.GetType() {
	case ast.Number:
		return "double"

	case ast.String:
		return "char*"

	case ast.Boolean:
		return "bool"

	case ast.Void:
		return "void"

	case ast.Struct:
		return t.String() + "*"

	default:
		return "void*"
	}
}

func (b *Builder) genStmt(stmt ast.Stmt) {
	switch s := stmt.(type) {
	case ast.ExpressionStmt:
		b.write(b.genExpr(s.Expression) + ";\n")

	case ast.VarDeclStmt:
		t := b.genType(s.ExplicitType)
		if t == "void" {
			if ne, ok := s.AssignedValue.(ast.NewExpr); ok {
				t = ne.ClassName + "*"
			} else {
				t = "double"
			}
		}
		if s.AssignedValue != nil {
			b.write(fmt.Sprintf("%s %s = %s;\n", t, s.VariableName, b.genExpr(s.AssignedValue)))
		} else {
			b.write(fmt.Sprintf("%s %s;\n", t, s.VariableName))
		}

	case ast.FuncStmt:
		b.InMain = false
		b.write(fmt.Sprintf("%s %s(", b.genType(s.ReturnType), s.Name))
		for i, p := range s.Params {
			b.write(fmt.Sprintf("%s %s", b.genType(p.Type), p.Name))
			if i < len(s.Params)-1 {
				b.write(", ")
			}
		}
		b.write(") {\n")
		for _, bs := range s.Body.Body {
			b.genStmt(bs)
		}
		b.write("}\n\n")

	case ast.StructStmt:
		b.Structs.WriteString(fmt.Sprintf("typedef struct %s {\n", s.Name))
		for _, f := range s.Fields {
			b.Structs.WriteString(fmt.Sprintf("    %s %s;\n", b.genType(f.Type), f.Name))
		}
		b.Structs.WriteString(fmt.Sprintf("} %s;\n\n", s.Name))

	case ast.ImplStmt:
		b.StructName = s.Name
		for _, m := range s.Methods {
			b.InMain = false
			methodName := m.Name
			b.write(fmt.Sprintf("%s %s(%s* this", b.genType(m.ReturnType), methodName, s.Name))
			if len(m.Params) > 0 {
				b.write(", ")
			}
			for i, p := range m.Params {
				b.write(fmt.Sprintf("%s %s", b.genType(p.Type), p.Name))
				if i < len(m.Params)-1 {
					b.write(", ")
				}
			}
			b.write(") {\n")
			for _, bs := range m.Body.Body {
				b.genStmt(bs)
			}
			b.write("}\n\n")
		}
		b.StructName = ""

	case ast.ReturnStmt:
		b.write(fmt.Sprintf("return %s;\n", b.genExpr(s.ReturnValue)))

	case ast.IfStmt:
		b.write(fmt.Sprintf("if (%s) {\n", b.genExpr(s.Condition)))
		for _, bs := range s.Then.Body {
			b.genStmt(bs)
		}
		b.write("}\n")
		if s.Else != nil {
			b.write("else {\n")
			if bs, ok := s.Else.(ast.BlockStmt); ok {
				for _, stmt := range bs.Body {
					b.genStmt(stmt)
				}
			} else {
				b.genStmt(s.Else)
			}
			b.write("}\n")
		}

	case ast.WhileStmt:
		b.write(fmt.Sprintf("while (%s) {\n", b.genExpr(s.Condition)))
		for _, bs := range s.Body.Body {
			b.genStmt(bs)
		}
		b.write("}\n")

	case ast.ForStmt:
		b.write("{\n")
		if s.Init != nil {
			b.genStmt(s.Init)
		}
		cond := "true"
		if s.Cond != nil {
			cond = b.genExpr(s.Cond)
		}
		b.write(fmt.Sprintf("while (%s) {\n", cond))
		for _, bs := range s.Body.Body {
			b.genStmt(bs)
		}
		if s.Post != nil {
			b.write(fmt.Sprintf("%s;\n", b.genExpr(s.Post)))
		}
		b.write("}\n")
		b.write("}\n")

	case ast.ForEachStmt:
		b.write("{\n")
		indexName := s.Index
		if indexName == "" {
			indexName = "_idx"
		}
		iterable := b.genExpr(s.Iterable)
		b.write(fmt.Sprintf("int %s = 0;\n", indexName))
		b.write(fmt.Sprintf("while (%s < sizeof(%s)/sizeof(%s[0])) {\n", indexName, iterable, iterable))
		b.write(fmt.Sprintf("double %s = %s[%s];\n", s.Item, iterable, indexName))
		for _, bs := range s.Body.Body {
			b.genStmt(bs)
		}
		b.write(fmt.Sprintf("%s++;\n", indexName))
		b.write("}\n")
		b.write("}\n")
	}
}

func (b *Builder) genExpr(expr ast.Expr) string {
	switch e := expr.(type) {
	case ast.NumberExpr:
		return fmt.Sprintf("%f", e.Val)

	case ast.StringExpr:
		return e.Val

	case ast.BooleanExpr:
		if e.Val {
			return "true"
		}
		return "false"

	case ast.SymbolExpr:
		return e.Val

	case ast.BinaryExpr:
		return fmt.Sprintf("(%s %s %s)", b.genExpr(e.Left), e.Opr.Val, b.genExpr(e.Right))

	case ast.AssignExpr:
		return fmt.Sprintf("(%s = %s)", b.genExpr(e.Assigne), b.genExpr(e.Value))

	case ast.CallExpr:
		if me, ok := e.Callee.(ast.MemberExpr); ok {
			obj := b.genExpr(me.Object)
			method := me.Property.(ast.SymbolExpr).Val

			args := make([]string, 0)
			args = append(args, obj)
			for _, arg := range e.Args {
				args = append(args, b.genExpr(arg))
			}
			return fmt.Sprintf("%s(%s)", method, strings.Join(args, ", "))
		} else {
			callee := b.genExpr(e.Callee)
			args := make([]string, 0)
			for _, arg := range e.Args {
				args = append(args, b.genExpr(arg))
			}
			return fmt.Sprintf("%s(%s)", callee, strings.Join(args, ", "))
		}

	case ast.NewExpr:
		return fmt.Sprintf("malloc(sizeof(%s))", e.ClassName)

	case ast.MemberExpr:
		obj := b.genExpr(e.Object)
		prop := e.Property.(ast.SymbolExpr).Val
		return fmt.Sprintf("%s->%s", obj, prop)

	case ast.ThisExpr:
		return "this"

	default:
		return "0"
	}
}
