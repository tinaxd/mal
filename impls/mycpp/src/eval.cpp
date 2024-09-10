#include "eval.h"

using namespace factory;

FPointer<FactoryValue> ast_eval(FPointer<FactoryValue> ast, Env &env)
{
    if (ast->isSymbol())
    {
        return env.get(static_cast<FPointer<FactorySymbol>>(ast)->getValue());
    }
}

FPointer<FactoryValue> factory::EVAL(FPointer<FactoryValue> ast, Env &env)
{
    return ast;
}