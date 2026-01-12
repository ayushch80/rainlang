package debug

import (
	"fmt"

	// l "rainlang/lexer"
	. "rainlang/parser"
)

func PrintAST(ast AST) {
	printSeq(ast.Seq, 0)
}

func printSeq(seq Seq, indent int) {
	printIndent("SEQ", indent)
	if len(seq.Stmts) <= 0 {
		printIndent("EMPTY", indent+1)
		return
	}

	for _, stmt := range seq.Stmts {
		printStmt(stmt, indent)
	}
}

func printStmt(stmt Stmt, indent int) {
	switch stmt.(type) {
	case ExprStmt:
		exprStmt, _ := stmt.(ExprStmt)
		printExpr(exprStmt, indent)
	default:
		fmt.Println("[DEBUG ERROR] Unknown type of statement")
	}
}

func PrintStmt(stmt Stmt, indent int) {
	switch stmt.(type) {
	case ExprStmt:
		exprStmt, _ := stmt.(ExprStmt)
		printExpr(exprStmt, indent)
	default:
		fmt.Println("[DEBUG ERROR] Unknown type of statement")
	}
}

func PrintExpr(expr ExprStmt, indent int) {
	printExpr(expr, indent)
}

func printExpr(expr ExprStmt, indent int) {
	e := expr.Expr
	switch e.(type) {
	case BinExpr:
		b, _ := e.(BinExpr)
		printBinExpr(b, indent+1)
	case AssignExpr:
		a, _ := e.(AssignExpr)
		printAssignExpr(a, indent+1)
	case DeclExpr:
		d, _ := e.(DeclExpr)
		printDeclExpr(d, indent+1)
	case Literal:
		l, _ := e.(Literal)
		printFactor(l, indent+1)
	default:
		fmt.Println("[DEBUG ERROR] Unknown type of expr")
	}
}

func printBinExpr(binExpr BinExpr, indent int) {
	printIndent(fmt.Sprintf("BINARY %s", binExpr.Op.Type.String()), indent)
	printFactor(binExpr.Left, indent+1)
	printFactor(binExpr.Right, indent+1)
}

func printAssignExpr(assignExpr AssignExpr, indent int) {
	printIndent(fmt.Sprintf("ASSIGN %s", assignExpr.Target.Lexeme), indent)
	switch assignExpr.Value.(type) {
	case BinExpr:
		b, _ := assignExpr.Value.(BinExpr)
		printBinExpr(b, indent+1)
	case Literal:
		l, _ := assignExpr.Value.(Literal)
		printFactor(l, indent+1)
	}
}

func printDeclExpr(declExpr DeclExpr, indent int) {
	printIndent(fmt.Sprintf("DECL %s", declExpr.Type.String()), indent)
	printAssignExpr(declExpr.Assign, indent+1)
}

func printFactor(f Factor, indent int) {
	switch f.(type) {
	case Literal:
		literal := f.(Literal)
		printIndent(fmt.Sprintf("%s %s", literal.Type.String(), literal.Value), indent)
	case BinExpr:
		b := f.(BinExpr)
		printBinExpr(b, indent)
	default:
		fmt.Println("[DEBUG ERROR] Unknown type of factor")
	}
}

func printIndent(str string, indent int) {
	for _ = range indent {
		fmt.Print("  ")
	}
	fmt.Println(str)
}
