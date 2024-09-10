#include <iostream>
#include <string>
#include "read.h"
#include "printer.h"
#include "env.h"
#include "eval.h"

using namespace factory;

std::string rep(const std::string &line, Env &env)
{
    const auto read = READ(line);
    const auto eval = EVAL(read);
    const auto print = PRINT(eval);
    return print;
}

int main()
{

    Env env;

    while (true)
    {
        std::cout << "user> ";
        std::flush(std::cout);
        std::string line;
        if (!std::getline(std::cin, line))
        {
            break;
        }
        try
        {
            std::cout << rep(line, env) << std::endl;
        }
        catch (const UnexpectedEOFError &e)
        {
            std::cout << "Unexpected EOF" << std::endl;
        }
    }
    return 0;
}