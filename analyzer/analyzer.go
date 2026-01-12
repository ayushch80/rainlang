package analyzer

import (
	"fmt"

	l "rainlang/lexer"
	p "rainlang/parser"
)

type Symbol struct {
	Name string
	Type l.TokenType
}

type SymbolTable struct {
	Symbols []Symbol
}

func matchSymbol(targetName string, sym Symbol) bool {
	if targetName == sym.Name {
		return true
	}
	return false
}

func newSymbol(symbolName string, symbolType l.TokenType, currentST *SymbolTable) error {
	_, exists := getSymbol(symbolName, currentST, nil)
	if exists {
		return fmt.Errorf("Symbol '%s' already exists in the current context", symbolName)
	}

	currentST.Symbols = append(currentST.Symbols, Symbol{
		Name: symbolName,
		Type: symbolType,
	})

	return nil
}

func getSymbol(symbolName string, currentST *SymbolTable, parentST *SymbolTable) (*Symbol, bool) {
	for _, symbol := range currentST.Symbols {
		if matchSymbol(symbolName, symbol) {
			return &symbol, true
		}
	}
	if parentST == nil {
		return nil, false
	}
	for _, symbol := range parentST.Symbols {
		if matchSymbol(symbolName, symbol) {
			return &symbol, true
		}
	}
	return nil, false
}

// func (st SymbolTable) getSymbol(symbolName string, parentST *SymbolTable) (*Symbol, bool) {

// }

func analyzeLiteral(literal p.Literal, currentST *SymbolTable, parentST *SymbolTable) (l.TokenType, error) {
	switch literal.Type {
	case l.Ident:
		s, exists := getSymbol(literal.Value, currentST, parentST)
		if !exists {
			return l.Unknown, fmt.Errorf("Identifier '%s' doesn't exist in the symbol table", literal.Value)
		}
		return s.Type, nil
	case l.IntLiteral, l.FloatLiteral, l.StringLiteral, l.CharLiteral, l.BoolLiteral:
		// fmt.Println("type", literal.Type)
		return literal.Type, nil
	}
	return l.Unknown, fmt.Errorf("Inavlid literal %s", literal)
}

func analyzeFactor(factor p.Factor, currentST *SymbolTable, parentST *SymbolTable) (l.TokenType, error) {
	switch factor.(type) {
	case p.Literal:
		return analyzeLiteral(factor.(p.Literal), currentST, parentST)
	case p.BinExpr:
		return analyzeBinExpr(factor.(p.BinExpr), currentST, parentST)
	default:
		return l.Unknown, fmt.Errorf("Invalid factor %s", factor)
	}
}

func analyzeBinExpr(binExpr p.BinExpr, currentST *SymbolTable, parentST *SymbolTable) (l.TokenType, error) {
	dtypeLeft, err := analyzeFactor(binExpr.Left, currentST, parentST)
	if err != nil {
		return l.Unknown, err
	}
	dtypeRight, err := analyzeFactor(binExpr.Right, currentST, parentST)
	if err != nil {
		return l.Unknown, err
	}

	if dtypeLeft != dtypeRight {
		return dtypeLeft, fmt.Errorf("Datatype should be same for all elements in a binary expression %s", binExpr)
	}

	return dtypeLeft, nil
}

func analyzeAssignExpr(assignExpr p.AssignExpr, currentST *SymbolTable, parentST *SymbolTable) (l.TokenType, error) {
	targetSym, exists := getSymbol(assignExpr.Target.Lexeme, currentST, parentST)
	if !exists {
		return l.Unknown, fmt.Errorf("Identifier '%s' doesn't exist in the symbol table", assignExpr.Target.Lexeme)
	}

	// fmt.Println(assignExpr.Value)

	resultType, err := analyzeExpr(assignExpr.Value, currentST, parentST)
	if err != nil {
		return l.Unknown, err
	}

	// fmt.Println(resultType, targetSym.Type)

	if resultType != targetSym.Type {
		return l.Unknown, fmt.Errorf("Incorrect datatype, cannot be assigned to the symbol '%s'(%s) --- %s", targetSym.Name, *targetSym, assignExpr.Value)
	}

	return resultType, nil
}

func analyzeDeclExpr(declExpr p.DeclExpr, currentST *SymbolTable, parentST *SymbolTable) (l.TokenType, error) {
	declType := declExpr.Type
	declSymName := declExpr.Assign.Target.Lexeme

	switch declType {
	case l.Int:
		declType = l.IntLiteral
	case l.Float:
		declType = l.FloatLiteral
	case l.String:
		declType = l.StringLiteral
	case l.Char:
		declType = l.CharLiteral
	case l.Bool:
		declType = l.BoolLiteral
	default:
		declType = l.Unknown
	}

	err := newSymbol(declSymName, declType, currentST)
	if err != nil {
		return l.Unknown, err
	}

	return analyzeAssignExpr(declExpr.Assign, currentST, parentST)
}

func analyzeExpr(expr p.Expr, currentST *SymbolTable, parentST *SymbolTable) (l.TokenType, error) {
	switch expr.(type) {
	case p.Literal:
		return analyzeLiteral(expr.(p.Literal), currentST, parentST)
	case p.BinExpr:
		return analyzeBinExpr(expr.(p.BinExpr), currentST, parentST)
	case p.AssignExpr:
		return analyzeAssignExpr(expr.(p.AssignExpr), currentST, parentST)
	case p.DeclExpr:
		return analyzeDeclExpr(expr.(p.DeclExpr), currentST, parentST)
	default:
		return l.Unknown, fmt.Errorf("Invalid expr %s", expr)
	}
	return l.Unknown, nil
}

func analyzeStmt(stmt p.Stmt, currentST *SymbolTable, parentST *SymbolTable) error {
	switch stmt.(type) {
	case p.ExprStmt:
		_, err := analyzeExpr(stmt.(p.ExprStmt).Expr, currentST, parentST)
		return err
	default:
		return fmt.Errorf("Invalid stmt %s", stmt)
	}
	return nil
}

func analyzeSeq(seq p.Seq, parentST *SymbolTable) error {
	var st SymbolTable
	for _, stmt := range seq.Stmts {
		err := analyzeStmt(stmt, &st, parentST)
		if err != nil {
			return err
		}
	}
	return nil
}

func Analyze(ast p.AST) error {
	return analyzeSeq(ast.Seq, nil)
}
