%define api.token.prefix {YY_}
%start program

%{
#include <stdio.h>
#include <stdlib.h>
#include "ast.h"
#include "token.h"

TokenList *g_tokens;
int pos = 0;

AST *root;

int yylex(void);
void yyerror(const char *s);
%}

%union {
    Token *tok;
    AST   *node;
}

%token
    LEFT_PAREN RIGHT_PAREN LEFT_BRACE RIGHT_BRACE
    COMMA DOT SEMICOLON AT
    STRING CHAR DECIMAL INT BOOL
    IF ELSE WHILE FOR RETURN VAR CONST FN NOPE
    eOF

%token <tok>
    PLUS MINUS STAR SLASH PERCENT
    IDENTIFIER
    BANG BANG_EQUAL
    EQUAL EQUAL_EQUAL
    GREATER GREATER_EQUAL
    LESS LESS_EQUAL
    LOGICAL_AND LOGICAL_OR
    TRUE FALSE
    STRING_LITERAL CHAR_LITERAL DECIMAL_LITERAL INT_LITERAL

%type <node>
    expr mul_expr factor stmt_list stmt
    decl specifier type fun_call args_list
    arg ident return_stmt or_expr and_expr
    equality comparison add_expr bool
    unary fun params param if_stmt else_stmt
    for_stmt while_stmt assign

%%

program:
      stmt_list eOF
        { root = $1; }
;

ident:
      IDENTIFIER
        { $$ = ast_ident($1); }
;

fun_call:
      ident AT args_list
        { $$ = ast_fun_call($1, $3); }
;

args_list:
      /* empty */
        { $$ = ast_args_list_new(); }
    | arg
        {
            $$ = ast_args_list_new();
            ast_args_list_push($$, $1);
        }
    | args_list COMMA arg
        {
            ast_args_list_push($1, $3);
            $$ = $1;
        }
;

arg:
      expr
        { $$ = $1; }
    | fun_call
        { $$ = $1; }
    | LEFT_PAREN arg RIGHT_PAREN
        { $$ = $2; }
;

return_stmt:
      RETURN arg
        { $$ = ast_return($2);}
;

stmt_list:
      /* empty */
        { $$ = ast_stmt_list_new(); }
    | stmt
        {
            $$ = ast_stmt_list_new();
            ast_stmt_list_push($$, $1);
        }
    | stmt_list stmt
        {
            ast_stmt_list_push($1, $2);
            $$ = $1;
        }
;

stmt:
      expr SEMICOLON
        { $$ = $1; }
    | decl SEMICOLON
        { $$ = $1; }
    | assign SEMICOLON
        { $$ = $1; }
    | fun_call SEMICOLON
        { $$ = $1; }
    | return_stmt SEMICOLON
        { $$ = $1; }
    | fun
        { $$ = $1; }
    | if_stmt
        { $$ = $1; }
    | for_stmt
        { $$ = $1; }
    | while_stmt
        { $$ = $1; }
;

assign:
      ident EQUAL arg
        { $$ = ast_assign($1, $3); }
;

for_stmt:
      FOR decl SEMICOLON expr SEMICOLON expr LEFT_BRACE stmt_list RIGHT_BRACE
        { $$ = ast_for_loop($2, $4, $6, $8); }
;

while_stmt:
      WHILE expr LEFT_BRACE stmt_list RIGHT_BRACE
        { $$ = ast_while_loop($2, $4); }
;

if_stmt:
      IF expr LEFT_BRACE stmt_list RIGHT_BRACE else_stmt
        { $$ = ast_if_stmt_else($2, $4, $6); }
    | IF expr LEFT_BRACE stmt_list RIGHT_BRACE
        { $$ = ast_if_stmt($2, $4); }
;

else_stmt:
      ELSE if_stmt
        { $$ = $2; }
    | ELSE LEFT_BRACE stmt_list RIGHT_BRACE
        { $$ = ast_else_stmt($3); }
;

fun:
      FN ident LEFT_PAREN params RIGHT_PAREN type LEFT_BRACE stmt_list RIGHT_BRACE
        { $$ = ast_function($2, $4, $6, $8); }
;

params:
      /* empty */
        { $$ = ast_params_list_new(); }
    | param
        {
            $$ = ast_params_list_new();
            ast_params_list_push($$, $1);
        }
    | params COMMA param
        {
            ast_params_list_push($1, $3);
            $$ = $1;
        }
;

param:
    ident type
        { $$ = ast_param($1, $2); }
;

decl:
      specifier ident type
        { $$ = ast_decl($1, $2, $3); }
    | specifier ident type EQUAL expr
        { $$ = ast_decl_expr($1, $2, $3, $5); }
;

specifier:
      VAR
        { $$ = ast_specifier_var(); }
    | CONST
        { $$ = ast_specifier_const(); }
;

type:
      INT
        { $$ = ast_type_int(); }
    | DECIMAL
        { $$ = ast_type_decimal(); }
    | BOOL
        { $$ = ast_type_bool(); }
    | CHAR
        { $$ = ast_type_char(); }
    | STRING
        { $$ = ast_type_string(); }
;

expr:
      or_expr
        { $$ = $1; }
;

or_expr:
      or_expr LOGICAL_OR and_expr
        { $$ = ast_binary($2, $1, $3); }
    | and_expr
        { $$ = $1; }
;

