#!/usr/bin/env node

process.stdin.resume();  
process.stdin.setEncoding('utf8');  
process.stdin.on('data', function(data) {
  var parsed_json = JSON.parse(data);
  // changing response
  parsed_json.response.status = 201;
  parsed_json.response.body = "body was replaced by JavaScript middleware\n";

  // stringifying JSON response
  var newJsonString = JSON.stringify(parsed_json);

  process.stdout.write(newJsonString);
});
