package ir

import (
	"fmt"
	"math/rand"
	"time"

	l "rainlang/lexer"
	p "rainlang/parser"
)

func genRandomName() string {
	const allowed = "_abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	length := 10

	seededRand := rand.New(rand.NewSource(time.Now().UnixNano()))

	b := make([]byte, length)
	for i := range b {
		b[i] = allowed[seededRand.Intn(len(allowed))]
	}
	return string(b)
}

func addHeaders(code *string, headers ...string) {
	for _, header := range headers {
		*code += fmt.Sprintf("#include <%s>\n", header)
	}
}

func getCType(lexerType l.TokenType) (string, error) {
	switch lexerType {
	case l.Int, l.IntLiteral:
		return "int", nil
	case l.Float, l.FloatLiteral:
		return "float", nil
	case l.Char, l.CharLiteral:
		return "char", nil
	case l.String, l.StringLiteral:
		return "char *", nil
	case l.Bool, l.BoolLiteral:
		return "bool", nil
	case l.Unknown:
		return "void", nil
	default:
		return "", fmt.Errorf("Unknown type %s", lexerType.String())
	}
}

func getOperatorSymbol(op l.TokenType) (string, error) {
	switch op {
	case l.Plus:
		return "+", nil
	case l.Minus:
		return "-", nil
	case l.Star:
		return "*", nil
	case l.Slash:
		return "/", nil
	default:
		return "", fmt.Errorf("Unknown operator %s", op.String())
	}

	return "", nil
}

func addCFunc(funcName string, funcType l.TokenType, code *string) error {
	stringFuncType, err := getCType(funcType)
	if err != nil {
		return err
	}

	*code += fmt.Sprintf("%s %s () {\n", stringFuncType, funcName)
	return nil
}

func endFunc(code *string) {
	*code += "}\n"
}

func genLiteralCodeC(literal p.Literal, code *string, independent bool) error {
	var value string
	switch literal.Type {
	case l.Char, l.CharLiteral:
		value = fmt.Sprintf("'%s'", literal.Value)
	case l.String, l.StringLiteral:
		value = fmt.Sprintf("\"%s\"", literal.Value)
	default:
		value = literal.Value
	}

	if independent {
		cType, err := getCType(literal.Type)
		if err != nil {
			return err
		}
		*code += fmt.Sprintf("%s %s = %s;\n", cType, genRandomName(), value)
	} else {
		*code += fmt.Sprintf("%s", value)
	}

	return nil
}

func genFactorCodeC(factor p.Factor, code *string, independent bool) error {
	switch factor.(type) {
	case p.Literal:
		return genLiteralCodeC(factor.(p.Literal), code, independent)
	case p.BinExpr:
		return genBinExprCodeC(factor.(p.BinExpr), code, independent)
	default:
		return fmt.Errorf("Unknown type of factor %s", factor)
	}
	return nil
}

func getFactorTypeC(factor p.Factor) (l.TokenType, error) {
	switch factor.(type) {
	case p.Literal:
		return factor.(p.Literal).Type, nil
	case p.BinExpr:
		return getFactorTypeC(factor.(p.BinExpr).Left)
	}
	return l.Unknown, fmt.Errorf("Unknown type of factor %s", factor)
}

func genBinExprCodeC(binExpr p.BinExpr, code *string, independent bool) error {
	value := "("

	err := genFactorCodeC(binExpr.Left, &value, false)
	if err != nil {
		return err
	}

	operator, err := getOperatorSymbol(binExpr.Op.Type)
	if err != nil {
		return err
	}
	value += operator

	err = genFactorCodeC(binExpr.Right, &value, false)
	if err != nil {
		return err
	}

	value += ")"

	if independent {
		exprType, err := getFactorTypeC(binExpr)
		if err != nil {
			return err
		}
		cType, err := getCType(exprType)
		if err != nil {
			return err
		}
		*code += fmt.Sprintf("%s %s = %s;\n", cType, genRandomName(), value)
	} else {
		*code += fmt.Sprintf("%s", value)
	}
	return nil
}

func genAssignExprCodeC(assignExpr p.AssignExpr, code *string) error {

	*code += fmt.Sprintf("%s = ", assignExpr.Target.Lexeme)
	err := genExprCodeC(assignExpr.Value, code, false)
	if err != nil {
		return err
	}
	*code += ";\n"

	return nil
}

func genDeclExprCodeC(declExpr p.DeclExpr, code *string) error {
	targetCType, err := getCType(declExpr.Type)
	if err != nil {
		return err
	}

	*code += fmt.Sprintf("%s ", targetCType)

	return genAssignExprCodeC(declExpr.Assign, code)
}

func genExprCodeC(expr p.Expr, code *string, independent bool) error {
	switch expr.(type) {
	case p.Literal:
		return genLiteralCodeC(expr.(p.Literal), code, independent)
	case p.BinExpr:
		return genBinExprCodeC(expr.(p.BinExpr), code, independent)
	case p.AssignExpr:
		return genAssignExprCodeC(expr.(p.AssignExpr), code)
	case p.DeclExpr:
		return genDeclExprCodeC(expr.(p.DeclExpr), code)
	default:
		return fmt.Errorf("Unknown type of expression %s", expr)
	}
	return nil
}

func genStmtCodeC(stmt p.Stmt, code *string) error {
	switch stmt.(type) {
	case p.ExprStmt:
		return genExprCodeC(stmt.(p.ExprStmt).Expr, code, true)
	default:
		return fmt.Errorf("Unknown type of statement %s", stmt)
	}
	return nil
}

func genSeqCodeC(seq p.Seq, code *string) error {
	for _, stmt := range seq.Stmts {
		err := genStmtCodeC(stmt, code)
		if err != nil {
			return err
		}
	}
	return nil
}

func GenerateCCode(ast p.AST) (string, error) {
	var code string

	addHeaders(&code, "stdio.h", "stdbool.h")

	err := addCFunc("main", l.Unknown, &code)
	if err != nil {
		return code, err
	}

	err = genSeqCodeC(ast.Seq, &code)
	if err != nil {
		return code, err
	}

	endFunc(&code)

	return code, nil
}
