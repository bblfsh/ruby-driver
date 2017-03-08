require 'test_helper'

class RubyDriverTest < Minitest::Test

  def setup
    @driver = RubyDriver::Driver.new
  end

  def test_wrong_request
    input = File.read("test/wrong_input.json")
    @driver.response_ast(JSON::parse(input))
    assert_raises
  end

  def test_parse_request
    requests = ['input.json', 'input_test1.json', 'input_test2.json', 'input_test3.json']
    requests.each do |req|
      input = File.read("test/#{req}")
      res = @driver.response_ast(JSON::parse(input))
      assert_kind_of(Message::Response, res)
      assert(res.to_hash.has_key?(:status), res.to_hash)
      assert_equal(:ok, res.status)
      assert(res.to_hash.has_key?(:ast), res.to_hash)
    end
  end

  def test_start
    output = StringIO.new('', 'a')
    infile = File.read("test/input_test_all.json")
    input = StringIO.new(infile, 'r')

    @driver.start(input, output)

    r = output.string
    responses = StringIO.new(r, 'r')

    responses.each_line do |line|
      hash_res = JSON.parse(line)
      assert(hash_res.has_key?('status'))
      assert_equal('ok', hash_res['status'])
      assert(hash_res.has_key?('ast'))
    end
  end

end
