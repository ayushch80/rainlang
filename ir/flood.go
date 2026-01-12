package ir

import (
	"encoding/binary"
	"fmt"

	l "rainlang/lexer"
	p "rainlang/parser"
)

var MAGIC = []byte{0x95, 0x36, 0x69, 0x67}
var VERSION = []byte{0x00, 0x0a}

const SEQ = 0xff

const ASSIGN = 0x21
const DECLARE = 0x22
const BINEXPR = 0x23
const LITERAL = 0x24

type Type int

const (
	INT Type = iota
	FLOAT
	STRING
	CHAR
	BOOL
	IDENT

	PLUS
	MINUS
	STAR
	SLASH
)

type Symbol struct {
	Index int
	Name  string
	Type  Type
}

type SymbolTable struct {
	Symbols []Symbol
}

type StringTable struct {
	Strings []string
}

func (sym Symbol) bytes() []byte {
	buf := make([]byte, 4)
	binary.LittleEndian.PutUint32(buf, uint32(sym.Index))

	name := []byte(sym.Name)
	nameLength := len(name)

	bytes := []byte{}
	bytes = append(bytes, buf...)
	bytes = append(bytes, byte(nameLength))
	bytes = append(bytes, name...)
	bytes = append(bytes, byte(sym.Type))

	return bytes
}

func (symt *SymbolTable) LookupSymbolTable(symbolName string) (Symbol, bool) {
	for _, sym := range symt.Symbols {
		if sym.Name == symbolName {
			return sym, true
		}
	}
	return Symbol{}, false
}

func (symt *SymbolTable) add(sym Symbol) (int, bool) {
	if _, ok := symt.LookupSymbolTable(sym.Name); ok {
		return -1, false
	}
	sym.Index = len(symt.Symbols)
	symt.Symbols = append(symt.Symbols, sym)
	return sym.Index, true
}

func (symt SymbolTable) bytes() []byte {
	bytes := []byte{}
	for _, sym := range symt.Symbols {
		bytes = append(bytes, sym.bytes()...)
	}
	return bytes
}

func (strt StringTable) Lookup(str string) int {
	for i, s := range strt.Strings {
		if s == str {
			return i
		}
	}
	return -1
}

func (strt *StringTable) add(str string) int {
	strt.Strings = append(strt.Strings, str)
	return len(strt.Strings) - 1
}

func (strt StringTable) bytes() []byte {
	bytes := []byte{}

	for _, str := range strt.Strings {
		bytes = append(bytes, byte(len(str)))
		bytes = append(bytes, []byte(str)...)
		bytes = append(bytes, 0x00)
	}

	return bytes
}

func GetType(t l.TokenType) (Type, error) {
	switch t {
	case l.Int, l.IntLiteral:
		return INT, nil
	case l.Float, l.FloatLiteral:
		return FLOAT, nil
	case l.Char, l.CharLiteral:
		return CHAR, nil
	case l.String, l.StringLiteral:
		return STRING, nil
	case l.Bool, l.BoolLiteral:
		return BOOL, nil
	case l.Ident:
		return IDENT, nil
	default:
		return Type(-1), fmt.Errorf("Unknown type %s", t.String())
	}
}

func getOperatorType(t l.TokenType) (Type, error) {
	switch t {
	case l.Plus:
		return PLUS, nil
	case l.Minus:
		return MINUS, nil
	case l.Star:
		return STAR, nil
	case l.Slash:
		return SLASH, nil
	default:
		return Type(-1), fmt.Errorf("Unknown operator %s", t.String())
	}
}

func GetFactorType(factor p.Factor) (Type, error) {
	switch factor.(type) {
	case p.Literal:
		return GetType(factor.(p.Literal).Type)
	case p.BinExpr:
		return GetFactorType(factor.(p.BinExpr).Left)
	}
	return Type(-1), fmt.Errorf("Unknown type of factor %s", factor)
}

func genLiteralBytecode(literal p.Literal, symt *SymbolTable, strt *StringTable, independent bool) ([]byte, error) {
	if independent {
		name := genRandomName()

		return genDeclExprBytecode(p.DeclExpr{
			Type: literal.Type,
			Assign: p.AssignExpr{
				Target: l.Token{
					Type:   l.Ident,
					Lexeme: name,
					Line:   -1,
					Col:    -1,
				},
				Value: literal,
			},
		}, symt, strt)
	}

	bytecode := []byte{LITERAL}

	literalType, err := GetFactorType(literal)
	if err != nil {
		return []byte{}, err
	}

	strtIndex := strt.Lookup(literal.Value)
	if strtIndex == -1 {
		strtIndex = strt.add(literal.Value)
	}

	buf := make([]byte, 4)
	binary.LittleEndian.PutUint32(buf, uint32(strtIndex))

	bytecode = append(bytecode, byte(literalType))
	bytecode = append(bytecode, buf...)

	return bytecode, nil
}

func genFactorBytecode(factor p.Factor, symt *SymbolTable, strt *StringTable, independent bool) ([]byte, error) {
	switch factor.(type) {
	case p.Literal:
		return genLiteralBytecode(factor.(p.Literal), symt, strt, independent)
	case p.BinExpr:
		return genBinExprBytecode(factor.(p.BinExpr), symt, strt, independent)
	default:
		return []byte{}, fmt.Errorf("Unknown type of factor")
	}
	return []byte{}, fmt.Errorf("Unknown type of factor")
}

