package parser

import (
	"fmt"

	l "rainlang/lexer"
)

type AST struct {
	Seq Seq
}

type Seq struct {
	Stmts []Stmt
}

type Stmt interface {
	isStmt()
}

type ExprStmt struct {
	Expr Expr
}

type Expr interface {
	isExpr()
}

type BinExpr struct {
	Left  Factor
	Right Factor
	Op    BinOp
}

type AssignExpr struct {
	Target l.Token
	Value  Expr
}

type DeclExpr struct {
	Type   l.TokenType
	Assign AssignExpr
}

type BinOp struct {
	Type l.TokenType
}

type Factor interface {
	isFactor()
}

type Literal struct {
	Type  l.TokenType
	Value string
}

func (ExprStmt) isStmt() {}

func (AssignExpr) isExpr() {}
func (BinExpr) isExpr()    {}
func (DeclExpr) isExpr()   {}
func (Literal) isExpr()    {}

func (BinExpr) isFactor() {}
func (Literal) isFactor() {}

func contains(arr []byte, val byte) bool {
	for _, a := range arr {
		if a == val {
			return true
		}
	}
	return false
}

func factor(tokens []l.Token, index int) (Factor, bool, int) {
	fmt.Println("factor")
	var result Factor

	start := 'a'
	final := []byte{'b'}

	currState := start
	currIdx := index

	failure := false

	for !failure && !contains(final, byte(currState)) && currIdx < len(tokens) {
		switch currState {
		case 'a':
			switch tokens[currIdx].Type {
			case l.IntLiteral, l.FloatLiteral, l.CharLiteral, l.StringLiteral, l.BoolLiteral, l.Ident:
				currState = 'b'
				result = Literal{
					Type:  tokens[currIdx].Type,
					Value: tokens[currIdx].Lexeme,
				}
			default:
				failure = true
			}
		}
		currIdx++
	}

	return result, !failure, currIdx
}

func binOp(tokens []l.Token, index int) (BinOp, bool, int) {
	fmt.Println("binOp")
	var result BinOp

	start := 'a'
	final := []byte{'b'}

	currState := start
	currIdx := index

	failure := false
	for !failure && !contains(final, byte(currState)) && currIdx < len(tokens) {
		switch currState {
		case 'a':
			switch tokens[currIdx].Type {
			case l.Plus, l.Minus, l.Star, l.Slash:
				currState = 'b'
				result = BinOp{
					Type: tokens[currIdx].Type,
				}
			default:
				failure = true
			}
		}
		currIdx++
	}

	return result, !failure, currIdx
}

func binExpr(tokens []l.Token, index int) (BinExpr, bool, int) {
	fmt.Println("binExpr")
	var result BinExpr

	start := 'a'
	final := []byte{'z'}

	currState := start
	currIdx := index

	failure := false

	for !failure && !contains(final, byte(currState)) && currIdx < len(tokens) {
		switch currState {
		case 'a':
			if f, exists, newI := factor(tokens, currIdx); exists {
				result.Left = f
				currState = 'b'
				currIdx = newI
			} else {
				failure = true
			}
		case 'b':
			if op, exists, newI := binOp(tokens, currIdx); exists {
				result.Op = op
				currIdx = newI
				switch op.Type {
				case l.Plus, l.Minus:
					currState = 'e'
				default:
					currState = 'c'
				}
			} else {
				failure = true
			}
		case 'c':
			if f, exists, newI := factor(tokens, currIdx); exists {
				result.Right = f
				currState = 'd'
				currIdx = newI
			} else {
				failure = true
			}
		case 'd':
			if op, exists, newI := binOp(tokens, currIdx); exists {
				result = BinExpr{
					Left: result,
					Op:   op,
				}
				currIdx = newI
				switch op.Type {
				case l.Plus, l.Minus:
					currState = 'e'
				default:
					currState = 'c'
				}
			} else {
				currState = 'z'
			}
		case 'e':
			if b, exists, newI := binExpr(tokens, currIdx); exists {
				result.Right = b
				currState = 'z'
				currIdx = newI
			} else {
				currState = 'c'
			}
		}
	}

	return result, !failure, currIdx
}

