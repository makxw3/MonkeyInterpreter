package ast

import (
	"fmt"
	"monkey/token"
)

type Node interface {
	Token_Literal() string // Returns the value of the associated token.Literal
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

/** The LET Statement **/
type LET_Statement struct {
	Token token.Token // The 'LET' token
	Name  *Identifier // The Identifier
	Value *Expression
}

func (ls *LET_Statement) Statement_Node()       {}
func (ls *LET_Statement) Token_Literal() string { return ls.Token.Literal }

/** The RETURN Statement **/
type RETURN_Statement struct {
	Token       token.Token // The token.RETURN token
	ReturnValue *Expression // The expression that is being returned
}

func (rs *RETURN_Statement) Statement_Node()       {}
func (rs *RETURN_Statement) Token_Literal() string { return rs.Token.Literal }

/** Helper fucntion to print all the Statements **/
func Print(st Statement) {
	switch st.(type) {
	case *LET_Statement:
		letStmt := st.(*LET_Statement)
		fmt.Printf("Token_Type: %s\nName: %s\n\n", letStmt.Token.Type, letStmt.Name.Value)
	case *RETURN_Statement:
		retStmt := st.(*RETURN_Statement)
		fmt.Printf("Token_Type: %s\n\n", retStmt.Token.Type)
	}
}