func genBinExprBytecode(binExpr p.BinExpr, symt *SymbolTable, strt *StringTable, independent bool) ([]byte, error) {
	if independent {
		binExprTypeC, err := getFactorTypeC(binExpr)
		if err != nil {
			return []byte{}, err
		}

		name := genRandomName()

		return genDeclExprBytecode(p.DeclExpr{
			Type: binExprTypeC,
			Assign: p.AssignExpr{
				Target: l.Token{
					Type:   l.Ident,
					Lexeme: name,
					Line:   -1,
					Col:    -1,
				},
				Value: binExpr,
			},
		}, symt, strt)
	}

	bytecode := []byte{BINEXPR}

	// binExprType, err := GetFactorType(binExpr)
	// if err != nil {
	// 	return []byte{}, err
	// }

	bytecodeLeft, err := genFactorBytecode(binExpr.Left, symt, strt, false)
	if err != nil {
		return bytecode, err
	}

	bytecodeRight, err := genFactorBytecode(binExpr.Right, symt, strt, false)
	if err != nil {
		return bytecode, err
	}

	operatorCode, err := getOperatorType(binExpr.Op.Type)
	if err != nil {
		return bytecode, err
	}

	// bytecode = append(bytecode, byte(binExprType))

	bytecode = append(bytecode, byte(operatorCode))
	bytecode = append(bytecode, bytecodeLeft...)
	bytecode = append(bytecode, bytecodeRight...)
	bytecode = append(bytecode, 0x00)

	return bytecode, nil
}

func genAssignExprBytecode(assignExpr p.AssignExpr, symt *SymbolTable, strt *StringTable) ([]byte, error) {
	bytecode := []byte{ASSIGN}

	sym, ok := symt.LookupSymbolTable(assignExpr.Target.Lexeme)
	if !ok {
		fmt.Println(*symt)
		return []byte{}, fmt.Errorf("Cannot find the expected symbol '%s'", assignExpr.Target.Lexeme)
	}

	buf := make([]byte, 4)
	binary.LittleEndian.PutUint32(buf, uint32(sym.Index))

	bytecode = append(bytecode, buf...)

	exprBytecode, err := genExprBytecode(assignExpr.Value, symt, strt, false)
	if err != nil {
		return bytecode, nil
	}

	bytecode = append(bytecode, exprBytecode...)

	return bytecode, nil
}

func genDeclExprBytecode(declExpr p.DeclExpr, symt *SymbolTable, strt *StringTable) ([]byte, error) {
	bytecode := []byte{DECLARE}

	symbolType, err := GetType(declExpr.Type)
	if err != nil {
		return []byte{}, err
	}

	symtIndex, ok := symt.add(Symbol{
		Name: declExpr.Assign.Target.Lexeme,
		Type: symbolType,
	})
	if !ok {
		return bytecode, fmt.Errorf("Smth went off while adding new symbol %s")
	}

	buf := make([]byte, 4)
	binary.LittleEndian.PutUint32(buf, uint32(symtIndex))

	bytecode = append(bytecode, buf...)

	exprBytecode, err := genExprBytecode(declExpr.Assign.Value, symt, strt, false)
	if err != nil {
		return bytecode, nil
	}

	bytecode = append(bytecode, exprBytecode...)

	return bytecode, nil
}

func genExprBytecode(expr p.Expr, symt *SymbolTable, strt *StringTable, independent bool) ([]byte, error) {
	switch expr.(type) {
	case p.Literal:
		return genLiteralBytecode(expr.(p.Literal), symt, strt, independent)
	case p.BinExpr:
		return genBinExprBytecode(expr.(p.BinExpr), symt, strt, independent)
	case p.AssignExpr:
		return genAssignExprBytecode(expr.(p.AssignExpr), symt, strt)
	case p.DeclExpr:
		return genDeclExprBytecode(expr.(p.DeclExpr), symt, strt)
	default:
		return []byte{}, fmt.Errorf("Unknown type of expression %s", expr)
	}
	return []byte{}, nil
}

func genStmtBytecode(stmt p.Stmt, symt *SymbolTable, strt *StringTable) ([]byte, error) {
	switch stmt.(type) {
	case p.ExprStmt:
		return genExprBytecode(stmt.(p.ExprStmt).Expr, symt, strt, true)
	default:
		return []byte{}, fmt.Errorf("Unknown type of statement %s", stmt)
	}
	return []byte{}, nil
}

func genSeqBytecode(seq p.Seq, symt *SymbolTable, strt *StringTable) ([]byte, error) {
	bytecode := []byte{SEQ}
	for _, stmt := range seq.Stmts {
		stmtBytecode, err := genStmtBytecode(stmt, symt, strt)
		if err != nil {
			return bytecode, err
		}
		bytecode = append(bytecode, stmtBytecode...)
	}
	bytecode = append(bytecode, 0x00)
	return bytecode, nil
}

func GenerateFloodBytecode(ast p.AST) ([]byte, error) {
	bytecode := append(MAGIC, VERSION...)

	symt := SymbolTable{}
	strt := StringTable{}

	seqBytecode, err := genSeqBytecode(ast.Seq, &symt, &strt)
	if err != nil {
		return bytecode, err
	}

	symtBytes := symt.bytes()
	strtBytes := strt.bytes()

	bufA := make([]byte, 4)
	bufB := make([]byte, 4)
	bufC := make([]byte, 4)

	binary.LittleEndian.PutUint32(bufA, uint32(len(seqBytecode)))
	binary.LittleEndian.PutUint32(bufB, uint32(len(symtBytes)))
	binary.LittleEndian.PutUint32(bufC, uint32(len(strtBytes)))

	bytecode = append(bytecode, bufA...)
	bytecode = append(bytecode, bufB...)
	bytecode = append(bytecode, bufC...)

	bytecode = append(bytecode, seqBytecode...)
	bytecode = append(bytecode, symtBytes...)
	bytecode = append(bytecode, strtBytes...)

	return bytecode, nil
}
