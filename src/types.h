#ifndef RAINLANG_TYPES_H
#define RAINLANG_TYPES_H

#include "tokentype.h"

typedef struct Token {
    enum TokenType type;
    char * lexeme;
    void * literal;
    int line;
} Token;

typedef struct TokenList {
    Token * tokens;
    int length;
} TokenList;

typedef struct Scanner {
    const char * source;
    int size;
    TokenList * tokens;
    int start;
    int current;
    int line;
} Scanner;

#endif //RAINLANG_TYPES_H