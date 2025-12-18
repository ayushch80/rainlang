#include "ast.h"

#include <stdio.h>
#include <stdlib.h>
#include <string.h>

AST * ast_list_new_with_kind(ASTKind kind) {
    AST * n = malloc(sizeof(AST));
    n->kind = kind;
    n->list.items = NULL;
    n->list.count = 0;
    n->list.capacity = 0;
    return n;
}

void ast_list_push(AST * list, AST * stmt) {
    if (list->list.count == list->list.capacity) {
        list->list.capacity = list->list.capacity ? list->list.capacity * 2 : 4;
        list->list.items = realloc(
            list->list.items,
            list->list.capacity * sizeof(AST *)
        );
    }
    list->list.items[list->list.count++] = stmt;
}

AST * ast_stmt_list_new(void) {
    return ast_list_new_with_kind(AST_SEQ);
}

void ast_stmt_list_push(AST * list, AST * arg) {
    ast_list_push(list, arg);
}

AST * ast_args_list_new(void) {
    return ast_list_new_with_kind(AST_ARGS_LIST);
}

void ast_args_list_push(AST * list, AST * arg) {
    ast_list_push(list, arg);
}

AST * ast_params_list_new(void) {
    return ast_list_new_with_kind(AST_PARAMS);
}

void ast_params_list_push(AST * list, AST * arg) {
    ast_list_push(list, arg);
}

AST * ast_ident(Token * ident) {
    AST * n = (AST *) malloc(sizeof(AST));
    n->kind = AST_IDENT;
    
    n->ident.name = ident;
    return n;
}

AST * ast_int(Token * tok) {
    AST * n = malloc(sizeof(AST));
    n->kind = AST_INT;

    n->_int = *(long *) tok->literal;
    return n;
}

AST * ast_decimal(Token * tok) {
    AST * n = malloc(sizeof(AST));
    n->kind = AST_DECIMAL;

    n->_decimal = *(double *) tok->literal;;
    return n;
}

AST * ast_char(Token * tok) {
    AST * n = malloc(sizeof(AST));
    n->kind = AST_CHAR;

    n->_char = *(char *) tok->literal;
    return n;
}

AST * ast_string(Token * tok) {
    AST * n = malloc(sizeof(AST));
    n->kind = AST_STRING;

    char *src = (char *) tok->literal;
    char *dst = malloc(strlen(src) + 1);
    strcpy(dst, src);

    n->string = dst;
    return n;
}

AST * ast_nope(void) {
    AST * n = malloc(sizeof(AST));
    n->kind = AST_NOPE;
    n->nope = NULL;

    return n;
}

AST * ast_bool(Token * tok) {
    AST * n = malloc(sizeof(AST));
    n->kind = AST_BOOL;

    n->_bool = tok->type == TRUE ? true : false;
    return n;
}

AST * ast_binary(Token * op, AST * l, AST * r) {
    AST * n = (AST *) malloc(sizeof(AST));
    n->kind = AST_BINARY;
    
    n->binary.op = op;
    n->binary.left = l;
    n->binary.right = r;
    return n;
}

AST * ast_unary(Token * op, AST * expr) {
    AST * n = (AST *) malloc(sizeof(AST));
    n->kind = AST_UNARY;
    
    n->unary.op = op;
    n->unary.expr = expr;
    return n;
}

AST * ast_specifier_var(void) {
    AST * n = (AST *) malloc(sizeof(AST));
    n->kind = AST_VAR;
    
    return n;
}

AST * ast_specifier_const(void) {
    AST * n = (AST *) malloc(sizeof(AST));
    n->kind = AST_CONST;
    
    return n;
}

AST * ast_decl(AST * specifier, AST * name, AST * type) {
    return ast_decl_expr(specifier, name, type, NULL);
}

AST * ast_decl_expr(AST * specifier, AST * name, AST * type, AST * expr) {
    AST * n = (AST *) malloc(sizeof(AST));
    n->kind = AST_DECL;
    
    n->decl.specifier = specifier;
    n->decl.name = name;
    n->decl.type = type;
    n->decl.expr = expr;
    return n;
}

AST * ast_assign(AST * lhs, AST * rhs) {
    AST * n = (AST *) malloc(sizeof(AST));
    n->kind = AST_ASSIGN;
    
    n->assign.lhs = lhs;
    n->assign.rhs = rhs;
    return n;
}

AST * ast_empty_stmt(void) {
    AST * n = (AST *) malloc(sizeof(AST));
    n->kind = AST_EMPTY;
    
    return n;
}

AST * ast_fun_call(AST * name, AST * argslist) {
    AST * n = (AST *) malloc(sizeof(AST));
    n->kind = AST_CALL;
    
    n->call.name = name;
    n->call.argslist = argslist;
    return n;
}

