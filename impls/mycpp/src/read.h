#pragma
#include <string>
#include <vector>
#include "types.h"
#include <optional>
#include <stdexcept>

namespace factory
{
    FPointer<FactoryValue> READ(const std::string &line);

    class Reader
    {
    private:
        std::vector<std::string> tokens;
        size_t position;

    public:
        Reader(std::vector<std::string> tokens);
        const std::optional<std::string> peek();
        const std::optional<std::string> next();
    };

    class UnexpectedEOFError : public std::runtime_error
    {
    public:
        UnexpectedEOFError() : std::runtime_error("Unexpected EOF") {}
    };
}