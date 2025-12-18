#include "main.h"

#include <stdbool.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>

#include "ast.h"
#include "scanner.h"
#include "token.h"
#include "types.h"

static bool debug = false;
static bool debugScanner = false;
static bool debugParser = false;
static bool hadError = false;

extern AST *root;
extern int pos;
int yyparse(void);

void report(int line, const char * where, const char * message) {
    printf("[line %d] Error %s : %s\n", line, where, message);
}

void error(int line, const char * message) {
    report(line, "", message);
}

void run(const char *buffer, const int size) {
    // printf("%s\n", buffer);

    Scanner * scanner = (Scanner *) malloc(sizeof(Scanner));
    scanner->size = size;
    scanner->source = buffer;
    scanner->start = 0;
    scanner->current = 0;
    scanner->line = 1;

    scanTokens(scanner);

    if (debug || debugScanner) {
        for (int i = 0; i < scanner->tokens->length; i++) {
            const Token token = scanner->tokens->tokens[i];
            printToken(token);
        }
    }

    extern TokenList *g_tokens;
    g_tokens = scanner->tokens;
    pos = 0;
    root = NULL;

    if (yyparse() == 0 && root && (debug || debugParser)) {
        if (debug) puts("===== AST =====");
        ast_print(root, 0);
        puts("\n");
    }

    ast_free(root);
    free(scanner->tokens->tokens);
    free(scanner->tokens);
    free(scanner);
}

void runFile(const char * filename) {
    FILE *f = fopen(filename, "r");
    if (!f) {
        puts("Could not open file");
        exit(64);
    }

    fseek(f, 0, SEEK_END);
    long size = ftell(f);
    rewind(f);

    char * buffer = (char *) malloc(size + 1);
    fread(buffer, 1, size, f);

    fclose(f);

    // printf("Size of file: %ld\n", size);
    // printf("Buffer: %s\n", buffer);

    run(buffer, size);

    if (hadError) {
        exit(65);
    }

    free(buffer);
}

void runPrompt() {
    for (;;) {
        printf(">  ");
        char line[102400];
        fgets(line, sizeof(line), stdin);

        if (feof(stdin)) {
            break;
        }

        run(line, (int)strlen(line));
        hadError = false;
    }
}

int main(int argc, char *argv[]) {
    char * file = NULL;
    bool repl = false;

    for (int i = 0; i < argc; i++) {
        if (strcmp(argv[i], "--help") == 0 || strcmp(argv[i], "-h") == 0) {
            break;
        } else if (strcmp(argv[i], "--debug") == 0 || strcmp(argv[i], "-d") == 0) {
            debug = true;
        } else if (strcmp(argv[i], "--debugscanner") == 0 || strcmp(argv[i], "-ds") == 0) {
            debugScanner = true;
        } else if (strcmp(argv[i], "--debugparser") == 0 || strcmp(argv[i], "-dp") == 0) {
            debugParser = true;
        } else if ((strcmp(argv[i], "--file") == 0 || strcmp(argv[i], "-f") == 0) && i+1 < argc) {
            file = argv[i+1];
        } else if (strcmp(argv[i], "--repl") == 0 || strcmp(argv[i], "-r") == 0) {
            repl = true;
        } 
    }

    if (file != NULL) {
        runFile(file);
    } else if (repl) {
        runPrompt();
    } else {
        puts("Usage: rain -f [script] <options>");
        puts("");
        puts("\t--debug         | -d           Prints debug statements");
        puts("\t--debugscanner  | -ds          Prints debug statements for scanner");
        puts("\t--debugparser   | -dp          Prints debug statements for parser");
        puts("\t--file <file>   | -f <file>    Takes file as input");
        puts("\t--repl          | -r           Repl mode");
        puts("\t--help          | -h           Prints this help message");
        exit(64);
    }

    return 0;
}
