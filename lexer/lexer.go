package lexer

import (
	"fmt"
	"strings"
	"unicode"
)

type TokenType int

const (
	EOF TokenType = iota
	Unknown

	Ident

	Plus
	Minus
	Star
	Slash
	Equal
	Stop

	IntLiteral
	FloatLiteral
	CharLiteral
	StringLiteral
	BoolLiteral

	Int
	Float
	Char
	String
	Bool
)

func (t TokenType) String() string {
	switch t {
	case EOF:
		return "EOF"
	case Unknown:
		return "UNKNOWN"
	case Ident:
		return "IDENT"
	case Plus:
		return "PLUS"
	case Minus:
		return "MINUS"
	case Star:
		return "STAR"
	case Slash:
		return "SLASH"
	case Equal:
		return "EQUAL"
	case Stop:
		return "STOP"
	case Int:
		return "INT"
	case Float:
		return "FLOAT"
	case Char:
		return "CHAR"
	case String:
		return "STRING"
	case Bool:
		return "BOOL"
	case IntLiteral:
		return "INT LITERAL"
	case FloatLiteral:
		return "FLOAT LITERAL"
	case CharLiteral:
		return "CHAR LITERAL"
	case StringLiteral:
		return "STRING LITERAL"
	case BoolLiteral:
		return "BOOL LITERAL"
	default:
		return "UNKNOWN"
	}
}

type Token struct {
	Type   TokenType
	Lexeme string
	Line   int
	Col    int
}

func Tokenizer(code string) ([]Token, error) {
	currLine := 0

	var tokens []Token
	lines := strings.Split(code, "\n")

	for _, line := range lines {
		currLine++

		for col := 0; col < len(line); {
			ch := line[col]
			col++

			token := Token{
				Type: Unknown,
				Line: currLine,
				Col:  col,
			}

			// fmt.Println(match(line, col, "true"), line[col-1:])

			// IGNORE WHITESPACE
			if unicode.IsSpace(rune(ch)) {
				continue
			} else

			// SINGLE CHAR TOKENS
			if ch == '+' {
				token.Type = Plus
				token.Lexeme = "+"
			} else if ch == '-' {
				token.Type = Minus
				token.Lexeme = "-"
			} else if ch == '*' {
				token.Type = Star
				token.Lexeme = "*"
			} else if ch == '/' {
				token.Type = Slash
				token.Lexeme = "/"
			} else if ch == '=' {
				token.Type = Equal
				token.Lexeme = "="
			} else if ch == ';' {
				token.Type = Stop
				token.Lexeme = ";"
			} else

			// CHAR AND STRING
			if ch == '\'' {
				stop := findNext(line, col, '\'')
				if stop == -1 {
					return tokens, fmt.Errorf("[%d:%d] Unterminated character literal\n", currLine, col)
				}

				token.Type = CharLiteral
				token.Lexeme = strings.Clone(line[col:stop])

				if len(token.Lexeme) != 1 {
					return tokens, fmt.Errorf("[%d:%d] Invalid character literal\n", currLine, col)
				}

				col = stop + 1
			} else if ch == '"' {
				stop := findNext(line, col, '"')
				if stop == -1 {
					return tokens, fmt.Errorf("[%d:%d] Unterminated string literal\n", currLine, col)
				}

				token.Type = StringLiteral
				token.Lexeme = strings.Clone(line[col:stop])

				col = stop + 1
			} else

			// INT, FLOAT and BOOL values
			if match(line, col, "true") {
				token.Type = BoolLiteral
				token.Lexeme = "true"
				col += 3
			} else if match(line, col, "false") {
				token.Type = BoolLiteral
				token.Lexeme = "false"
				col += 4
			} else if unicode.IsDigit(rune(ch)) {
				token.Type = IntLiteral
				token.Lexeme = string(ch)
				if col < len(line) && line[col] == '.' {
					token.Lexeme += "."
					token.Type = FloatLiteral
					col++
					if col < len(line) && !unicode.IsDigit(rune(line[col])) {
						return tokens, fmt.Errorf("[%d:%d] Invalid literal\n", currLine, col)
					}
				}
				for col < len(line) {
					if unicode.IsDigit(rune(line[col])) {
						token.Lexeme += string(line[col])
						col++
						if col < len(line) && line[col] == '.' {
							if token.Type == IntLiteral {
								token.Lexeme += "."
								token.Type = FloatLiteral
								col++
							} else {
								return tokens, fmt.Errorf("[%d:%d] Invalid literal\n", currLine, col)
							}
						}
					} else {
						break
					}
				}
			} else

			// IDENTIFIERS and types
			if unicode.IsLetter(rune(ch)) {
				if match(line, col, "int") {
					token.Type = Int
					token.Lexeme = "int"
					col += 3
				} else if match(line, col, "float") {
					token.Type = Float
					token.Lexeme = "float"
					col += 5
				} else if match(line, col, "char") {
					token.Type = Char
					token.Lexeme = "char"
					col += 4
				} else if match(line, col, "string") {
					token.Type = String
					token.Lexeme = "string"
					col += 6
				} else if match(line, col, "bool") {
					token.Type = Bool
					token.Lexeme = "bool"
					col += 4
				} else {
					token.Type = Ident
					token.Lexeme += string(ch)
					for col < len(line) && (unicode.IsLetter(rune(line[col])) || unicode.IsDigit(rune(line[col])) || line[col] == '_') {
						token.Lexeme += string(line[col])
						col++
					}
				}
			}

			tokens = append(tokens, token)
		}
	}

	tokens = append(tokens, Token{
		Type:   EOF,
		Lexeme: "EOF",
		Line:   currLine + 1,
		Col:    0,
	})

	return tokens, nil
}

func match(line string, col int, target string) bool {
	var i int
	for i = col - 1; i < len(line); i++ {
		ch := line[i]
		// if i - col + 1 >= len(target) {
		// 	if line[i]
		// 	break
		// }
		if unicode.IsSpace(rune(ch)) || (!unicode.IsLetter(rune(ch)) && !unicode.IsDigit(rune(ch))) {
			break
		}
	}

	// fmt.Println("HERERR", target, line[col-1 : i], col-1, i)

	str := strings.Clone(line[col-1 : i])
	if str == target {
		return true
	}
	return false
}

func findNext(line string, col int, ch byte) int {
	for i := col + 1; i < len(line); i++ {
		if line[i] == ch {
			return i
		}
	}
	return -1
}
