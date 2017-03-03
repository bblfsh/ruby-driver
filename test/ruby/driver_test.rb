require 'test_helper'

class Ruby::DriverTest < Minitest::Test

    def setup
        @parser = Yajl::Parser.new(:symbolize_keys => true)
    end

    def test_that_it_has_a_version_number
        refute_nil ::Ruby::Driver::VERSION
    end

    def test_wrong_request
        @parser.on_parse_complete = method(:callback_wrong_request)
        input = File.read("test/wrong_input.json")
        @parser.parse(input)
    end

    def callback_wrong_request(obj)
        Ruby::Driver::RequestMessage.new(obj)
        assert_raises
    end

    def test_parse_request
        requests = ['input.json', 'input_test1.json', 'input_test2.json', 'input_test3.json']
        requests.each do |req|
            @parser.on_parse_complete = method(:callback_parse_request)
            input = File.read("test/#{req}")
            @parser.parse(input)
        end
    end

    def test_multi_request
        @parser.on_parse_complete = method(:callback_parse_request)
        input = File.read("test/input_test_all.json")
        @parser.parse(input)
    end

    def callback_parse_request(obj)
        req = Ruby::Driver::RequestMessage.new(obj)
        assert_equal(Ruby::Driver::RequestMessage::PARSE_AST, req.action)
    end

    def test_start
        output = StringIO.new('', 'a')
        infile = File.read("test/input_test_all.json")
        input = StringIO.new(infile, 'r')

        @driver_version = 'driver-minitest'
        Ruby::Driver::start(@driver_version, input, output)

        @parser.on_parse_complete = method(:callback_start)
        responses = StringIO.new()
        r = output.string
        responses = StringIO.new(r, 'r')
        @parser.parse(responses)
    end

    def callback_start(json_res)
        assert_equal(Ruby::Driver::ResponseMessage::STATUS_OK, json_res[:status])
        assert_equal(@driver_version, json_res[:driver])
        assert_equal('ruby', json_res[:language])
        assert_equal(RUBY_VERSION, json_res[:language_version])
    end

end
