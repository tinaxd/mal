#include <iostream>
#include <string>

std::string READ(const std::string &line)
{
    return line;
}

std::string EVAL(const std::string &read)
{
    return read;
}

std::string PRINT(const std::string &eval)
{
    return eval;
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