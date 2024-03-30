package parser

import (
	"strings"

	"github.com/DeRuneLabs/jane/ast/models"
	"github.com/DeRuneLabs/jane/package/jn"
	"github.com/DeRuneLabs/jane/package/jnapi"
)

type function struct {
	Ast          *Func
	Desc         string
	used         bool
	checked      bool
	isEntryPoint bool
}

func (f *function) outId() string {
	if f.isEntryPoint {
		return jnapi.OutId(f.Ast.Id, nil)
	}
	return f.Ast.OutId()
}

func (f *function) getTracePointStatements() []models.Statement {
	var trace strings.Builder
	trace.WriteString(`___trace.push(`)
	var tracepoint strings.Builder
	tracepoint.WriteString(f.Ast.Id)
	tracepoint.WriteString(f.Ast.DataTypeString())
	tracepoint.WriteString("\n\t")
	tracepoint.WriteString(f.Ast.Tok.File.Path())
	trace.WriteString(jnapi.ToStr([]byte(tracepoint.String())))
	trace.WriteByte(')')
	statements := []models.Statement{{}, {}}
	statements[0].Data = models.ExprStatement{
		Expr: models.Expr{Model: exprNode{trace.String()}},
	}
	statements[1].Data = models.ExprStatement{
		Expr: models.Expr{Model: exprNode{"DEFER(___trace.ok())"}},
	}
	return statements
}

func (f function) String() string {
	var cpp strings.Builder
	cpp.WriteString(f.Head())
	cpp.WriteByte(' ')
	block := f.Ast.Block
	vars := f.Ast.RetType.Vars()
	if vars != nil {
		statements := make([]models.Statement, len(vars))
		for i, v := range vars {
			statements[i] = models.Statement{Tok: v.IdTok, Data: *v}
		}
		block.Tree = append(statements, block.Tree...)
	}
	if f.Ast.Receiver != nil && !typeIsPtr(*f.Ast.Receiver) {
		s := f.Ast.Receiver.Tag.(*jnstruct)
		self := s.selfVar(*f.Ast.Receiver)
		statements := make([]models.Statement, 1)
		statements[0] = models.Statement{Tok: s.Ast.Tok, Data: self}
		block.Tree = append(statements, block.Tree...)
	}
	block.Tree = append(f.getTracePointStatements(), block.Tree...)
	cpp.WriteString(block.String())
	return cpp.String()
}

func (f *function) Head() string {
	var cpp strings.Builder
	cpp.WriteString(f.declHead())
	cpp.WriteString(paramsToCpp(f.Ast.Params))
	return cpp.String()
}

func (f *function) declHead() string {
	var cpp strings.Builder
	cpp.WriteString(genericsToCpp(f.Ast.Generics))
	if cpp.Len() > 0 {
		cpp.WriteByte('\n')
		cpp.WriteString(models.IndentString())
	}
	cpp.WriteString(attributesToString(f.Ast.Attributes))
	cpp.WriteString(f.Ast.RetType.String())
	cpp.WriteByte(' ')
	cpp.WriteString(f.outId())
	return cpp.String()
}

func (f *function) Prototype() string {
	var cpp strings.Builder
	cpp.WriteString(f.declHead())
	cpp.WriteString(f.PrototypeParams())
	cpp.WriteByte(';')
	return cpp.String()
}

func (f *function) PrototypeParams() string {
	if len(f.Ast.Params) == 0 {
		return "(void)"
	}
	var cpp strings.Builder
	cpp.WriteByte('(')
	for _, p := range f.Ast.Params {
		cpp.WriteString(p.Prototype())
		cpp.WriteByte(',')
	}
	return cpp.String()[:cpp.Len()-1] + ")"
}

func isOutableAttribute(kind string) bool {
	return kind == jn.Attribute_Inline
}

func attributesToString(attributes []models.Attribute) string {
	var cpp strings.Builder
	for _, attr := range attributes {
		if isOutableAttribute(attr.Tag.Kind) {
			cpp.WriteString(attr.String())
			cpp.WriteByte(' ')
		}
	}
	return cpp.String()
}

func paramsToCpp(params []Param) string {
	if len(params) == 0 {
		return "(void)"
	}
	var cpp strings.Builder
	cpp.WriteByte('(')
	for _, p := range params {
		cpp.WriteString(p.String())
		cpp.WriteByte(',')
	}
	return cpp.String()[:cpp.Len()-1] + ")"
}
