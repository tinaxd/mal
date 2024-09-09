#include "read.h"
#include <regex>
#include "types.h"
#include <cctype>
#include <iostream>

using namespace factory;

std::vector<std::string> tokenize(const std::string &str);
FPointer<FactoryValue> read_form(Reader &reader);
FPointer<FactoryValue> read_list(Reader &reader);
FPointer<FactoryValue> read_atom(Reader &reader);

Reader::Reader(std::vector<std::string> tokens) : tokens(tokens), position(0) {}

const std::optional<std::string> Reader::peek()
{
    if (this->position >= this->tokens.size())
    {
        return std::nullopt;
    }
    return this->tokens[this->position];
}

const std::optional<std::string> Reader::next()
{
    if (this->position >= this->tokens.size())
    {
        return std::nullopt;
    }
    return this->tokens[this->position++];
}

FPointer<FactoryValue> read_str(const std::string &str)
{
    auto reader = Reader(tokenize(str));
    return read_form(reader);
}

std::vector<std::string> tokenize(const std::string &str)
{
    const auto pattern = std::regex(R"***([\s,]*(~@|[\[\]{}()'`~^@]|"(?:\\.|[^\\"])*"?|;.*|[^\s\[\]{}('"`,;)]*))***");

    std::string::const_iterator searchStart(str.cbegin());

    std::vector<std::string> tokens;
    std::smatch matches;
    while (std::regex_search(searchStart, str.cend(), matches, pattern))
    {
        // std::cout << "match found: " << matches[1] << std::endl;
        if (matches[1] == "")
            break;
        tokens.push_back(matches[1]);
        searchStart = matches.suffix().first;
    }

    return tokens;
}

FPointer<FactoryValue> read_form(Reader &reader)
{
    const auto &firstTokenOpt = reader.peek();

    if (!firstTokenOpt.has_value())
    {
        return f_null<FactoryValue>;
    }

    const auto &firstToken = firstTokenOpt.value();

    if (firstToken == "(")
    {
        return read_list(reader);
    }
    return read_atom(reader);
}

FPointer<FactoryValue> read_list(Reader &reader)
{
    auto list = make_factory_list();
    reader.next(); // consume "("
    while (true)
    {
        if (reader.peek() == ")")
        {
            break;
        }

        const auto value = read_form(reader);
        if (value == f_null<FactoryValue>)
        {
            throw std::runtime_error("unexpected EOF while reading");
        }
        list->append(value);
    }
    reader.next(); // consume ")"
    return list;
}

FPointer<FactoryValue> read_atom(Reader &reader)
{
    const auto tokenOpt = reader.next();
    if (!tokenOpt.has_value())
    {
        return f_null<FactoryValue>;
    }
    const auto &token = tokenOpt.value();
    if (std::isdigit(token[0]) || (token[0] == '-' && std::isdigit(token[1])))
    {
        return make_factory_int(std::stoll(token));
    }

    return make_factory_symbol(token);
}

FPointer<FactoryValue> factory::READ(const std::string &line)
{
    return read_str(line);
}