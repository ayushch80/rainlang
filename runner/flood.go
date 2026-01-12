package runner

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"strings"

	ir "rainlang/ir"
	l "rainlang/lexer"
	p "rainlang/parser"
)

var SUPPORTED_VERSIONS = [][]byte{
	[]byte{0x00, 0x0a},
}

func decodeSymbolTable(symtBytes []byte) (ir.SymbolTable, error) {
	symt := ir.SymbolTable{}

	for i := 0; i < len(symtBytes); {
		symIndex := int(binary.LittleEndian.Uint32(symtBytes[i : i+4]))
		nameLength := int(symtBytes[i+4])
		name := strings.Clone(string(symtBytes[i+5 : i+5+nameLength]))
		symType := ir.Type(symtBytes[i+5+nameLength])

		sym := ir.Symbol{
			Index: symIndex,
			Name:  name,
			Type:  symType,
		}

		symt.Symbols = append(symt.Symbols, sym)
		i += nameLength + 6
	}

	return symt, nil
}

func decodeStringTable(strtBytes []byte) (ir.StringTable, error) {
	strt := ir.StringTable{}
	for i := 0; i < len(strtBytes); {
		strLen := int(strtBytes[i])
		str := strings.Clone(string(strtBytes[i+1 : i+1+strLen]))
		strt.Strings = append(strt.Strings, str)
		i += 2 + strLen
	}
	return strt, nil
}

func checkMagic(bytecode []byte) bool {
	if len(bytecode) > 6 {
		if bytes.Equal(bytecode[:len(ir.MAGIC)], ir.MAGIC) {
			return true
		}
		return false
	}
	return false
}

func checkVersion(bytecode []byte) bool {
	if len(bytecode) > 6 {
		for _, version := range SUPPORTED_VERSIONS {
			if bytes.Equal(bytecode[len(ir.MAGIC):len(ir.MAGIC)+2], version) {
				return true
			}
		}
		return false
	}
	return true
}

func getLexerType(t ir.Type) (l.TokenType, error) {
	switch t {
	case ir.INT:
		return l.Int, nil
	case ir.FLOAT:
		return l.Float, nil
	case ir.BOOL:
		return l.Bool, nil
	case ir.STRING:
		return l.String, nil
	case ir.CHAR:
		return l.Char, nil
	default:
		return l.Unknown, fmt.Errorf("Unknown type")
	}
}

func getLexerLiteralType(t ir.Type) (l.TokenType, error) {
	switch t {
	case ir.INT:
		return l.IntLiteral, nil
	case ir.FLOAT:
		return l.FloatLiteral, nil
	case ir.BOOL:
		return l.BoolLiteral, nil
	case ir.STRING:
		return l.StringLiteral, nil
	case ir.CHAR:
		return l.CharLiteral, nil
	case ir.IDENT:
		return l.Ident, nil
	default:
		return l.Unknown, fmt.Errorf("Unknown type")
	}
}

func getLexerOperatorType(t ir.Type) (l.TokenType, error) {
	switch t {
	case ir.PLUS:
		return l.Plus, nil
	case ir.MINUS:
		return l.Minus, nil
	case ir.STAR:
		return l.Star, nil
	case ir.SLASH:
		return l.Slash, nil
	default:
		return l.Unknown, fmt.Errorf("Unknown operator")
	}
}

func converExprToFactor(expr p.Expr) (p.Factor, error) {
	switch expr.(type) {
	case p.Literal:
		return expr.(p.Literal), nil
	case p.BinExpr:
		return expr.(p.BinExpr), nil
	default:
		return p.Literal{}, fmt.Errorf("Cannot convert this expr to factor %s", expr)
	}
}

func decodeLiteral(bytecode []byte, symt *ir.SymbolTable, strt *ir.StringTable, index int) (p.Expr, int, error) {
	literalType := ir.Type(bytecode[index])
	lexerLiteralType, err := getLexerLiteralType(literalType)
	if err != nil {
		return p.Literal{}, -1, err
	}

	strtIndex := int(binary.LittleEndian.Uint32(bytecode[index+1 : index+5]))
	str := strt.Strings[strtIndex]

	return p.Literal{
		Type:  lexerLiteralType,
		Value: str,
	}, index + 5, nil
}

func decodeBinExpr(bytecode []byte, symt *ir.SymbolTable, strt *ir.StringTable, index int) (p.Expr, int, error) {
	binExpr := p.BinExpr{}

	operatorType := ir.Type(bytecode[index])
	operator, err := getLexerOperatorType(operatorType)

	left, index, err := decodeExpr(bytecode, symt, strt, index+1)
	if err != nil {
		return p.BinExpr{}, -1, err
	}
	right, index, err := decodeExpr(bytecode, symt, strt, index)
	if err != nil {
		return p.BinExpr{}, -1, err
	}

	leftFactor, err := converExprToFactor(left.(p.ExprStmt).Expr)
	if err != nil {
		return p.BinExpr{}, -1, err
	}
	rightFactor, err := converExprToFactor(right.(p.ExprStmt).Expr)
	if err != nil {
		return p.BinExpr{}, -1, err
	}

	binExpr.Left = leftFactor
	binExpr.Right = rightFactor
	binExpr.Op = p.BinOp{
		Type: operator,
	}

	if index >= len(bytecode) || bytecode[index] != 0x00 {
		return p.BinExpr{}, -1, fmt.Errorf("BINEXPR: missing terminator")
	}
	return binExpr, index + 1, nil
}