and_expr:
      and_expr LOGICAL_AND equality
        { $$ = ast_binary($2, $1, $3); }
    | equality
        { $$ = $1; }
;

equality:
    equality EQUAL_EQUAL comparison
        { $$ = ast_binary($2, $1, $3); }
    | equality BANG_EQUAL comparison
        { $$ = ast_binary($2, $1, $3); }
    | comparison
        { $$ = $1; }
;

comparison:
    comparison LESS add_expr
        { $$ = ast_binary($2, $1, $3); }
    | comparison LESS_EQUAL add_expr
        { $$ = ast_binary($2, $1, $3); }
    | comparison GREATER add_expr
        { $$ = ast_binary($2, $1, $3); }
    | comparison GREATER_EQUAL add_expr
        { $$ = ast_binary($2, $1, $3); }
    | add_expr
        { $$ = $1; }
;

add_expr:
      add_expr PLUS mul_expr
        { $$ = ast_binary($2, $1, $3); }
    | add_expr MINUS mul_expr
        { $$ = ast_binary($2, $1, $3); }
    | mul_expr
        { $$ = $1; }
;

mul_expr:
      mul_expr STAR unary
        { $$ = ast_binary($2, $1, $3); }
    | mul_expr SLASH unary
        { $$ = ast_binary($2, $1, $3); }
    | mul_expr PERCENT unary
        { $$ = ast_binary($2, $1, $3); }
    | unary
        { $$ = $1; }
;

unary:
      BANG unary
        { $$ = ast_unary($1, $2); }
    | MINUS unary
        { $$ = ast_unary($1, $2); }
    | factor
        { $$ = $1; }
;

bool:
      TRUE
        { $$ = ast_bool($1); }
    | FALSE
        { $$ = ast_bool($1); }
;

factor:
      fun_call
        { $$ = $1; }
    | INT_LITERAL
        { $$ = ast_int($1); }
    | DECIMAL_LITERAL
        { $$ = ast_decimal($1); }
    | CHAR_LITERAL
        { $$ = ast_char($1); }
    | STRING_LITERAL
        { $$ = ast_string($1); }
    | bool
        { $$ = $1; }
    | NOPE
        { $$ = ast_nope(); }
    | ident
        { $$ = $1; }
    | LEFT_PAREN expr RIGHT_PAREN
        { $$ = $2; }
;
%%

int yylex(void) {
    if (pos >= g_tokens->length)
        return 0;

    Token *t = &g_tokens->tokens[pos++];
    yylval.tok = t;

    switch (t->type) {
        case LEFT_PAREN:        return YY_LEFT_PAREN;
        case RIGHT_PAREN:       return YY_RIGHT_PAREN;
        case LEFT_BRACE:        return YY_LEFT_BRACE;
        case RIGHT_BRACE:       return YY_RIGHT_BRACE;

        case SEMICOLON:         return YY_SEMICOLON;
        case COMMA:             return YY_COMMA;
        case PLUS:              return YY_PLUS;
        case MINUS:             return YY_MINUS;
        case STAR:              return YY_STAR;
        case SLASH:             return YY_SLASH;
        case PERCENT:           return YY_PERCENT;
        case AT:                return YY_AT;

        case BANG:              return YY_BANG;
        case BANG_EQUAL:        return YY_BANG_EQUAL;
        case EQUAL:             return YY_EQUAL;
        case EQUAL_EQUAL:       return YY_EQUAL_EQUAL;
        case GREATER:           return YY_GREATER;
        case GREATER_EQUAL:     return YY_GREATER_EQUAL;
        case LESS:              return YY_LESS;
        case LESS_EQUAL:        return YY_LESS_EQUAL;
        case LOGICAL_AND:       return YY_LOGICAL_AND;
        case LOGICAL_OR:        return YY_LOGICAL_OR;

        case INT:               return YY_INT;
        case DECIMAL:           return YY_DECIMAL;
        case BOOL:              return YY_BOOL;
        case STRING:            return YY_STRING;
        case CHAR:              return YY_CHAR;
        case NOPE:              return YY_NOPE;
        case TRUE:              return YY_TRUE;
        case FALSE:             return YY_FALSE;

        case INT_LITERAL:       return YY_INT_LITERAL;
        case DECIMAL_LITERAL:   return YY_DECIMAL_LITERAL;
        case STRING_LITERAL:    return YY_STRING_LITERAL;
        case CHAR_LITERAL:      return YY_CHAR_LITERAL;
        case IDENTIFIER:        return YY_IDENTIFIER;

        case FN:                return YY_FN;
        case IF:                return YY_IF;
        case ELSE:              return YY_ELSE;
        case FOR:               return YY_FOR;
        case WHILE:             return YY_WHILE;
        case VAR:               return YY_VAR;
        case CONST:             return YY_CONST;
        case RETURN:            return YY_RETURN;

        case eOF:               return YY_eOF;

        default:
            fprintf(stderr, "Unknown token: ");
            printToken(* t);
            return YY_YYerror;
    }
}


void yyerror(const char * s) {
    if (pos > 0 && pos <= g_tokens->length) {
        Token * t = &g_tokens->tokens[pos - 1];
        char * str = tokenToString(* t);
        fprintf(stderr, "[line %d] parse error near '%s'\n%s\n", t->line, t->lexeme ? t->lexeme : "<eof>", str);
        free(str);
    } else {
        fprintf(stderr, "parse error\n");
    }
}