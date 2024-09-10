#include "eval.h"

using namespace factory;

FPointer<FactoryValue> eval_ast(FPointer<FactoryValue> ast, Env &env)
{
    if (ast->isSymbol())
    {
        return env.get(f_ptr_cast<FactorySymbol>(ast)->getValue());
    }
    else if (ast->isList())
    {
        auto lst = f_ptr_cast<FactoryList>(ast);
        const auto s = lst->size();
        auto new_list = make_factory_list();
        for (size_t i = 0; i < s; i++)
        {
            auto item = (*lst)[i];
            auto evaled = EVAL(item, env);
            new_list->append(evaled);
        }
        return new_list;
    }

    return ast;
}

FPointer<FactoryValue> factory::EVAL(FPointer<FactoryValue> ast, Env &env)
{
    if (!ast->isList())
    {
        return eval_ast(ast, env);
    }
    auto list = f_ptr_cast<FactoryList>(ast);
    if (list->size() == 0)
    {
        return ast;
    }

    auto evaled = f_ptr_cast<FactoryList>(eval_ast(list, env));
    auto f = (*evaled)[0];
    auto args = make_factory_list();
    for (size_t i = 1; i < evaled->size(); i++)
    {
        args->append((*evaled)[i]);
    }
}