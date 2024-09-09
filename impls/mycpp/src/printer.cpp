#include "printer.h"

using namespace factory;

std::string factory::pr_str(FPointer<FactoryValue> value)
{
    return value->printString();
}

std::string factory::PRINT(FPointer<FactoryValue> eval)
{
    return pr_str(eval);
}