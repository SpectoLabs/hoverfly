#!/usr/bin/env ruby
# encoding: utf-8

require 'rubygems'
require 'json'

while payload = STDIN.gets
  next unless payload

  jsonPayload = JSON.parse(payload)

  jsonPayload["response"]["body"] = "body was replaced by middleware\n"

  STDOUT.puts jsonPayload.to_json

end