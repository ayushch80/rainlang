#ifndef RAINLANG_AST_H
#define RAINLANG_AST_H

#include "types.h"

#include <stddef.h>
#include <stdbool.h>

typedef enum {
    AST_IDENT,
    AST_INT,
    AST_DECIMAL,
    AST_STRING,
    AST_CHAR,
    AST_BOOL,
    AST_NOPE,

    AST_BINARY,
    AST_UNARY,

    AST_VAR,
    AST_CONST,
    
    AST_DECL,
    AST_ASSIGN,
    AST_SEQ,
    AST_EMPTY,
    AST_CALL,
    AST_RETURN,
    AST_FUNCTION,
    AST_IF_ELSE,
    AST_FOR_LOOP,
    AST_WHILE_LOOP,

    AST_PARAM,
    AST_PARAMS,
    AST_ARGS_LIST,

    AST_TYPE_INT,
    AST_TYPE_DECIMAL,
    AST_TYPE_STRING,
    AST_TYPE_CHAR,
    AST_TYPE_BOOL
} ASTKind;

typedef struct AST AST;

typedef struct {
    AST ** items;
    size_t count;
    size_t capacity;
} ASTList;

struct AST {
    ASTKind kind;

    union {

        long _int;
        double _decimal;
        bool _bool;
        char _char;
        char * string;
        void * nope;

        /* AST_INDET */
        struct {
            Token * name;
        } ident;

        /* AST_BINARY */
        struct {
            Token * op;
            AST * left;
            AST * right;
        } binary;

        /* AST_UNARY */
        struct {
            Token * op;
            AST * expr;
        } unary;

        /* AST_DECL */
        struct {
            AST * specifier;
            AST * name;
            AST * type;
            AST * expr;
        } decl;

        /* AST_ASSIGN */
        struct {
            AST * rhs;
            AST * lhs;
        } assign;

        /* AST_SEQ, AST_ARGS_LIST and AST_PARAMS */
        ASTList list;

        /* AST_CALL */
        struct {
            AST * name;
            AST * argslist;
        } call;

        /* AST_PARAM */
        struct {
            AST * name;
            AST * type;
        } param;

        /* AST_RETURN */
        AST * arg;

        /* AST_FUNCTION */
        struct {
            AST * name;
            AST * params;
            AST * type;
            AST * body;
        } function;

        /* AST_IF_ELSE */
        struct {
            AST * cond;
            AST * body;
            AST * _else;
        } if_else;

        /* AST_FOR_LOOP */
        struct {
            AST * decl;
            AST * cond;
            AST * stmt;
            AST * body;
        } for_loop;

        /* AST_WHILE_LOOP */
        struct {
            AST * cond;
            AST * body;
        } while_loop;
    };
};

AST * ast_stmt_list_new(void);
void ast_stmt_list_push(AST * list, AST * stmt);
AST * ast_args_list_new(void);
void ast_args_list_push(AST * list, AST * arg);
AST * ast_params_list_new(void);
void ast_params_list_push(AST * list, AST * arg);

AST * ast_ident(Token * ident);
AST * ast_int(Token * tok);
AST * ast_decimal(Token * tok);
AST * ast_char(Token * tok);
AST * ast_string(Token * tok);
AST * ast_nope(void);
AST * ast_bool(Token * tok);

AST * ast_binary(Token * op, AST * l, AST * r);
AST * ast_unary(Token * op, AST * expr);

AST * ast_specifier_var(void);
AST * ast_specifier_const(void);

AST * ast_decl(AST * specifier, AST * name, AST * type);
AST * ast_decl_expr(AST * specifier, AST * name, AST * type, AST * expr);
AST * ast_assign(AST * lhs, AST * rhs);
AST * ast_empty_stmt(void);
AST * ast_fun_call(AST * name, AST * argslist);
AST * ast_return(AST * arg);
AST * ast_function(AST * ident, AST * params, AST * type, AST * body);
AST * ast_if_stmt(AST * cond, AST * body);
AST * ast_if_stmt_else(AST * cond, AST * body, AST * _else);
AST * ast_else_stmt(AST * body);
AST * ast_for_loop(AST * decl, AST * cond, AST * stmt, AST * body);
AST * ast_while_loop(AST * cond, AST * body);

AST * ast_param(AST * ident, AST * type);

AST * ast_type_int(void);
AST * ast_type_decimal(void);
AST * ast_type_string(void);
AST * ast_type_char(void);
AST * ast_type_bool(void);

void ast_print(AST * n, int indent);
void ast_free(AST * node);

#endif