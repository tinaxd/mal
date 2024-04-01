#include "reader.h"
#include <stdio.h>
#include <pcre.h>

Token *tokenize(const char *input)
{
    const re_str = "[\s,]*(~@|[\[\]{}()'`~^@]|\"(?:\\.|[^\\\"])*\"?|;.*|[^\s\[\]{}('\"`,;)]*)";

    pcre *re;
    const char *err_str;
    int err_offset;
    re = pcre_compile(re_str, 0, &err_str, &err_offset, NULL);
    if (re == NULL)
    {
        fprintf(stderr, "pcre error: %s\n", err_str);
        return NULL;
    }

    int ovector = 0;
    int matched = pcre_exec(re, NULL, input, strlen(input), 0, 0, &ovector, 1);
    pcre_free(re);

    if (matched < 0)
    {
        return NULL;
    }

    Token *tokens = malloc(sizeof(Token) * matched);
}

Reader *read_str(const char *input)
{
    Token *tokens = tokenize(input);
    Reader *r = malloc(sizeof(Reader));
    r->position = 0;
    r->tokens = tokens;
    return r;
}
