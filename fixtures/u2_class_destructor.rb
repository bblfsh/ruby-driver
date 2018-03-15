class Testcls1
    def self.testfnc1
    end
    def initialize
        ObjectSpace.define_finalizer(self, self.class.testfnc1)
    end
end
