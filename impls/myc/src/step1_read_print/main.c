#include <stdio.h>
#include <string.h>

char *read(char *input)
{
    return input;
}

char *eval(char *eval)
{
    return eval;
}

char *print(char *print)
{
    return print;
}

void repl_loop(void)
{
    while (1)
    {
        printf("user> ");
        char input[1024];
        char *g = fgets(input, 1024, stdin);
        if (g == NULL)
        {
            break;
        }

        // trim newline from input
        input[strcspn(input, "\n")] = 0;

        char *step1 = read(input);
        char *step2 = eval(step1);
        char *step3 = print(step2);
        printf("%s\n", step3);
    }
}

int main(void)
{
    repl_loop();
}