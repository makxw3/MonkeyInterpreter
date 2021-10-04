package ast

import (
	"bytes"
	"monkey/token"
	"strings"
)

type Node interface {
	Token_Literal() string // Returns the value of the associated token.Literal
	Node_String() string   // A function that prints out the Node as a string
}

type Statement interface {
	Node
	Statement_Node()
}

type Expression interface {
	Node
	Expression_Node()
}

// The entire program is an array of Statements
type Program struct {
	Statements []Statement
}

func (p *Program) Node_String() string {
	var out bytes.Buffer
	for _, stmt := range p.Statements {
		out.WriteString(stmt.Node_String())
		out.WriteString("\n")
	}
	return out.String()
}

func (p *Program) Token_Literal() string {
	if len(p.Statements) > 0 {
		return p.Statements[0].Token_Literal()
	}
	return ""
}

/** The Identifier **/
type Identifier struct {
	Token token.Token // the token.IDENTIFIER token
	Value string
}

func (id *Identifier) Expression_Node()      {}
func (id *Identifier) Token_Literal() string { return id.Token.Literal }
func (id *Identifier) Node_String() string {
	var out bytes.Buffer
	out.WriteString(id.Token_Literal())
	return out.String()
}

/** The LET Statement **/
type LET_Statement struct {
	Token token.Token // The 'LET' token
	Name  *Identifier // The Identifier
	Value Expression
}

func (ls *LET_Statement) Statement_Node()       {}
func (ls *LET_Statement) Token_Literal() string { return ls.Token.Literal }
func (ls *LET_Statement) Node_String() string {
	var out bytes.Buffer
	out.WriteString(ls.Token_Literal() + " " + ls.Name.Node_String() + " = ")
	if ls.Value != nil {
		out.WriteString(ls.Value.Node_String())
	}
	out.WriteString(";")
	return out.String()
}

/** The RETURN Statement **/
type RETURN_Statement struct {
	Token       token.Token // The token.RETURN token
	ReturnValue Expression  // The expression that is being returned
}

func (rs *RETURN_Statement) Statement_Node()       {}
func (rs *RETURN_Statement) Token_Literal() string { return rs.Token.Literal }
func (rs *RETURN_Statement) Node_String() string {
	var out bytes.Buffer
	out.WriteString(rs.Token_Literal() + " ")
	if rs.ReturnValue != nil {
		out.WriteString(rs.ReturnValue.Node_String())
	}
	out.WriteString(";")
	return out.String()
}

/** The EXPRESSION_Statement **/
type EXPRESSION_Statement struct {
	Token      token.Token // The first token of the expression
	Expression Expression  // The expression
}

func (es *EXPRESSION_Statement) Statement_Node()       {}
func (es *EXPRESSION_Statement) Token_Literal() string { return es.Token.Literal }
func (es *EXPRESSION_Statement) Node_String() string {
	if es.Expression != nil {
		return es.Expression.Node_String()
	}
	return ""
}

/** INTEGER Literal **/
type INTEGER_Literal struct {
	Token token.Token // the token.INT token
	Value int64
}

func (il *INTEGER_Literal) Expression_Node()      {}
func (il *INTEGER_Literal) Token_Literal() string { return il.Token.Literal }
func (il *INTEGER_Literal) Node_String() string {
	return il.Token_Literal()
}

/** PREFIX Expression **/
type PREFIX_Expression struct {
	Token token.Token // the prefix token, e.g. !
	Right Expression
}

func (pe *PREFIX_Expression) Expression_Node()      {}
func (pe *PREFIX_Expression) Token_Literal() string { return pe.Token.Literal }
func (pe *PREFIX_Expression) Node_String() string {
	var out bytes.Buffer
	out.WriteString("(")
	out.WriteString(pe.Token_Literal())
	out.WriteString(pe.Right.Node_String())
	out.WriteString(")")
	return out.String()
}

/** Infix Expression **/
type INFIX_Expression struct {
	Token token.Token
	Left  Expression
	Right Expression
}

func (ie *INFIX_Expression) Expression_Node()      {}
func (ie *INFIX_Expression) Token_Literal() string { return ie.Token.Literal }
func (ie *INFIX_Expression) Node_String() string {
	var out bytes.Buffer
	out.WriteString("(")
	out.WriteString(ie.Left.Node_String())
	out.WriteString(ie.Token_Literal())
	out.WriteString(ie.Right.Node_String())
	out.WriteString(")")
	return out.String()
}

/** Boolean Literals **/
type Boolean struct {
	Token token.Token
	Value bool
}

func (bl *Boolean) Expression_Node()      {}
func (bl *Boolean) Token_Literal() string { return bl.Token.Literal }
func (bl *Boolean) Node_String() string   { return bl.Token.Literal }

/** IF Expressions **/
type IF_Expression struct {
	Token       token.Token // The 'if' token
	Condition   Expression
	Consequence *BlockStatement
	Alternative *BlockStatement
}

func (ie *IF_Expression) Expression_Node()      {}
func (ie *IF_Expression) Token_Literal() string { return ie.Token.Literal }
func (ie *IF_Expression) Node_String() string {
	var out bytes.Buffer
	out.WriteString(ie.Token_Literal() + " " + ie.Condition.Node_String() + " ")
	out.WriteString(ie.Consequence.Node_String())
	if ie.Alternative != nil {
		out.WriteString("else")
		out.WriteString(ie.Alternative.Node_String())
	}
	return out.String()
}

/** Block Statement **/
type BlockStatement struct {
	Token     token.Token // the '{' token
	Statemens []Statement
}

func (bs *BlockStatement) Statement_Node()       {}
func (bs *BlockStatement) Token_Literal() string { return bs.Token.Literal }
func (bs *BlockStatement) Node_String() string {
	var out bytes.Buffer
	for _, stmt := range bs.Statemens {
		out.WriteString(stmt.Node_String())
	}
	return out.String()
}

// TODO: Finsish up the function Literal
/** Function Literals **/
type FunctionLiteral struct {
	Token      token.Token // the 'function' token
	Parameters []*Identifier
	Body       *BlockStatement
}

func (fl *FunctionLiteral) Expression_Node()      {}
func (fl *FunctionLiteral) Token_Literal() string { return fl.Token.Literal }
func (fl *FunctionLiteral) Node_String() string {
	var out bytes.Buffer
	params := []string{}
	for _, p := range fl.Parameters {
		params = append(params, p.Node_String())
	}
	out.WriteString(fl.Token_Literal())
	out.WriteString("(")
	out.WriteString(strings.Join(params, ","))
	out.WriteString(")")
	out.WriteString(fl.Body.Node_String())
	return out.String()
}

/** CALL Expressions **/
type CALL_Expression struct {
	Token     token.Token // the '(' token
	Function  Expression
	Arguments []Expression
}

func (ce *CALL_Expression) Expression_Node()      {}
func (ce *CALL_Expression) Token_Literal() string { return ce.Token.Literal }
func (ce *CALL_Expression) Node_String() string {
	var out bytes.Buffer
	args := []string{}
	for _, a := range ce.Arguments {
		args = append(args, a.Node_String())
	}
	out.WriteString(ce.Function.Node_String())
	out.WriteString("(")
	out.WriteString(strings.Join(args, ","))
	out.WriteString(")")
	return out.String()
}