func decodeAssignExpr(bytecode []byte, symt *ir.SymbolTable, strt *ir.StringTable, index int) (p.Expr, int, error) {
	symtIndex := int(binary.LittleEndian.Uint32(bytecode[index : index+4]))
	sym := symt.Symbols[symtIndex]

	value, index, err := decodeExpr(bytecode, symt, strt, index+4)
	if err != nil {
		return p.AssignExpr{}, -1, err
	}

	assignExpr := p.AssignExpr{
		Target: l.Token{
			Type:   l.Ident,
			Lexeme: sym.Name,
			Col:    -1,
			Line:   -1,
		},
		Value: value.(p.ExprStmt).Expr,
	}

	return assignExpr, index, nil
}

func decodeDeclExpr(bytecode []byte, symt *ir.SymbolTable, strt *ir.StringTable, index int) (p.Expr, int, error) {
	symtIndex := int(binary.LittleEndian.Uint32(bytecode[index : index+4]))
	sym := symt.Symbols[symtIndex]

	lexerType, err := getLexerType(sym.Type)
	if err != nil {
		return p.DeclExpr{}, -1, err
	}
	value, index, err := decodeExpr(bytecode, symt, strt, index+4)
	if err != nil {
		return p.DeclExpr{}, -1, err
	}

	declExpr := p.DeclExpr{
		Type: lexerType,
		Assign: p.AssignExpr{
			Target: l.Token{
				Type:   l.Ident,
				Lexeme: sym.Name,
				Col:    -1,
				Line:   -1,
			},
			Value: value.(p.ExprStmt).Expr,
		},
	}

	return declExpr, index, nil
}

func decodeExpr(bytecode []byte, symt *ir.SymbolTable, strt *ir.StringTable, index int) (p.Stmt, int, error) {
	var expr p.Expr
	var newIdx int
	var err error
	if index >= len(bytecode) {
		return p.ExprStmt{}, -1, fmt.Errorf("EXPR: index %d out of range", index)
	}
	if bytecode[index] == 0x00 {
		return p.ExprStmt{}, index, nil
	}
	switch int(bytecode[index]) {
	case 0x21:
		expr, newIdx, err = decodeAssignExpr(bytecode, symt, strt, index+1)
	case 0x22:
		expr, newIdx, err = decodeDeclExpr(bytecode, symt, strt, index+1)
	case 0x23:
		expr, newIdx, err = decodeBinExpr(bytecode, symt, strt, index+1)
	case 0x24:
		expr, newIdx, err = decodeLiteral(bytecode, symt, strt, index+1)
	default:
		return p.ExprStmt{}, -1, fmt.Errorf("EXPR: Unexpected byte 0x%x at index %d", bytecode[index], index)
	}
	return p.ExprStmt{
		Expr: expr,
	}, newIdx, err
}

func decodeStmts(bytecode []byte, symt *ir.SymbolTable, strt *ir.StringTable, index int) ([]p.Stmt, error) {
	var stmts []p.Stmt

	for {
		if index >= len(bytecode) {
			return stmts, fmt.Errorf("SEQ: missing terminator")
		}
		if bytecode[index] == 0x00 {
			return stmts, nil
		}
		stmt, newIdx, err := decodeExpr(bytecode, symt, strt, index)
		if err != nil {
			return stmts, err
		}
		stmts = append(stmts, stmt)
		index = newIdx
		// d.PrintStmt(stmt, 0)
	}

	return stmts, nil
}

func decodeSeq(bytecode []byte, symt *ir.SymbolTable, strt *ir.StringTable) (p.Seq, error) {
	index := 0
	if bytecode[index] == ir.SEQ {
		stmts, err := decodeStmts(bytecode, symt, strt, index+1)
		if err != nil {
			return p.Seq{}, err
		}
		return p.Seq{
			Stmts: stmts,
		}, nil
	}
	return p.Seq{}, fmt.Errorf("SEQ: Unexpected byte 0x%x at index %d", bytecode[index], index)
}

func DecodeFloodBytecode(bytecode []byte) (p.AST, error) {
	ast := p.AST{}

	if !checkMagic(bytecode) {
		return ast, fmt.Errorf("Invalid magic")
	}

	if !checkVersion(bytecode) {
		return ast, fmt.Errorf("Unsupported version")
	}

	lenSeqBytecode := int(binary.LittleEndian.Uint32(bytecode[6:10]))
	lenSymtBytes := int(binary.LittleEndian.Uint32(bytecode[10:14]))
	lenStrtBytes := int(binary.LittleEndian.Uint32(bytecode[14:18]))

	seqBytecode := make([]byte, lenSeqBytecode)
	symtBytes := make([]byte, lenSymtBytes)
	strtBytes := make([]byte, lenStrtBytes)

	copy(seqBytecode, bytecode[18:18+lenSeqBytecode])
	copy(symtBytes, bytecode[18+lenSeqBytecode:18+lenSeqBytecode+lenSymtBytes])
	copy(strtBytes, bytecode[18+lenSeqBytecode+lenSymtBytes:18+lenSeqBytecode+lenSymtBytes+lenStrtBytes])

	symt, err := decodeSymbolTable(symtBytes)
	if err != nil {
		return ast, err
	}

	strt, err := decodeStringTable(strtBytes)
	if err != nil {
		return ast, err
	}

	seq, err := decodeSeq(seqBytecode, &symt, &strt)
	if err != nil {
		return ast, err
	}

	ast.Seq = seq

	return ast, nil
}
