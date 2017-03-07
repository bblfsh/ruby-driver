require 'ruby_driver/message'
require 'ripper'
require 'yajl'
require 'json'

module RubyDriver
  class Driver
    # json parser callback
    def response_ast(obj)
      res = Message::Response.new
      begin
        req = Message::Request.new(obj)
        res.ast = Ripper.sexp_raw(req.content)
        res.status = :ok
      rescue Message::BadRequest => e
        res.status = 'error'
        res.errors = [e.message]
      end

      @output.puts(JSON.generate(res.to_hash))
    end

    def start(input, output)
      @input = input
      @output = output

      begin
        parser = Yajl::Parser.new(:symbolize_keys => true)
        parser.on_parse_complete = method(:response_ast)
        parser.parse(@input)
      rescue Exception => e
        res = Message::Response.new
        res.status = 'fatal'
        res.errors = [e.message]
        @output.puts(JSON.generate(res.to_hash))
      end

    end

  end
end
