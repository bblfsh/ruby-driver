module FooMod

class Foo
end

class Bar < Foo
    @@classvar = 10
    def initialize(a)
        @instance_var = a
    end

    def foo()
    end

    def use()
        puts @@classvar
        puts @instance_var
        foo()
        self.foo()
    end

    def selfinstance()
        b = Bar.new(1)
    end
end
end