AST * ast_return(AST * arg) {
    AST * n = (AST *) malloc(sizeof(AST));
    n->kind = AST_RETURN;
    
    n->arg = arg;
    return n;
}

AST * ast_function(AST * ident, AST * params, AST * type, AST * body) {
    AST * n = (AST *) malloc(sizeof(AST));
    n->kind = AST_FUNCTION;
    
    n->function.name = ident;
    n->function.params = params;
    n->function.type = type;
    n->function.body = body;
    return n;
}

AST * ast_if_stmt(AST * cond, AST * body) {
    return ast_if_stmt_else(cond, body, NULL);
}

AST * ast_if_stmt_else(AST * cond, AST * body, AST * _else) {
    AST * n = (AST *) malloc(sizeof(AST));
    n->kind = AST_IF_ELSE;
    
    n->if_else.cond = cond;
    n->if_else.body = body;
    n->if_else._else = _else;
    return n;
}

AST * ast_else_stmt(AST * body) {
    return ast_if_stmt_else(NULL, body, NULL);
}

AST * ast_for_loop(AST * decl, AST * cond, AST * stmt, AST * body) {
    AST * n = (AST *) malloc(sizeof(AST));
    n->kind = AST_FOR_LOOP;
    
    n->for_loop.decl = decl;
    n->for_loop.cond = cond;
    n->for_loop.stmt = stmt;
    n->for_loop.body = body;
    return n;
}

AST * ast_while_loop(AST * cond, AST * body) {
    AST * n = (AST *) malloc(sizeof(AST));
    n->kind = AST_WHILE_LOOP;
    
    n->while_loop.cond = cond;
    n->while_loop.body = body;
    return n;
}

AST * ast_param(AST * name, AST * type) {
    AST * n = (AST *) malloc(sizeof(AST));
    n->kind = AST_PARAM;
    
    n->param.name = name;
    n->param.type = type;
    return n;
}

AST * ast_type_int(void) {
    AST * n = (AST *) malloc(sizeof(AST));
    n->kind = AST_TYPE_INT;
    
    return n;
}

AST * ast_type_decimal(void) {
    AST * n = (AST *) malloc(sizeof(AST));
    n->kind = AST_TYPE_DECIMAL;
    
    return n;
}

AST * ast_type_string(void) {
    AST * n = (AST *) malloc(sizeof(AST));
    n->kind = AST_TYPE_STRING;
    
    return n;
}

AST * ast_type_char(void) {
    AST * n = (AST *) malloc(sizeof(AST));
    n->kind = AST_TYPE_CHAR;
    
    return n;
}

AST * ast_type_bool(void) {
    AST * n = (AST *) malloc(sizeof(AST));
    n->kind = AST_TYPE_BOOL;
    
    return n;
}

