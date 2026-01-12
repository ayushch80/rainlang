package runner

import (
	"fmt"
	"strconv"

	ir "rainlang/ir"
	l "rainlang/lexer"
	p "rainlang/parser"
)

type Symbol struct {
	Name  string
	Value string
	Type  ir.Type
}

type SymbolTable struct {
	Symbols []Symbol
	Parent  *SymbolTable
}

func (sym Symbol) print() {
	fmt.Printf("%s = `%s`\n", sym.Name, sym.Value)
}

func (st SymbolTable) print() {
	for _, sym := range st.Symbols {
		sym.print()
	}
}

func (st *SymbolTable) lookupLocal(name string) (*Symbol, error) {
	for _, sym := range st.Symbols {
		if sym.Name == name {
			return &sym, nil
		}
	}
	return nil, fmt.Errorf("Unable to find symbol '%s' in the local symbol table", name)
}

func (st *SymbolTable) lookup(name string) (*Symbol, error) {
	for _, sym := range st.Symbols {
		if sym.Name == name {
			return &sym, nil
		}
	}
	if st.Parent != nil {
		return st.Parent.lookup(name)
	}
	return nil, fmt.Errorf("Unable to find symbol '%s' in the symbol table", name)
}

func (symt *SymbolTable) add(sym Symbol) error {
	if _, err := symt.lookupLocal(sym.Name); err == nil {
		return fmt.Errorf("Symbol already exists")
	}
	symt.Symbols = append(symt.Symbols, sym)
	return nil
}

func computePlus(exprType ir.Type, left, right string) (string, error) {
	switch exprType {
	case ir.INT:
		leftInt, err := strconv.Atoi(left)
		if err != nil {
			return "", err
		}
		rightInt, err := strconv.Atoi(right)
		if err != nil {
			return "", err
		}
		result := leftInt + rightInt
		resultStr := strconv.Itoa(result)
		return resultStr, nil
	case ir.FLOAT:
		leftFloat, err := strconv.ParseFloat(left, 64)
		if err != nil {
			return "", err
		}
		rightFloat, err := strconv.ParseFloat(right, 64)
		if err != nil {
			return "", err
		}
		result := leftFloat + rightFloat
		resultStr := strconv.FormatFloat(result, 'f', -1, 64)
		return resultStr, nil
	case ir.BOOL:
		return "", fmt.Errorf("Cannot do binary arithmetic operations in booleans")
	case ir.CHAR:
		leftRune := rune(left[0])
		rightRune := rune(right[0])
		result := leftRune + rightRune
		return string(result), nil
	case ir.STRING:
		result := left + right
		return result, nil
	default:
		fmt.Println(exprType, left, right)
		return "", fmt.Errorf("Unknown kind of binary operation")
	}
}

func computeMinus(exprType ir.Type, left, right string) (string, error) {
	switch exprType {
	case ir.INT:
		leftInt, err := strconv.Atoi(left)
		if err != nil {
			return "", err
		}
		rightInt, err := strconv.Atoi(right)
		if err != nil {
			return "", err
		}
		result := leftInt - rightInt
		resultStr := strconv.Itoa(result)
		return resultStr, nil
	case ir.FLOAT:
		leftFloat, err := strconv.ParseFloat(left, 64)
		if err != nil {
			return "", err
		}
		rightFloat, err := strconv.ParseFloat(right, 64)
		if err != nil {
			return "", err
		}
		result := leftFloat - rightFloat
		resultStr := strconv.FormatFloat(result, 'f', -1, 64)
		return resultStr, nil
	case ir.BOOL:
		return "", fmt.Errorf("Cannot do binary arithmetic operations in booleans")
	case ir.CHAR:
		leftRune := rune(left[0])
		rightRune := rune(right[0])
		result := leftRune - rightRune
		return string(result), nil
	case ir.STRING:
		return "", fmt.Errorf("Cannot do `-` operation on strings")
	default:
		fmt.Println(exprType, left, right)
		return "", fmt.Errorf("Unknown kind of binary operation")
	}
}

func computeStar(exprType ir.Type, left, right string) (string, error) {
	switch exprType {
	case ir.INT:
		leftInt, err := strconv.Atoi(left)
		if err != nil {
			return "", err
		}
		rightInt, err := strconv.Atoi(right)
		if err != nil {
			return "", err
		}
		result := leftInt * rightInt
		resultStr := strconv.Itoa(result)
		return resultStr, nil
	case ir.FLOAT:
		leftFloat, err := strconv.ParseFloat(left, 64)
		if err != nil {
			return "", err
		}
		rightFloat, err := strconv.ParseFloat(right, 64)
		if err != nil {
			return "", err
		}
		result := leftFloat * rightFloat
		resultStr := strconv.FormatFloat(result, 'f', -1, 64)
		return resultStr, nil
	case ir.BOOL:
		return "", fmt.Errorf("Cannot do binary arithmetic operations in booleans")
	case ir.CHAR:
		leftRune := rune(left[0])
		rightRune := rune(right[0])
		result := leftRune * rightRune
		return string(result), nil
	case ir.STRING:
		return "", fmt.Errorf("Cannot do `*` operation on strings")
	default:
		fmt.Println(exprType, left, right)
		return "", fmt.Errorf("Unknown kind of binary operation")
	}
}

