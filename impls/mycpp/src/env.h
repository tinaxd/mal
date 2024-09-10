#pragma once

#include <vector>
#include <functional>
#include "types.h"
#include <unordered_map>
#include <string>

namespace factory
{
    class Env
    {
        std::unordered_map<std::string, FPointer<FactoryValue>> env;

    public:
        Env() = default;

        void set(const std::string &key, FPointer<FactoryValue> value);
        FPointer<FactoryValue> get(const std::string &key);
    };
}