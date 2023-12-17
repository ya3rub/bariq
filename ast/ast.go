package ast

import (
	"bytes"
	"strings"

	"bariq/token"
)

type Node interface {
	// used only in debugging and testing
	TokenLiteral() string
	String() string
}
type Stmt interface {
	Node
	// a no op method to help the compiler to find if we
	// used Statement instead of Expression and vice versa
	statementNode()
}

type Expr interface {
	Node
	// a no op method to help the compiler to find if we
	// used Statement instead of Expression and vice versa
	expressionNode()
}

// the root Node
type Program struct {
	Stmts []Stmt
}

func (p *Program) String() string {
	var out bytes.Buffer
	for _, s := range p.Stmts {
		out.WriteString(s.String())
	}
	return out.String()
}

func (p *Program) TokenLiteral() string {
	if len(p.Stmts) > 0 {
		return p.Stmts[0].TokenLiteral()
	} else {
		return ""
	}
}

type ExprStmt struct {
	Token token.Token // RET
	Expr
}

func (*ExprStmt) statementNode()          {}
func (es *ExprStmt) TokenLiteral() string { return es.Token.Literal }

func (es *ExprStmt) String() string {
	if es.Expr != nil {
		return es.Expr.String()
	}
	return ""
}

type ReturnStmt struct {
	Token token.Token // RET
	Value Expr
}

func (ls *ReturnStmt) String() string {
	var out bytes.Buffer
	out.WriteString(ls.TokenLiteral() + " ")
	if ls.Value != nil {
		out.WriteString(ls.Value.String())
	}
	out.WriteString(";")
	return out.String()
}
func (*ReturnStmt) statementNode()          {}
func (rs *ReturnStmt) TokenLiteral() string { return rs.Token.Literal }

type LetStmt struct {
	Token token.Token // LET
	Name  *Ident
	Value Expr
}

func (ls *LetStmt) String() string {
	var out bytes.Buffer
	out.WriteString(ls.TokenLiteral() + " ")
	out.WriteString(ls.Name.String() + " ")
	out.WriteString("= ")
	if ls.Value != nil {
		out.WriteString(ls.Value.String())
	}
	out.WriteString(";")
	return out.String()
}
func (*LetStmt) statementNode()          {}
func (ls *LetStmt) TokenLiteral() string { return ls.Token.Literal }

type Ident struct {
	Token token.Token // IDENT
	Value string
}

func (i *Ident) expressionNode() {}
func (i *Ident) String() string {
	return i.Value
}
func (id *Ident) TokenLiteral() string { return id.Token.Literal }

type IndexExpr struct {
	Token token.Token
	Left  Expr
	Index Expr
}

func (ie *IndexExpr) expressionNode()      {}
func (ie *IndexExpr) TokenLiteral() string { return ie.Token.Literal }
func (ie *IndexExpr) String() string {
	var out bytes.Buffer
	out.WriteString("(")
	out.WriteString(ie.Left.String())
	out.WriteString("[")
	out.WriteString(ie.Index.String())
	out.WriteString("])")
	return out.String()
}

type HashLiteral struct {
	Token token.Token //{
	Pairs map[Expr]Expr
}

func (hl *HashLiteral) expressionNode()      {}
func (hl *HashLiteral) TokenLiteral() string { return hl.Token.Literal }
func (hl *HashLiteral) String() string {
	var out bytes.Buffer
	pairs := []string{}
	for k, v := range hl.Pairs {
		pairs = append(pairs, k.String()+":"+v.String())
	}
	out.WriteString("{")
	out.WriteString(strings.Join(pairs, ", "))
	out.WriteString("}")
	return out.String()
}

type ArrayLiteral struct {
	Token  token.Token
	Elmnts []Expr
}

func (al *ArrayLiteral) expressionNode()      {}
func (al *ArrayLiteral) TokenLiteral() string { return al.Token.Literal }
func (al *ArrayLiteral) String() string {
	var out bytes.Buffer
	elmnts := []string{}
	for _, el := range al.Elmnts {
		elmnts = append(elmnts, el.String())
	}
	out.WriteString("[")
	out.WriteString(strings.Join(elmnts, ", "))
	out.WriteString("]")
	return out.String()
}

type StringLiteral struct {
	Token token.Token
	Value string
}

func (s *StringLiteral) expressionNode()      {}
func (s *StringLiteral) String() string       { return s.Token.Literal }
func (s *StringLiteral) TokenLiteral() string { return s.Token.Literal }

type IntLiteral struct {
	Token token.Token
	// now you know why using value along with
	// token literal
	Value int64
}

func (i *IntLiteral) expressionNode()       {}
func (i *IntLiteral) String() string        { return i.Token.Literal }
func (id *IntLiteral) TokenLiteral() string { return id.Token.Literal }

