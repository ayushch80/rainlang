#include "token.h"

#include <stdio.h>
#include <stdlib.h>

const char * tokentype_to_string(enum TokenType t) {
    switch(t) {
        case LEFT_PAREN: return "LEFT_PAREN";
        case RIGHT_PAREN: return "RIGHT_PAREN";
        case LEFT_BRACE: return "LEFT_BRACE";
        case RIGHT_BRACE: return "RIGHT_BRACE";
        case COMMA: return "COMMA";
        case DOT: return "DOT";
        case MINUS: return "MINUS";
        case PLUS: return "PLUS";
        case STAR: return "STAR";
        case SLASH: return "SLASH";
        case PERCENT: return "PERCENT";
        case SEMICOLON: return "SEMICOLON";
        case AT: return "AT";
        case BANG: return "BANG";
        case BANG_EQUAL: return "BANG_EQUAL";
        case EQUAL: return "EQUAL";
        case EQUAL_EQUAL: return "EQUAL_EQUAL";
        case GREATER: return "GREATER";
        case GREATER_EQUAL: return "GREATER_EQUAL";
        case LESS: return "LESS";
        case LESS_EQUAL: return "LESS_EQUAL";
        case IDENTIFIER: return "IDENTIFIER";
        case STRING: return "STRING";
        case CHAR: return "CHAR";
        case DECIMAL: return "DECIMAL";
        case INT: return "INT";
        case BOOL: return "BOOL";
        case STRING_LITERAL: return "STRING_LITERAL";
        case CHAR_LITERAL: return "CHAR_LITERAL";
        case DECIMAL_LITERAL: return "DECIMAL_LITERAL";
        case INT_LITERAL: return "INT_LITERAL";
        case LOGICAL_AND: return "LOGICAL_AND";
        case LOGICAL_OR: return "LOGICAL_OR";
        case NOPE: return "NOPE";
        case IF: return "IF";
        case ELSE: return "ELSE";
        case WHILE: return "WHILE";
        case FOR: return "FOR";
        case RETURN: return "RETURN";
        case VAR: return "VAR";
        case CONST: return "CONST";
        case TRUE: return "TRUE";
        case FALSE: return "FALSE";
        case FN: return "FN";
        case eOF: return "eOF";
        default: return "UNKNOWN";
    }
}

char * tokenToString(Token token) {
    char * string = (char *) malloc(1024);
    sprintf(string, "[%d] %s %s", token.line, tokentype_to_string(token.type), token.lexeme);
    return string;
}

void printToken(Token t) {
    char * string = tokenToString(t);
    puts(string);
    free(string);
}