func computeSlash(exprType ir.Type, left, right string) (string, error) {
	switch exprType {
	case ir.INT:
		leftInt, err := strconv.Atoi(left)
		if err != nil {
			return "", err
		}
		rightInt, err := strconv.Atoi(right)
		if err != nil {
			return "", err
		}
		result := leftInt / rightInt
		resultStr := strconv.Itoa(result)
		return resultStr, nil
	case ir.FLOAT:
		leftFloat, err := strconv.ParseFloat(left, 64)
		if err != nil {
			return "", err
		}
		rightFloat, err := strconv.ParseFloat(right, 64)
		if err != nil {
			return "", err
		}
		result := leftFloat / rightFloat
		resultStr := strconv.FormatFloat(result, 'f', -1, 64)
		return resultStr, nil
	case ir.BOOL:
		return "", fmt.Errorf("Cannot do binary arithmetic operations in booleans")
	case ir.CHAR:
		leftRune := rune(left[0])
		rightRune := rune(right[0])
		result := leftRune / rightRune
		return string(result), nil
	case ir.STRING:
		return "", fmt.Errorf("Cannot do `/` operation on strings")
	default:
		fmt.Println(exprType, left, right)
		return "", fmt.Errorf("Unknown kind of binary operation")
	}
}

func computeLiteralResult(literal p.Literal, st *SymbolTable) (string, error) {
	switch literal.Type {
	case l.IntLiteral, l.FloatLiteral, l.StringLiteral, l.CharLiteral, l.BoolLiteral:
		return literal.Value, nil
	case l.Ident:
		sym, err := st.lookup(literal.Value)
		if err != nil {
			return "", err
		}
		return sym.Value, nil
	default:
		return "", fmt.Errorf("Invalid type of literal")
	}
}

func computeBinExprResult(binExpr p.BinExpr, st *SymbolTable, exprType ir.Type) (string, error) {
	resultLeft, err := computeFactorResult(binExpr.Left, st, exprType)
	if err != nil {
		return "", err
	}

	resultRight, err := computeFactorResult(binExpr.Right, st, exprType)
	if err != nil {
		return "", err
	}

	switch binExpr.Op.Type {
	case l.Plus:
		return computePlus(exprType, resultLeft, resultRight)
	case l.Minus:
		return computeMinus(exprType, resultLeft, resultRight)
	case l.Star:
		return computeStar(exprType, resultLeft, resultRight)
	case l.Slash:
		return computeSlash(exprType, resultLeft, resultRight)
	default:
		return "", fmt.Errorf("Unknown operator %s", binExpr.Op.Type)
	}
}

func runAssignExpr(assignExpr p.AssignExpr, st *SymbolTable) error {
	sym, err := st.lookup(assignExpr.Target.Lexeme)
	if err != nil {
		return err
	}

	result, err := computeExprResult(assignExpr.Value, st, sym.Type)
	if err != nil {
		return err
	}

	sym.Value = result
	return nil
}

func runDeclExpr(declExpr p.DeclExpr, st *SymbolTable) error {
	_, err := st.lookupLocal(declExpr.Assign.Target.Lexeme)
	if err == nil {
		return fmt.Errorf("Symbol '%s' already exists", declExpr.Assign.Target.Lexeme)
	}

	irType, err := ir.GetType(declExpr.Type)
	if err != nil {
		return err
	}

	value, err := computeExprResult(declExpr.Assign.Value, st, irType)
	if err != nil {
		return err
	}

	err = st.add(Symbol{
		Name:  declExpr.Assign.Target.Lexeme,
		Type:  irType,
		Value: value,
	})
	if err != nil {
		return err
	}

	return nil
}

func computeExprResult(expr p.Expr, st *SymbolTable, exprType ir.Type) (string, error) {
	switch expr.(type) {
	case p.Literal:
		return computeLiteralResult(expr.(p.Literal), st)
	case p.BinExpr:
		return computeBinExprResult(expr.(p.BinExpr), st, exprType)
	default:
		return "", fmt.Errorf("Unknown type of computable expr")
	}
}

func computeFactorResult(factor p.Factor, st *SymbolTable, exprType ir.Type) (string, error) {
	switch factor.(type) {
	case p.Literal:
		return computeLiteralResult(factor.(p.Literal), st)
	case p.BinExpr:
		return computeBinExprResult(factor.(p.BinExpr), st, exprType)
	default:
		return "", fmt.Errorf("Unknown type of computable factor")
	}
}

func runExpr(expr p.Expr, st *SymbolTable) error {
	switch expr.(type) {
	case p.AssignExpr:
		return runAssignExpr(expr.(p.AssignExpr), st)
	case p.DeclExpr:
		return runDeclExpr(expr.(p.DeclExpr), st)
	default:
		return fmt.Errorf("Unknown type of runnable expr")
	}
}

func runStmt(stmt p.Stmt, st *SymbolTable) error {
	switch stmt.(type) {
	case p.ExprStmt:
		return runExpr(stmt.(p.ExprStmt).Expr, st)
	default:
		return fmt.Errorf("Unknown type of stmt")
	}
}

func runSeq(seq p.Seq, st *SymbolTable) error {
	for _, stmt := range seq.Stmts {
		err := runStmt(stmt, st)
		if err != nil {
			return err
		}
	}
	return nil
}

func Run(ast p.AST) error {
	st := SymbolTable{
		Parent: nil,
	}

	err := runSeq(ast.Seq, &st)
	if err != nil {
		return err
	}

	st.print()

	return nil
}