func assignExpr(tokens []l.Token, index int) (AssignExpr, bool, int) {
	fmt.Println("assignExpr")
	var result AssignExpr

	start := 'a'
	final := []byte{'e'}

	currState := start
	currIdx := index

	failure := false
	for !failure && !contains(final, byte(currState)) && currIdx < len(tokens) {
		switch currState {
		case 'a':
			if tokens[currIdx].Type == l.Ident {
				currState = 'b'
				result.Target = tokens[currIdx]
				currIdx++
			} else {
				failure = true
			}
		case 'b':
			if tokens[currIdx].Type == l.Equal {
				currState = 'c'
				currIdx++
			} else {
				failure = true
			}
		case 'c':
			if b, exists, newI := binExpr(tokens, currIdx); exists {
				currState = 'e'
				result.Value = b
				currIdx = newI
			} else {
				currState = 'd'
			}
		case 'd':
			if f, exists, newI := factor(tokens, currIdx); exists {
				result.Value = f.(Literal)
				currState = 'e'
				currIdx = newI
			} else {
				failure = true
			}
		}
	}

	return result, !failure, currIdx
}

func declExpr(tokens []l.Token, index int) (DeclExpr, bool, int) {
	fmt.Println("declExpr")
	var result DeclExpr

	start := 'a'
	final := []byte{'c'}

	currState := start
	currIdx := index

	failure := false
	for !failure && !contains(final, byte(currState)) && currIdx < len(tokens) {
		switch currState {
		case 'a':
			switch tokens[currIdx].Type {
			case l.Int, l.Float, l.Bool, l.String, l.Char:
				result.Type = tokens[currIdx].Type
				currState = 'b'
				currIdx++
			default:
				failure = true
			}
		case 'b':
			if a, exists, newI := assignExpr(tokens, currIdx); exists {
				currState = 'c'
				result.Assign = a
				currIdx = newI
			} else {
				failure = true
			}
		}
	}

	return result, !failure, currIdx
}

func expr(tokens []l.Token, index int) (Expr, bool, int) {
	fmt.Println("expr")
	var result Expr

	start := 'a'
	final := []byte{'z'}

	currState := start
	currIdx := index

	failure := false
	for !failure && !contains(final, byte(currState)) && currIdx < len(tokens) {
		switch currState {
		case 'a':
			if d, exists, newI := declExpr(tokens, currIdx); exists {
				currState = 'y'
				result = d
				currIdx = newI
			} else {
				currState = 'b'
			}
		case 'b':
			if b, exists, newI := assignExpr(tokens, currIdx); exists {
				currState = 'y'
				result = b
				currIdx = newI
			} else {
				currState = 'c'
			}
		case 'c':
			if b, exists, newI := binExpr(tokens, currIdx); exists {
				currState = 'y'
				result = b
				currIdx = newI
			} else {
				currState = 'd'
			}
		case 'd':
			if f, exists, newI := factor(tokens, currIdx); exists {
				switch f.(type) {
				case Literal:
					result = f.(Literal)
					currState = 'y'
					currIdx = newI
				default:
					failure = true
				}
			} else {
				failure = true
			}
		case 'y':
			if tokens[currIdx].Type == l.Stop {
				currState = 'z'
				currIdx++
			} else {
				failure = true
			}
		}
	}

	return result, !failure, currIdx
}

func Parse(tokens []l.Token) (AST, error) {
	ast := AST{}
	seq := Seq{}
	stmts := []Stmt{}

	for i := 0; i < len(tokens); {
		if tokens[i].Type == l.EOF {
			break
		}

		e, exists, newI := expr(tokens, i)
		if exists {
			stmts = append(stmts, Stmt(ExprStmt{
				Expr: e,
			}))
			i = newI
		} else {
			token := tokens[i]

			seq.Stmts = stmts
			ast.Seq = seq

			return ast, fmt.Errorf("[%d:%d] Unknown token %s %s\n", token.Line, token.Col, token.Type.String(), token.Lexeme)
		}
	}

	seq.Stmts = stmts
	ast.Seq = seq
	return ast, nil
}
