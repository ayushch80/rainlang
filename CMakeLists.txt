cmake_minimum_required(VERSION 3.31)
project(rainlang C)

set(CMAKE_C_STANDARD 17)
set(CMAKE_C_STANDARD_REQUIRED ON)

find_package(BISON REQUIRED)

BISON_TARGET(
        RainParser
        ${CMAKE_SOURCE_DIR}/src/parser.y
        ${CMAKE_BINARY_DIR}/parser.tab.c
        DEFINES_FILE ${CMAKE_BINARY_DIR}/parser.tab.h
)

add_executable(rainlang src/main.c
        src/scanner.c
        src/scanner.h
        src/types.h
        src/types.c
        src/tokentype.h
        src/token.c
        src/token.h
        src/main.h
        src/ast.c
        src/ast.h

        ${BISON_RainParser_OUTPUTS}
)

target_include_directories(rainlang PRIVATE
        src
        ${CMAKE_BINARY_DIR}
)

#add_dependencies(rainlang RainParser)

target_compile_options(rainlang PRIVATE
        -Wall -Wextra -Wpedantic
)