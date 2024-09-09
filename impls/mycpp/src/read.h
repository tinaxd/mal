#pragma
#include <string>

namespace factory
{
    std::string READ(const std::string &line);

    class Reader
    {
    public:
        char peek();
        char next();
    };
}