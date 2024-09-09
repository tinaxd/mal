#include "read.h"
#include <regex>

using namespace factory;

std::string READ(const std::string &line)
{
    const auto pattern = std::regex(R"***([\s,]*(~@|[\[\]{}()'`~^@]|"(?:\\.|[^\\"])*"?|;.*|[^\s\[\]{}('"`,;)]*))***");
}