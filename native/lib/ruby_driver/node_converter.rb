require 'parser/current'

# Nodes doc:
# https://github.com/whitequark/parser/blob/master/doc/AST_FORMAT.md

module NodeConverter
  class Converter
    @@typekey = "@type"
    @@tokenkey = "@token"
    @@statement_send = ["continue", "lambda", "require", "each", "public", "protected",
                        "private"]
    @@operators = ["+", "-", "*", "/", "%", "**", "&", "|", "^", "~", "<<", ">>",
                   "==", "===", "<=", ">=", "!=", "!", "eql?", "equal?", "<==>"]

    def initialize(node, comments)
      @empty_with_comments = false
      if node.is_a?(NilClass) and comments != nil
        # Since Ruby parses comments and "normal" nodes separately, it will consider
        # a file with only comments as a NilNode and refuse to parse anymore. We fix it
        # changing the NilNode to a empty module node
        @empty_with_comments = true
      elsif not node.is_a?(Parser::AST::Node)
        raise "Object is not a Parser::AST::Node, is: #{node.class.name}"
      end

      @root = node
      @comments = comments
      @dict = {}
      @qualified_stack = []
    end

    def tohash()
      @dict["ast"] = {"file" => convert(@root), @@typekey => "module"}
      add_comments()
      return @dict["ast"]
    end

    private

    def convert(node)
      type = node_type(node)

      case type

      when "int", "float", "str"
        return sexp_to_hash(node, {"l_" + @@tokenkey => 0})

      when "lvar", "ivar", "cvar", "gvar", "arg", "kwarg", "restarg", "blockarg"
        return sexp_to_hash(node, {@@tokenkey + "." + @@tokenkey => 0}, 1, "children")

      when "pair", "irange", "erange", "alias"
        return sexp_to_hash(node, {"_1" => 0, "_2" => 1})

      when "iflipflop", "eflipflop"
        return process_flipflop(node)

      when "lvasgn", "ivasgn", "cvasgn", "or_asgn", "and_asgn"
        return sexp_to_hash(node, {"target" => 0, "value" => 1})

      when "block"
        return sexp_to_hash(node, {"blockdata" => 0, "args.children" => 1, "body" => 2})

      when "array", "hash"
        return sexp_to_hash(node, {}, 0, "contents")

      when "optarg", "kwoptarg"
        return sexp_to_hash(node, {@@tokenkey => 0, "default" => 1})

      when "splat", "kwsplat", "defined?", "kwrestarg"
        return sexp_to_hash(node, {@@tokenkey => 0})

      when "casgn"
        return sexp_to_hash(node, {"base" => 0, "selector" => 1, "value" => 2})

      when "csend", "send"
        return process_send(node)

      when "complex", "rational", "sym"
        return sexp_to_hash(node, {@@tokenkey + "." + @@tokenkey => 0})

      # the inner nodes of the above
      when "Complex", "Rational", "Symbol"
        return {@@typekey => node_type(node), @@tokenkey => node.to_s}

      when "masgn"
        return sexp_to_hash(node, {"targets" => 0, "values" => 1})

      when "op_asgn"
        return sexp_to_hash(node, {"target" => 0, "operator" => 1, "value" => 2})

      when "module"
        d = sexp_to_hash(node, {}, 1, "begin")
        d[@@tokenkey] = node.children[0].children[1].to_s
        return d

      when "class"
        d = sexp_to_hash(node, {"parent" => 1}, 2, "body")
        d[@@tokenkey] = node.children[0].children[1].to_s
        return d

      when "sclass"
        return sexp_to_hash(node, {"object" => 0}, 1, "body")

      when "def"
        return sexp_to_hash(node, {"s_" + @@tokenkey => 0, "args" => 1}, 2, "body")

      when "undef", "yield", "break", "next", "return"
        return sexp_to_hash(node, {"target" => 0})

      when "and", "or"
        return sexp_to_hash(node, {"left" => 0, "right" => 1})

      when "case"
        d = sexp_to_hash(node, {"casevar" => 0})
        if node.children.length > 2
          d["when_clauses"] = node.children[1..-2].map{ |x| convert(x) }.compact
        end
        d["else"] = convert(node.children[-1])
        return d

      when "when"
        d = {@@typekey => "when"}
        d["conditions"] = node.children[0..-2].map{ |x| convert(x) }.compact
        d["body"] = convert(node.children[-1])
        return d

      when "const"
        return sexp_to_hash(node, {"base" => 0, @@tokenkey => 1})

      when "while", "until", "while_post", "until_post"
        return sexp_to_hash(node, {"condition" => 0, "body" => 1})

      when "begin", "kwbegin", "preexe", "postexe"
        return sexp_to_hash(node, {}, 0, "body")

      when "for"
        return sexp_to_hash(node, {"iterators" => 0, "iterated" => 1, "body" => 2})

      when "resbody"
        return sexp_to_hash(node, {"exceptions" => 0, "alias" => 1, "body" => 2})

      when "rescue"
        d = sexp_to_hash(node, {"body" => 0})
        if node.children.length > 2
          d["handlers"] = node.children[1..-2].map{ |x| convert(x) }.compact
        end
        d["else"] = convert(node.children[-1])
        return d

      when "ensure"
        return sexp_to_hash(node, {"body" => 0, "ensure_body" => 1})

      when "if"
        return sexp_to_hash(node, {"condition" => 0, "body" => 1, "else" => 2})

      when "defs" # "singleton method"
        return sexp_to_hash(node, {"base" => 0, @@tokenkey => 1, "args.children" => 2, "class" => 3})

      when "regexp"
        return sexp_to_hash(node, {"text" => 0, "options" => 1})

      when "regopt"
        return sexp_to_hash(node, {}, 0, "options")

      when "NilClass"
        if @empty_with_comments
          return {@@typekey=> "module", @@tokenkey => "empty_module"}
        else
          return {@@typekey=> "NilNode"}
        end

      else
        # default conversion
        return sexp_to_hash(node, {}, 0, "children")
      end
    end

    def node_type(node)
      (node.is_a?(Parser::AST::Node) ? node.type : node.class).to_s
    end

    # Convert positional children nodes to hashtable nodes keyed to a named
    # attribute using a attrname => position hashmap in the "table" argument.
    #
    # The attrnames can use a micro-DSL to specify further operations on the
    # children nodes:
    #
    # - Start with "s_": the children node wont be visited and will be converted
    #   to string as is.
    # - Starts with "l_": take child literal value, without calling convert.
    # - Ends with ".foo": the child node will be visited and the key "foo" from
    #   the resulting dictionary will be assigned.
    #
    # In any other case, the child will just be converted and assigned to the specified
    # key in the node dictionary.
    #
    # The cdr_index and cdr_key arguments, if present, specify that any other children
    # after cdr_index will be converted and assigned as a list of dictnodes to the cdr_key
    # property in the node.
    def sexp_to_hash(node, table, cdr_index=nil, cdr_key=nil)
      d = {@@typekey => node_type(node)}

      table.each do |propname, idx|
        if propname.start_with? "s_"
          d[propname[ 2..-1]] = node.children[idx].to_s

        elsif propname.start_with? "l_"
          d[propname[2..-1]] = node.children[idx]

        elsif propname.include? "."
          propname, childkey = propname.split(".")
          d[propname] = convert(node.children[idx])[childkey]

        elsif node.children[idx].is_a? Parser::AST::Node
          d[propname] = convert(node.children[idx])

        else
          d[propname] = node.children[idx]
        end
      end

      if cdr_index != nil and cdr_key != nil and node.respond_to?("children") and \
        node.children.length > cdr_index

        d[cdr_key] = node.children[cdr_index..-1].map{ |x| convert(x) }.compact
      end

      return add_position(node, d)
    end

    def process_send(node)
      # Same as sexp_to_hash but without visiting the base yet since we need
      # to visit it after visiting this node to add the selectors of qualified
      # identifiers to qnames in the same order
      hash_send = {}
      hash_send["base"] = node.children[0]
      hash_send["selector"] = node.children[1]

      if node.children.length > 2
        hash_send["values"] = node.children[2..-1].map{ |x| convert(x) }.compact
      end

      selector = hash_send["selector"]
      hash_send = add_position(node, hash_send)

      def convert_base(hash)
        if hash["base"].is_a? Parser::AST::Node
          hash["base"] = convert(hash["base"])
        end
        return hash
      end

      def push_identifier(id_node, selector=nil)
        qual_node = {
          @@typekey => "uast:Identifier",
          "@pos" => id_node["@pos"]
        }

        if selector.nil?
          qual_node["Name"] = id_node["selector"]
        else
          qual_node["Name"] = selector
        end

        @qualified_stack.push(qual_node)
      end

      if @@statement_send.include? selector
        if selector == "require"
            hash_send[@@typekey] = "send_require"
        else
            hash_send[@@typekey] = "send_statement"
        end

      elsif @@operators.include? selector
        hash_send[@@typekey] = "send_operator"

      elsif selector[-1] == "=" and not hash_send["values"].nil?
        hash_send[@@typekey] = "send_assign"
        hash_send[@@tokenkey] = selector[0..-2]

        if not hash_send["base"].nil?
          push_identifier(hash_send, selector[0..-2])
        end

      elsif hash_send["base"].nil?
        hash_send[@@typekey] = "send_qualified"
        process_qualified(hash_send)

      else
        if selector == "[]"
          hash_send[@@typekey] = "send_array"
        elsif hash_send.key?("values")
          hash_send[@@typekey] = "send_call"
        else
          hash_send[@@typekey] = "send_identifier"
        end

        push_identifier(hash_send)
      end

      return convert_base(hash_send)
    end

    def process_qualified(node)
      node["qnames"] = Array.new

      @qualified_stack.each do |x|
        node["qnames"].push(x)
      end

      node["qnames"].push(node["selector"])
      @qualified_stack = []
    end

    def process_flipflop(node)
      hash_ff = sexp_to_hash(node, {"flip_1" => 0, "flip_2" => 1})

      # Change the type of the flip-parts from "send" to "flip_1/2" and get the token
      ["flip_1", "flip_2"].each do |key|
        hash_ff[key][@@typekey] = key
        hash_ff[key][@@tokenkey] = hash_ff[key]["selector"]
        hash_ff[key].delete("selector")
      end

      return hash_ff
    end

    def add_from_subelem(node, hash, key)
      subelem = node.loc.send(key)
      if subelem != nil
        hash["@pos"] = {
          @@typekey => "uast:Positions",
          "start" => {
            @@typekey => "uast:Position",
            "line" => subelem.begin.line,
            "col" => subelem.begin.column + 1
          },
          "end" => {
            @@typekey => "uast:Position",
            "line" => subelem.end.line,
            # str inside str have cols set at 0 from the native AST
            "col" => subelem.end.column > 0 ? subelem.end.column : 1
          },
        }
      end
    end

    def add_position(node, hash)
      case hash[@@typekey]

      when "defined?", "module", "class", "sclass", "def", "defs",
        "undef", "alias", "super", "zsuper", "yield", "if", "when",
        "case", "while", "while-post", "for", "break", "next", "redo",
        "return", "rescue", "ensure", "retry", "preexe", "postexe"
        subelem = "keyword"

      when "optarg", "restarg", "blockarg", "kwarg", "kwoptarg",
           "kwrestarg"
        subelem = "name"

      when "not"
        subelem = "operator"

      else
        subelem = "expression"
      end

      if hash[@@typekey] == "if" and not node.loc.respond_to?("keyword")
        subelem = "question"
      end

      if node != nil
        add_from_subelem(node, hash, subelem)
      end

      return hash
    end

    # Add comments inside the root node "comments" field
    def add_comments()
      if @comments == nil
        return
      end

      comments = []

      @comments.each do |comment|
        commentdict = {
          @@typekey => "comment",
          @@tokenkey => comment.text[1..-1],
          "inline" => comment.inline?,
          "documentation" => comment.document?,
          "@pos" => {
            @@typekey => "uast:Positions",
            "@start" => {
              @@typekey => "uast:Position",
              "line" => comment.loc.first_line,
              "col" => comment.loc.column + 1
            },
            "@end" => {
              @@typekey => "uast:Position",
              "line" => comment.loc.last_line,
              "col" => comment.loc.last_column
            },
          }
        }
        comments.push(commentdict)
      end

      if comments.length > 0
        @dict["ast"]["file"][:comments] = comments
      end
    end


  end
end