type PrefixExpr struct {
	Token    token.Token // prefix token , ex: !
	Operator string
	Right    Expr
}

func (pe *PrefixExpr) expressionNode() {}
func (pe *PrefixExpr) String() string {
	var out bytes.Buffer
	out.WriteString("(")
	out.WriteString(pe.Operator)
	out.WriteString(pe.Right.String())
	out.WriteString(")")
	return out.String()
}
func (pe *PrefixExpr) TokenLiteral() string { return pe.Token.Literal }

type InfixExpr struct {
	Token    token.Token // prefix token , ex: !
	Left     Expr
	Operator string
	Right    Expr
}

func (ie *InfixExpr) expressionNode() {}
func (ie *InfixExpr) String() string {
	var out bytes.Buffer
	out.WriteString("(")
	out.WriteString(ie.Left.String())
	out.WriteString(" " + ie.Operator + " ")
	out.WriteString(ie.Right.String())
	out.WriteString(")")
	return out.String()
}
func (ie *InfixExpr) TokenLiteral() string { return ie.Token.Literal }

type Boolean struct {
	Token token.Token
	// now you know why using value along with
	// token literal
	Value bool
}

func (b *Boolean) expressionNode()      {}
func (b *Boolean) String() string       { return b.Token.Literal }
func (b *Boolean) TokenLiteral() string { return b.Token.Literal }

type IfExpr struct {
	Token       token.Token
	Condition   Expr
	Consequence *BlockStmt
	Alternative *BlockStmt
}

func (ie *IfExpr) expressionNode()      {}
func (ie *IfExpr) TokenLiteral() string { return ie.Token.Literal }
func (ie *IfExpr) String() string {
	var out bytes.Buffer
	out.WriteString("if")
	out.WriteString(ie.Condition.String())
	out.WriteString(" ")
	out.WriteString(ie.Consequence.String())
	if ie.Alternative != nil {
		out.WriteString("else ")
		out.WriteString(ie.Alternative.String())
	}
	return out.String()
}

type BlockStmt struct {
	Token token.Token
	Stmts []Stmt
}

func (bs *BlockStmt) expressionNode()      {}
func (bs *BlockStmt) TokenLiteral() string { return bs.Token.Literal }
func (bs *BlockStmt) String() string {
	var out bytes.Buffer
	for _, s := range bs.Stmts {
		out.WriteString(s.String())
	}
	return out.String()
}

type YieldExpr struct {
	Token token.Token // function
	Arg   Expr
}

func (ye *YieldExpr) expressionNode()      {}
func (ye *YieldExpr) TokenLiteral() string { return ye.Token.Literal }
func (ye *YieldExpr) String() string {
	var out bytes.Buffer
	out.WriteString(ye.TokenLiteral())
	out.WriteString(" ")
	out.WriteString(ye.Arg.String())
	return out.String()
}

type AwaitExpr struct {
	Token token.Token // function
	Arg   Expr
}

func (a *AwaitExpr) expressionNode()      {}
func (a *AwaitExpr) TokenLiteral() string { return a.Token.Literal }
func (a *AwaitExpr) String() string {
	var out bytes.Buffer
	out.WriteString(a.TokenLiteral())
	out.WriteString("(")
	out.WriteString(a.Arg.String())
	out.WriteString(") ")
	return out.String()
}

type FunctionLiteral struct {
	Token      token.Token // function
	Parameters []*Ident
	Body       *BlockStmt
	Async      bool
	Gen        bool
}

func (fl *FunctionLiteral) expressionNode()      {}
func (fl *FunctionLiteral) TokenLiteral() string { return fl.Token.Literal }
func (fl *FunctionLiteral) String() string {
	var out bytes.Buffer
	params := []string{}
	for _, p := range fl.Parameters {
		params = append(params, p.String())
	}
	out.WriteString(fl.TokenLiteral())
	out.WriteString("(")
	out.WriteString(strings.Join(params, ", "))
	out.WriteString(") ")
	out.WriteString(fl.Body.String())
	return out.String()
}

type CallExpr struct {
	Token    token.Token
	Function Expr // ident or function literal
	Args     []Expr
}

func (ce *CallExpr) expressionNode()      {}
func (ce *CallExpr) TokenLiteral() string { return ce.Token.Literal }
func (ce *CallExpr) String() string {
	var out bytes.Buffer
	args := []string{}
	for _, p := range ce.Args {
		args = append(args, p.String())
	}
	out.WriteString(ce.Function.String())
	out.WriteString("(")
	out.WriteString(strings.Join(args, ", "))
	out.WriteString(")")
	return out.String()
}
