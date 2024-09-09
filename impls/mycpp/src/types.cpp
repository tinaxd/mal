#include "types.h"

using namespace factory;

FactoryList::FactoryList() {}

void FactoryList::append(FPointer<FactoryValue> value)
{
    this->values.push_back(value);
}

std::string FactoryList::printString()
{
    std::string result = "(";
    for (size_t i = 0; i < this->values.size(); i++)
    {
        if (i > 0)
        {
            result += " ";
        }
        result += this->values[i]->printString();
    }
    result += ")";
    return result;
}

FactoryInt::FactoryInt(int64_t value) : value(value) {}

std::string FactoryInt::printString()
{
    return std::to_string(this->value);
}

FactorySymbol::FactorySymbol(std::string value) : value(std::move(value)) {}

std::string FactorySymbol::printString()
{
    return this->value;
}

FPointer<FactoryList> factory::make_factory_list()
{
    return new FactoryList();
}

FPointer<FactoryInt> factory::make_factory_int(int64_t value)
{
    return new FactoryInt(value);
}

FPointer<FactorySymbol> factory::make_factory_symbol(std::string value)
{
    return new FactorySymbol(std::move(value));
}