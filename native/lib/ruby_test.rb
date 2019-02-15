require_relative './ruby_driver/node_converter'

require 'parser/current'

node, comments = Parser::CurrentRuby.parse_with_comments(File.read("../../fixtures/func_with_comments.rb"))
ast = NodeConverter::Converter.new(node, comments).tohash()
puts ast
