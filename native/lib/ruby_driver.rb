require 'ruby_driver/message'
require 'ruby_driver/node_converter'

require 'json'

require 'parser/current'

module RubyDriver
  # Driver implements the functionality to parse a JSON object per line, which
  # represents a Request with ruby source code, and reply with a JSON object
  # response containing the AST of the code.
  class Driver

    # response_ast extracts the AST from a request, and returns a response.
    def response_ast(hash_req)
      res = Message::Response.new
      begin
        res = Message::Request.new(hash_req)
        node = Parser::CurrentRuby.parse(res.content)
        res.ast = NodeConverter::Converter.new(node).visit.tohash()
        res.status = :ok
      rescue Message::BadRequest => e
        res.status = :error
        res.errors = [e.message]
      end

      return res
    end

    # start unmarshal requests from input, and marshal responses to output.
    def start(input, output)
      @input = input
      @output = output

      @output.sync = true
      @input.each_line do |line|
        begin
          res = response_ast(JSON.parse(line))
          @output.puts(JSON.generate(res.to_hash))
        rescue Exception => e
          res = Message::Response.new
          res.status = :fatal
          res.errors = [e.message]
          @output.puts(JSON.generate(res.to_hash))
        end
      end
    end
  end
end
