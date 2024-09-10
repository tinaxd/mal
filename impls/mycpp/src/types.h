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

        virtual inline bool isList() { return false; }
        virtual inline bool isInt() { return false; }
        virtual inline bool isSymbol() { return false; }
    };

    class FactoryList : public FactoryValue
    {
    private:
        std::vector<FPointer<FactoryValue>> values;

    public:
        FactoryList();
        void append(FPointer<FactoryValue> value);
        std::string printString() override;

        inline bool isList() override { return true; }
    };

    class FactoryInt : public FactoryValue
    {
    private:
        int64_t value;

    public:
        FactoryInt(int64_t value);
        std::string printString() override;

        inline bool isInt() override { return true; }
    };

    class FactorySymbol : public FactoryValue
    {
    private:
        std::string value;

    public:
        FactorySymbol(std::string value);
        std::string printString() override;

        inline bool isSymbol() override { return true; }
        inline std::string getValue() { return this->value; }
    };

    FPointer<FactoryList> make_factory_list();
    FPointer<FactoryInt> make_factory_int(int64_t value);
    FPointer<FactorySymbol> make_factory_symbol(std::string value);
    template <typename T>
    FPointer<T> f_null = nullptr;
    // FPointer<FactoryValue
}