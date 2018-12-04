def n_queens(n)
  if n == 1
    return "Q"
  elsif n < 4
    puts "no solutions for n=#{n}"
    return ""
  end
 
  evens = (2..n).step(2).to_a
  odds = (1..n).step(2).to_a
 
  rem = n % 12  # (1)
  nums = evens  # (2)
 
  nums.push(nums.shift) if rem == 3 or rem == 9  # (3)
 
  # (4)
  if rem == 8
    odds = odds.each_slice(2).inject([]) {|ary, (a,b)| ary += [b,a]}
  end
  nums.concat(odds)
 
  # (5)
  if rem == 2
    idx = []
    [1,3,5].each {|i| idx[i] = nums.index(i)}
    nums[idx[1]], nums[idx[3]] = nums[idx[3]], nums[idx[1]]
    nums.slice!(idx[5])
    nums.push(5)
  end
 
  # (6)
  if rem == 3 or rem == 9
    [1,3].each do |i|
      nums.slice!( nums.index(i) )
      nums.push(i)
    end
  end
 
  # (7)
  board = Array.new(n) {Array.new(n) {"."}}
  n.times {|i| board[i][nums[i] - 1] = "Q"}
  board.inject("") {|str, row| str << row.join(" ") << "\n"}
end
 
(1 .. 15).each {|n| puts "n=#{n}"; puts n_queens(n); puts}
