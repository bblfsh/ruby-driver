class Foo
    def foo(a)
    end

    def bar()
        this.foo(this.bar(5))
    end
end

a = Foo.new()
a.foo(1)
a.bar + 1
a.foo = 3
