#include <iostream>
#include <string>
#include "read.h"
#include "printer.h"

using namespace factory;

FPointer<FactoryValue> EVAL(FPointer<FactoryValue> read)
{
    return read;
}

std::string rep(const std::string &line)
{
    const auto read = READ(line);
    const auto eval = EVAL(read);
    const auto print = PRINT(eval);
    return print;
}

int main()
{

    while (true)
    {
        std::cout << "user> ";
        std::flush(std::cout);
        std::string line;
        if (!std::getline(std::cin, line))
        {
            break;
        }
        std::cout << rep(line) << std::endl;
    }
    return 0;
}