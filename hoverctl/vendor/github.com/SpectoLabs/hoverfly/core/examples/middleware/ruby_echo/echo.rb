#!/usr/bin/env ruby
# encoding: utf-8
while payload = STDIN.gets
  next unless payload

  STDOUT.puts payload

  STDERR.puts "Payload data: #{payload}"

end