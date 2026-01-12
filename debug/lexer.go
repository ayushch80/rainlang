package debug

import (
	"fmt"

	. "rainlang/lexer"
)

func PrintTokens(tokens []Token) {
	for _, token := range tokens {
		printToken(token)
	}
}

func printToken(token Token) {
	switch token.Type {
	case EOF:
		fmt.Println("EOF")
	case Unknown:
		fmt.Printf("UNKNOWN %s\n", token.Lexeme)
	case Ident:
		fmt.Printf("IDENT %s\n", token.Lexeme)
	case Plus:
		fmt.Println("PLUS")
	case Minus:
		fmt.Println("MINUS")
	case Star:
		fmt.Println("STAR")
	case Slash:
		fmt.Println("SLASH")
	case Equal:
		fmt.Println("EQUAL")
	case Stop:
		fmt.Println("STOP")
	case Int:
		fmt.Println("INT")
	case Float:
		fmt.Println("FLOAT")
	case Char:
		fmt.Println("CHAR")
	case String:
		fmt.Println("STRING")
	case Bool:
		fmt.Println("BOOL")
	case IntLiteral:
		fmt.Printf("INT LITERAL %s\n", token.Lexeme)
	case FloatLiteral:
		fmt.Printf("FLOAT LITERAL %s\n", token.Lexeme)
	case CharLiteral:
		fmt.Printf("CHAR LITERAL '%s'\n", token.Lexeme)
	case StringLiteral:
		fmt.Printf("STRING LITERAL \"%s\"\n", token.Lexeme)
	case BoolLiteral:
		fmt.Printf("BOOL LITERAL %s\n", token.Lexeme)
	default:
		fmt.Println("UNKNOWN")
		return
	}
}
