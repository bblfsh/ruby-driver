begin
    raise 'foo'
rescue Exception => e
    a = e.message
else
    b = 0
end
