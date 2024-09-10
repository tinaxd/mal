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

        virtual inline bool isList() const { return false; }
        virtual inline bool isInt() const { return false; }
        virtual inline bool isSymbol() const { return false; }
    };

    class FactoryList : public FactoryValue
    {
    private:
        std::vector<FPointer<FactoryValue>> values;

    public:
        FactoryList();
        void append(FPointer<FactoryValue> value);
        std::string printString() override;

        inline bool isList() const override { return true; }

        size_t size() const;
        FPointer<FactoryValue> &operator[](size_t index);
        const FPointer<FactoryValue> &operator[](size_t index) const;
    };

    class FactoryInt : public FactoryValue
    {
    private:
        int64_t value;

    public:
        FactoryInt(int64_t value);
        std::string printString() override;

        inline bool isInt() const override { return true; }
    };

    class FactorySymbol : public FactoryValue
    {
    private:
        std::string value;

    public:
        FactorySymbol(std::string value);
        std::string printString() override;

        inline bool isSymbol() const override { return true; }
        inline std::string getValue() { return this->value; }
    };

    FPointer<FactoryList> make_factory_list();
    FPointer<FactoryInt> make_factory_int(int64_t value);
    FPointer<FactorySymbol> make_factory_symbol(std::string value);
    template <typename T>
    FPointer<T> f_null = nullptr;
    // FPointer<FactoryValue

    template <typename T>
    FPointer<T> f_ptr_cast(FPointer<FactoryValue> value)
    {
        return static_cast<FPointer<T>>(value);
    }
}