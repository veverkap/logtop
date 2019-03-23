at_exit do

  # do cleanup
  puts "exit"
  # now exit for real
  exit
end
while true
  begin
    sleep 1
  rescue SystemExit
    puts "system"
  rescue Interrupt
    puts "inter"
    exit
  end
end
