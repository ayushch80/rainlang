#include "scanner.h"

#include <ctype.h>
#include <stdlib.h>
#include <stdbool.h>
#include <stdio.h>
#include <string.h>

#include "main.h"
#include "tokentype.h"
#include "types.h"

bool isAtEnd(const Scanner * scanner) {
    return scanner->current >= scanner->size;
}

void addToken(const Scanner * scanner, enum TokenType tokenType, void * literal) {
    int size = scanner->current-scanner->start;
    char * text = (char*) malloc((size+1) * sizeof(char));
    memcpy(text, scanner->source + scanner->start, size);
    text[size] = '\0';

    Token token;
    token.type = tokenType;
    token.lexeme = text;
    token.literal = literal;
    token.line = scanner->line;

    if (scanner->tokens->length >= MAX_TOKENS) {
        fprintf(stderr, "Max Token limit reached (%d)\n", MAX_TOKENS);
        exit(66);
    }

    scanner->tokens->tokens[scanner->tokens->length++] = token;
}

bool match(Scanner * scanner, char expected) {
    if (isAtEnd(scanner)) return false;
    if (scanner->source[scanner->current] != expected) return false;

    scanner->current++;
    return true;
}

char advance(Scanner * scanner) {
    return scanner->source[scanner->current++];
}

char peek(const Scanner * scanner) {
    if (isAtEnd(scanner)) return '\0';
    return scanner->source[scanner->current];
}

char peekNext(const Scanner * scanner) {
    if (scanner->current + 1 >= scanner->size) return '\0';
    return scanner->source[scanner->current + 1];
}

void character(Scanner * scanner) {
    while (peek(scanner) != '\'' && !isAtEnd(scanner)) {
        if (peek(scanner) == '\n') scanner->line++;
        advance(scanner);
    }

    if (isAtEnd(scanner)) {
        error(scanner->line, "Unterminated character literal.");
        exit(67);
    }

    advance(scanner);

    const int size = scanner->current-scanner->start-2;

    if (size != 1) {
        error(scanner->line, "Invalid character literal.");
        exit(67);
    }

    char * value = (char*) malloc(sizeof(char)*2);
    value[0] = scanner->source[scanner->start + 1];
    value[1] = '\0';

    addToken(scanner, CHAR_LITERAL, (void *)value);
}

void string(Scanner * scanner) {
    while (peek(scanner) != '"' && !isAtEnd(scanner)) {
        if (peek(scanner) == '\n') scanner->line++;
        advance(scanner);
    }

    if (isAtEnd(scanner)) {
        error(scanner->line, "Unterminated string.");
        exit(67);
    }

    advance(scanner);

    int size = scanner->current-scanner->start - 2;
    char * value = (char*) malloc((size+1) * sizeof(char));
    memcpy(value, scanner->source + scanner->start + 1, size);
    value[size] = '\0';

    addToken(scanner, STRING_LITERAL, (void *)value);
}

void number(Scanner * scanner) {
    enum TokenType type = INT_LITERAL;

    while (isdigit(peek(scanner))) advance(scanner);

    if (peek(scanner) == '.' && isdigit(peekNext(scanner))) {
        type = DECIMAL_LITERAL;
        advance(scanner);

        while (isdigit(peek(scanner))) advance(scanner);
    }

    const int size = scanner->current-scanner->start;
    char * value = (char*) malloc((size+1) * sizeof(char));
    memcpy(value, scanner->source + scanner->start, size);
    value[size] = '\0';

    void * number = NULL;
    char * endptr;

    if (type == INT_LITERAL) {
        number = (void *) malloc(sizeof(long));
        const long num = strtol(value, &endptr, 10);
        memcpy(number, &num, sizeof(long));
    } else {
        number = (void *) malloc(sizeof(double));
        const double num = strtod(value, &endptr);
        memcpy(number, &num, sizeof(double));
    }

    addToken(scanner, type, number);
}

