#!/usr/bin/env node

process.stdin.resume();  
process.stdin.setEncoding('utf8');  
process.stdin.on('data', function(data) {
  var parsed_json = JSON.parse(data);

  parsed_json.response.body = "body was replaced by middleware\n";

  process.stdout.write(JSON.stringify(parsed_json));
});
