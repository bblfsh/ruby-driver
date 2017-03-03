require "ruby/driver/version"
require 'active_support'
require 'ripper'
require 'yajl'

module Ruby
  module Driver
      class MyException < StandardError
      end

      class RequestMessage
          attr_reader :action
          attr_reader :language
          attr_reader :language_version
          attr_reader :content
          BAD_REQUEST = "Bad Request"
          PARSE_AST = "ParseAST"

          def initialize(json_req)
                  check_req(json_req)
                  @action = json_req[:action]
                  @language = json_req[:language]
                  @language_version = json_req[:language_version]
                  @content = json_req[:content]
          end

          def do_action()
              case action
              when PARSE_AST
                  Ripper.sexp_raw(content).to_s()
              else
                  raise MyException, BAD_REQUEST
              end
          end

          private
          def check_req(req)
              raise MyException, BAD_REQUEST + ": Missing action" unless req.has_key?(:action)
              raise MyException, BAD_REQUEST + ": Missing content"unless req.has_key?(:content)
              req.each do |key, value|
                  raise MyException, BAD_REQUEST + ": " + key + ": " + "the value for this key is not a String" unless value.class == String
              end
          end
      end

      class ResponseMessage
          attr_accessor :status
          attr_accessor :errors
          attr_accessor :driver
          attr_accessor :language
          attr_accessor :language_version
          attr_accessor :ast

          def initialize(driver)
              @driver = driver
              @language = "ruby"
              @language_version = RUBY_VERSION
          end

          STATUS_OK = "ok"
          STATUS_ERROR = "error"
          STATUS_FATAL = "fatal"
      end

      # json parser callback
      def self.response_ast(obj)
          res = ResponseMessage.new($DRIVER)
          begin
              req = RequestMessage.new(obj)
              if req.action == RequestMessage::PARSE_AST
                  res.ast = Ripper.sexp_raw(req.content).to_s()
              else
                  raise MyException, RequestMessage::BAD_REQUEST + ": Unknown action"
              end

              res.status = ResponseMessage::STATUS_OK
          rescue MyException => e
              res.status = ResponseMessage::STATUS_ERROR
              res.errors = [e.message]
          rescue Exception => e
              res.status = ResponseMessage::STATUS_FATAL
              res.errors = [e.message]
              fatal_msg = e.message
          ensure
              @output.write(ActiveSupport::JSON.encode(res))
              if fatal_msg != nil then abort(fatal_msg) end
          end
      end

      def self.start(driver, input, output)
          $DRIVER = driver
          if $DRIVER == nil
              $DRIVER = "driver-test"
          end

          @input = input
          @output = output

          parser = Yajl::Parser.new(:symbolize_keys => true)
          parser.on_parse_complete = method(:response_ast)
          parser.parse(@input)
      end
  end
end