void identifier(Scanner * scanner) {
    while (isalnum(peek(scanner)) || peek(scanner) == '_') advance(scanner);

    const int size = scanner->current-scanner->start;
    char * text = (char*) malloc((size+1) * sizeof(char));
    memcpy(text, scanner->source + scanner->start, size);
    text[size] = '\0';

    if (strcmp(text, "nope") == 0) {
        addToken(scanner, NOPE, NULL);
    } else if (strcmp(text, "if") == 0) {
        addToken(scanner, IF, NULL);
    } else if (strcmp(text, "else") == 0) {
        addToken(scanner, ELSE, NULL);
    } else if (strcmp(text, "while") == 0) {
        addToken(scanner, WHILE, NULL);
    } else if (strcmp(text, "for") == 0) {
        addToken(scanner, FOR, NULL);
    } else if (strcmp(text, "return") == 0) {
        addToken(scanner, RETURN, NULL);
    } else if (strcmp(text, "var") == 0) {
        addToken(scanner, VAR, NULL);
    } else if (strcmp(text, "const") == 0) {
        addToken(scanner, CONST, NULL);
    } else if (strcmp(text, "true") == 0) {
        addToken(scanner, TRUE, NULL);
    } else if (strcmp(text, "false") == 0) {
        addToken(scanner, FALSE, NULL);
    } else if (strcmp(text, "fn") == 0) {
        addToken(scanner, FN, NULL);
    } else if (strcmp(text, "str") == 0) {
        addToken(scanner, STRING, NULL);
    } else if (strcmp(text, "char") == 0) {
        addToken(scanner, CHAR, NULL);
    } else if (strcmp(text, "deci") == 0) {
        addToken(scanner, DECIMAL, NULL);
    } else if (strcmp(text, "int") == 0) {
        addToken(scanner, INT, NULL);
    } else if (strcmp(text, "bool") == 0) {
        addToken(scanner, BOOL, NULL);
    } else {
        addToken(scanner, IDENTIFIER, NULL);
    }

    free(text);
}

void scanToken(Scanner * scanner) {
    char c = advance(scanner);
    switch (c) {
        case '(':
            addToken(scanner, LEFT_PAREN, NULL);
            break;
        case ')':
            addToken(scanner, RIGHT_PAREN, NULL);
            break;
        case '{':
            addToken(scanner, LEFT_BRACE, NULL);
            break;
        case '}':
            addToken(scanner, RIGHT_BRACE, NULL);
            break;
        case '[':
            addToken(scanner, LEFT_BRACKET, NULL);
            break;
        case ']':
            addToken(scanner, RIGHT_BRACKET, NULL);
            break;
        case ',':
            addToken(scanner, COMMA, NULL);
            break;
        case '.':
            addToken(scanner, DOT, NULL);
            break;
        case '-':
            addToken(scanner, MINUS, NULL);
            break;
        case '+':
            addToken(scanner, PLUS, NULL);
            break;
        case ';':
            addToken(scanner, SEMICOLON, NULL);
            break;
        case '*':
            addToken(scanner, STAR, NULL);
            break;
        case '%':
            addToken(scanner, PERCENT, NULL);
            break;
        case '@':
            addToken(scanner, AT, NULL);
            break;
        case '!':
            addToken(scanner, match(scanner, '=') ? BANG_EQUAL : BANG, NULL);
            break;
        case '=':
            addToken(scanner, match(scanner, '=') ? EQUAL_EQUAL : EQUAL, NULL);
            break;
        case '<':
            addToken(scanner, match(scanner, '=') ? LESS_EQUAL : LESS, NULL);
            break;
        case '>':
            addToken(scanner, match(scanner, '=') ? GREATER_EQUAL : GREATER, NULL);
            break;
        case '&':
            if (match(scanner, '&')) {
                addToken(scanner, LOGICAL_AND, NULL);
            } else error(scanner->line, "Unexpected '&'");
            break;
        case '|':
            if (match(scanner, '|')) {
                addToken(scanner, LOGICAL_OR, NULL);
            } else error(scanner->line, "Unexpected '|'");
            break;
        case '/':
            if (match(scanner, '/')) {
                while (peek(scanner) != '\n' && !isAtEnd(scanner)) advance(scanner);
            } else {
                addToken(scanner, SLASH, NULL);
            }
            break;
        case ' ':
        case '\r':
        case '\t':
            break;
        case '\n':
            scanner->line++;
            break;
        case '\'':
            character(scanner);
            break;
        case '"':
            string(scanner);
            break;
        default:
            if (isdigit(c)) {
                number(scanner);
                break;
            } else if (isalpha(c) || c == '_') {
                identifier(scanner);
                break;
            }
            error(scanner->line, "Unexpected character.");
            break;
    }
}

void scanTokens(Scanner * scanner) {
    TokenList * tokenList = (TokenList *) malloc(sizeof(TokenList));
    scanner->tokens = tokenList;
    scanner->tokens->length = 0;
    scanner->tokens->tokens = (Token *) malloc(sizeof(Token) * MAX_TOKENS);

    // printf("Size : %d\n", scanner->size);

    while (!isAtEnd(scanner)) {
        scanner->start = scanner->current;
        scanToken(scanner);
    }

    Token EOFToken;
    EOFToken.type = eOF;
    EOFToken.lexeme = "";
    EOFToken.literal = NULL;
    EOFToken.line = scanner->line;

    scanner->tokens->tokens[scanner->tokens->length++] = EOFToken;
}
