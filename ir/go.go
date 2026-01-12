package ir

import (
	"fmt"

	l "rainlang/lexer"
	p "rainlang/parser"
)

func getGoType(lexerType l.TokenType) (string, error) {
	switch lexerType {
	case l.Int, l.IntLiteral:
		return "int", nil
	case l.Float, l.FloatLiteral:
		return "float", nil
	case l.Char, l.CharLiteral:
		return "byte", nil
	case l.String, l.StringLiteral:
		return "string", nil
	case l.Bool, l.BoolLiteral:
		return "bool", nil
	case l.Unknown:
		return "", nil
	default:
		return "", fmt.Errorf("Unknown type %s", lexerType.String())
	}
}

func addGoFunc(funcName string, funcType l.TokenType, code *string) error {
	stringFuncType, err := getGoType(funcType)
	if err != nil {
		return err
	}

	*code += fmt.Sprintf("func %s() %s {\n", funcName, stringFuncType)
	return nil
}

func genLiteralCode(literal p.Literal, code *string, independent bool) error {
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
		name := genRandomName()
		*code += fmt.Sprintf("%s := %s\n", name, value)
		*code += fmt.Sprintf("__prevent_go_error__(%s)\n", name)
	} else {
		*code += fmt.Sprintf("%s", value)
	}

	return nil
}

func genFactorCode(factor p.Factor, code *string, independent bool) error {
	switch factor.(type) {
	case p.Literal:
		return genLiteralCode(factor.(p.Literal), code, independent)
	case p.BinExpr:
		return genBinExprCode(factor.(p.BinExpr), code, independent)
	default:
		return fmt.Errorf("Unknown type of factor %s", factor)
	}
	return nil
}

func genBinExprCode(binExpr p.BinExpr, code *string, independent bool) error {
	value := "("

	err := genFactorCode(binExpr.Left, &value, false)
	if err != nil {
		return err
	}

	operator, err := getOperatorSymbol(binExpr.Op.Type)
	if err != nil {
		return err
	}
	value += operator

	err = genFactorCode(binExpr.Right, &value, false)
	if err != nil {
		return err
	}

	value += ")"

	if independent {
		name := genRandomName()
		*code += fmt.Sprintf("%s := %s\n", name, value)
		*code += fmt.Sprintf("__prevent_go_error__(%s)\n", name)
	} else {
		*code += fmt.Sprintf("%s", value)
	}
	return nil
}

func genAssignExprCode(assignExpr p.AssignExpr, code *string, isNew bool) error {
	if isNew {
		*code += fmt.Sprintf("%s := ", assignExpr.Target.Lexeme)
	} else {
		*code += fmt.Sprintf("%s = ", assignExpr.Target.Lexeme)
	}
	err := genExprCode(assignExpr.Value, code, false)
	if err != nil {
		return err
	}
	*code += "\n"

	if isNew {
		*code += fmt.Sprintf("__prevent_go_error__(%s)\n", assignExpr.Target.Lexeme)
	}

	return nil
}

func genDeclExprCodeGo(declExpr p.DeclExpr, code *string) error {
	return genAssignExprCode(declExpr.Assign, code, true)
}

func genExprCode(expr p.Expr, code *string, independent bool) error {
	switch expr.(type) {
	case p.Literal:
		return genLiteralCode(expr.(p.Literal), code, independent)
	case p.BinExpr:
		return genBinExprCode(expr.(p.BinExpr), code, independent)
	case p.AssignExpr:
		return genAssignExprCode(expr.(p.AssignExpr), code, false)
	case p.DeclExpr:
		return genDeclExprCodeGo(expr.(p.DeclExpr), code)
	default:
		return fmt.Errorf("Unknown type of expression %s", expr)
	}
	return nil
}

func genStmtCode(stmt p.Stmt, code *string) error {
	switch stmt.(type) {
	case p.ExprStmt:
		return genExprCode(stmt.(p.ExprStmt).Expr, code, true)
	default:
		return fmt.Errorf("Unknown type of statement %s", stmt)
	}
	return nil
}

func genSeqCode(seq p.Seq, code *string) error {
	for _, stmt := range seq.Stmts {
		err := genStmtCode(stmt, code)
		if err != nil {
			return err
		}
	}
	return nil
}

func GenerateGoCode(ast p.AST) (string, error) {
	code := "package main\n\n"

	// err := addGoFunc("__prevent_go_error__", l.Unknown)
	code += "func __prevent_go_error__(a any) {}\n"

	err := addGoFunc("main", l.Unknown, &code)
	if err != nil {
		return code, err
	}

	err = genSeqCode(ast.Seq, &code)
	if err != nil {
		return code, err
	}

	endFunc(&code)

	return code, nil
}
