package builder

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/rafaeldepontes/comp/ast"
)

type Builder struct {
	Root []ast.Stmt
}

func newBuilder(body []ast.Stmt) *Builder {
	return &Builder{
		Root: body,
	}
}

var OutputFilePaths = []string{
	"output.asm",
}

func (b Builder) tokensToAsm() string {
	var sb strings.Builder

	sb.WriteString("global _start\n_start:\n")
	sb.WriteString("    mov rax, 60\n")
	fmt.Fprintf(&sb, "    mov rdi, %s\n", b.Root[0].(ast.ExitStmt).Expr.(ast.NodeExpr).Value)
	sb.WriteString("    syscall\n")

	return sb.String()
}

func Compile(root []ast.Stmt, i int) {
	b := newBuilder(root)

	content := b.tokensToAsm()
	f, err := os.OpenFile("./bin/"+OutputFilePaths[i], os.O_RDWR|os.O_APPEND|os.O_TRUNC|os.O_CREATE, 0644)
	if err != nil {
		panic(fmt.Errorf("[ERROR] could not create the file %s because %v", OutputFilePaths[i], err))
	}

	_, err = f.WriteString(content)
	if err != nil {
		panic(fmt.Errorf("[ERROR] writing content failed because %v", err))
	}

	output := strings.TrimSuffix(OutputFilePaths[i], ".asm")
	objFile := output + ".o"

	nasm := exec.Command("nasm", "-felf64", "./bin/"+OutputFilePaths[i])
	nasm.Stdout, nasm.Stderr = os.Stdout, os.Stderr
	if err := nasm.Run(); err != nil {
		panic(err)
	}

	ld := exec.Command("ld", "./bin/"+objFile, "-o", "./bin/"+output)
	ld.Stdout, ld.Stderr = os.Stdout, os.Stderr
	if err := ld.Run(); err != nil {
		panic(err)
	}
}
