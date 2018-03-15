class Testcls1
    attr_accessor :a
end
class Testcls2
    attr_reader :b
    attr_writer :b
end
class Testcls3
    @c = 1
    attr_reader :c
    def c=(new_c)
        @c = new_c
    end
end
