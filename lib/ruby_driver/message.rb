module Message

  class BadRequest < StandardError
    def to_s
      'Bad Request'
    end
  end

  class Request
    attr_reader :content

    def initialize(req)
      check_req(req)
      @content = req[:content]
    end

    private
    def check_req(req)
      raise BadRequest unless req.has_key?(:content)
      raise BadRequest unless req[:content].is_a? String
    end

  end

  class Response
    attr_accessor :status
    attr_accessor :errors
    attr_accessor :ast

    def to_hash
      hash = {}

      if defined? @status
        hash['status'] = @status
      end

      if defined? @errors
        hash[:errors] = @errors
      end

      if defined? @ast
        hash[:ast] = @ast
      end

      return hash
    end

  end

end
