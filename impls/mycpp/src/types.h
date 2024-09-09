#pragma once
#include <cstdint>
#include <string>
#include <vector>

namespace factory
{
    template <typename T>
    using FPointer = T *;

    class FactoryValue
    {
    public:
        virtual std::string printString() = 0;
    };

    class FactoryList : public FactoryValue
    {
    private:
        std::vector<FPointer<FactoryValue>> values;

    public:
        FactoryList();
        void append(FPointer<FactoryValue> value);
        std::string printString() override;
    };

    class FactoryInt : public FactoryValue
    {
    private:
        int64_t value;

    public:
        FactoryInt(int64_t value);
        std::string printString() override;
    };

    class FactorySymbol : public FactoryValue
    {
    private:
        std::string value;

    public:
        FactorySymbol(std::string value);
        std::string printString() override;
    };

    FPointer<FactoryList> make_factory_list();
    FPointer<FactoryInt> make_factory_int(int64_t value);
    FPointer<FactorySymbol> make_factory_symbol(std::string value);
    template <typename T>
    FPointer<T> f_null = nullptr;
    // FPointer<FactoryValue
}