void ast_print(AST *n, int indent) {
    if (!n) return;

    for (int i = 0; i < indent; i++)
        printf("  ");

    switch (n->kind) {
    case AST_IDENT:
        printf("IDENT %s\n", n->ident.name->lexeme);
        break;

    case AST_INT:
        printf("INT %ld\n", n->_int);
        break;

    case AST_DECIMAL:
        printf("DECIMAL %f\n", n->_decimal);
        break;

    case AST_STRING:
        printf("STRING \"%s\"\n", n->string);
        break;

    case AST_CHAR:
        printf("CHAR '%c'\n", n->_char);
        break;

    case AST_NOPE:
        printf("NOPE\n");
        break;

    case AST_BOOL:
        printf("%s\n", n->_bool ? "TRUE" : "FALSE");
        break;

    case AST_BINARY:
        printf("BINOP %s\n", n->binary.op->lexeme);
        ast_print(n->binary.left, indent + 1);
        ast_print(n->binary.right, indent + 1);
        break;

    case AST_UNARY:
        printf("UNOP %s\n", n->unary.op->lexeme);
        ast_print(n->unary.expr, indent + 1);
        break;

    case AST_VAR:
        printf("VARIABLE\n");
        break;

    case AST_CONST:
        printf("CONSTANT\n");
        break;

    case AST_DECL:
        printf("DECL\n");
        ast_print(n->decl.name, indent + 1);
        ast_print(n->decl.specifier, indent + 1);
        ast_print(n->decl.type, indent + 1);
        if (n->decl.expr != NULL) {
            ast_print(n->decl.expr, indent + 1);
        }
        break;

    case AST_ASSIGN:
        printf("ASSIGN\n");
        ast_print(n->assign.lhs, indent + 1);
        ast_print(n->assign.rhs, indent + 1);
        break;

    case AST_SEQ:
        printf("SEQ\n");
        for (size_t i = 0; i < n->list.count; i++) {
            ast_print(n->list.items[i], indent + 1);
        }
        break;

    case AST_EMPTY:
        printf("EMPTY\n");
        break;

    case AST_CALL:
        printf("CALL\n");
        ast_print(n->call.name, indent + 1);
        ast_print(n->call.argslist, indent + 1);
        break;

    case AST_RETURN:
        printf("RETURN\n");
        ast_print(n->arg, indent + 1);
        break;

    case AST_FUNCTION:
        printf("FUNCTION\n");
        ast_print(n->function.name, indent + 1);
        ast_print(n->function.params, indent + 1);
        ast_print(n->function.type, indent + 1);
        ast_print(n->function.body, indent + 1);
        break;

    case AST_IF_ELSE:
        // if (n == NULL) break;
        if (n->if_else.cond == NULL) {
            printf("ELSE\n");
        } else {
            printf("IF\n");
            ast_print(n->if_else.cond, indent + 1);
        }
        ast_print(n->if_else.body, indent + 1);
        ast_print(n->if_else._else, indent + 1);
        break;

    case AST_FOR_LOOP:
        printf("FOR\n");
        ast_print(n->for_loop.decl, indent + 1);
        ast_print(n->for_loop.cond, indent + 1);
        ast_print(n->for_loop.stmt, indent + 1);
        ast_print(n->for_loop.body, indent + 1);
        break;

    case AST_WHILE_LOOP:
        printf("WHILE\n");
        ast_print(n->while_loop.cond, indent + 1);
        ast_print(n->while_loop.body, indent + 1);
        break;

    case AST_PARAM:
        printf("PARAM\n");
        ast_print(n->param.name, indent + 1);
        ast_print(n->param.type, indent + 1);
        break;

    case AST_ARGS_LIST:
        if (n->list.count > 0) {
            printf("ARGS\n");
            for (size_t i = 0; i < n->list.count; i++) {
                ast_print(n->list.items[i], indent + 1);
            }
        }
        break;

    case AST_PARAMS:
        if (n->list.count > 0) {
            printf("PARAMS\n");
            for (size_t i = 0; i < n->list.count; i++) {
                ast_print(n->list.items[i], indent + 1);
            }
        }
        break;

    case AST_TYPE_INT:
        printf("TYPE int\n");
        break;

    case AST_TYPE_DECIMAL:
        printf("TYPE decimal\n");
        break;

    case AST_TYPE_BOOL:
        printf("TYPE bool\n");
        break;

    case AST_TYPE_CHAR:
        printf("TYPE char\n");
        break;

    case AST_TYPE_STRING:
        printf("TYPE string\n");
        break;
    }
}


void ast_free(AST * n) {
    if (!n) return;

    switch (n->kind) {
    case AST_INT:
    case AST_DECIMAL:
    case AST_BOOL:
    case AST_CHAR:
        break;

    case AST_STRING:
        free(n->string);
        break;

    case AST_BINARY:
        ast_free(n->binary.left);
        ast_free(n->binary.right);
        break;

    case AST_UNARY:
        ast_free(n->unary.expr);
        break;

    case AST_DECL:
        ast_free(n->decl.specifier);
        ast_free(n->decl.name);
        ast_free(n->decl.type);
        ast_free(n->decl.expr);
        break;

    case AST_ASSIGN:
        ast_free(n->assign.lhs);
        ast_free(n->assign.rhs);
        break;

    case AST_CALL:
        ast_free(n->call.name);
        ast_free(n->call.argslist);
        break;

    case AST_PARAM:
        ast_free(n->param.name);
        ast_free(n->param.type);
        break;

    case AST_RETURN:
        ast_free(n->arg);
        break;

    case AST_FUNCTION:
        ast_free(n->function.name);
        ast_free(n->function.params);
        ast_free(n->function.type);
        ast_free(n->function.body);
        break;

    case AST_IF_ELSE:
        ast_free(n->if_else.cond);
        ast_free(n->if_else.body);
        ast_free(n->if_else._else);
        break;

    case AST_FOR_LOOP:
        ast_free(n->for_loop.decl);
        ast_free(n->for_loop.cond);
        ast_free(n->for_loop.stmt);
        ast_free(n->for_loop.body);
        break;

    case AST_WHILE_LOOP:
        ast_free(n->while_loop.cond);
        ast_free(n->while_loop.body);
        break;

    case AST_SEQ:
    case AST_PARAMS:
    case AST_ARGS_LIST:
        for (size_t i = 0; i < n->list.count; i++)
            ast_free(n->list.items[i]);
        free(n->list.items);
        break;

    default:
        break;
    }

    free(n);
}