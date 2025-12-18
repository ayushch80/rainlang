#!/usr/bin/env bash
set -e

BIN=./build-linux/rainlang
TESTDIR=tests/scanner_and_parser

FAILED=0
TOTAL=0

for file in $TESTDIR/*.rain; do
    TOTAL=$((TOTAL+2))
    base=$(basename "$file" .rain)
    tokensExpected="$TESTDIR/scanner/$base.tokens"
    astExpected="$TESTDIR/parser/$base.ast"
    astActual="/tmp/$base.actual.ast"
    tokensActual="/tmp/$base.actual.tokens"

    echo "Running $base..."

    $BIN -dp -f "$file" > "$astActual"
    $BIN -ds -f "$file" > "$tokensActual"

    # $BIN -dp -f "$file" > "$astExpected"
    # $BIN -ds -f "$file" > "$tokensExpected"

    if ! diff -u "$tokensExpected" "$tokensActual"; then
        echo "[-] Scanner test failed: $base"
        FAILED=$((FAILED+1))
    fi
    if ! diff -u "$astExpected" "$astActual"; then
        echo "[-] Parser test failed: $base"
        FAILED=$((FAILED+1))
    fi
done

if [ $FAILED -ne 0 ]; then
    echo "$FAILED/$TOTAL tests failed"
    exit 1
else
    echo "All tests passed"
fi