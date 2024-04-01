#pragma once

typedef struct
{
    char *value;
} Token;

typedef struct
{
    Token *tokens;
    int position;
} Reader;

Reader *read_str(const char *input);
