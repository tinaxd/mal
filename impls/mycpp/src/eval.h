#pragma once
#include "types.h"
#include "env.h"

namespace factory
{
    FPointer<FactoryValue> EVAL(FPointer<FactoryValue> ast, Env &env);
}