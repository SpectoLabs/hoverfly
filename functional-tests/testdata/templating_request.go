package testdata

var TemplatingRequest = `{
	"data": {
		"pairs": [
			{
				"request": {
					"path": [
						{
							"matcher": "exact",
							"value": "/Request"
						}
					]
				},
				"response": {
					"status": 200,
					"body": "{{ Request }}",
					"encodedBody": false,
					"templated": true
				}
			},
			{
				"request": {
					"path": [
						{
							"matcher": "exact",
							"value": "/Request.Body_jsonpath"
						}
					]
				},
				"response": {
					"status": 200,
					"body": "{{ Request.Body \"jsonpath\" \"$.test\" }}",
					"encodedBody": false,
					"templated": true
				}
			},
			{
				"request": {
					"path": [
						{
							"matcher": "exact",
							"value": "/Request.Body_jsonpath_with_each"
						}
					]
				},
				"response": {
					"status": 200,
					"body": "{{#each (Request.Body 'jsonpath' '$.currencies')}} {{@index}} : {{this}} \n {{/each}}",
					"encodedBody": false,
					"templated": true
				}
			},
			{
				"request": {
					"path": [
						{
							"matcher": "glob",
							"value": "/Request.Path_with_each/*"
						}
					]
				},
				"response": {
					"status": 200,
					"body": "{{#each Vars.splitRequestPath}}{{@index}}:{{this}} {{/each}}",
					"encodedBody": false,
					"templated": true
				}
			},
			{
				"request": {
					"path": [
						{
							"matcher": "exact",
							"value": "/Request.Body_xpath"
						}
					]
				},
				"response": {
					"status": 200,
					"body": "{{ Request.Body \"xpath\" \"/root/text\" }}",
					"encodedBody": false,
					"templated": true
				}
			},
			{
				"request": {
					"path": [
						{
							"matcher": "exact",
							"value": "/Request.Method"
						}
					]
				},
				"response": {
					"status": 200,
					"body": "{{ Request.Method }}",
					"encodedBody": false,
					"templated": true
				}
			},
			{
				"request": {
					"path": [
						{
							"matcher": "exact",
							"value": "/Request.Scheme"
						}
					]
				},
				"response": {
					"status": 200,
					"body": "{{ Request.Scheme }}",
					"encodedBody": false,
					"templated": true
				}
			},
			{
				"request": {
					"path": [
						{
							"matcher": "exact",
							"value": "/one/two/three/Request.Path"
						}
					]
				},
				"response": {
					"status": 200,
					"body": "{{ Request.Path }}",
					"encodedBody": false,
					"templated": true
				}
			},
			{
				"request": {
					"path": [
						{
							"matcher": "exact",
							"value": "/one/two/three/Request.Path0"
						}
					]
				},
				"response": {
					"status": 200,
					"body": "{{ Request.Path.[0] }}",
					"encodedBody": false,
					"templated": true
				}
			},
			{
				"request": {
					"path": [
						{
							"matcher": "exact",
							"value": "/Request.QueryParam"
						}
					]
				},
				"response": {
					"status": 200,
					"body": "{{ Request.QueryParam }}",
					"encodedBody": false,
					"templated": true
				}
			},
			{
				"request": {
					"path": [
						{
							"matcher": "exact",
							"value": "/Request.QueryParam.query"
						}
					]
				},
				"response": {
					"status": 200,
					"body": "{{ Request.QueryParam.query }}",
					"encodedBody": false,
					"templated": true
				}
			}, 
			{
				"request": {
					"path": [
						{
							"matcher": "exact",
							"value": "/Request.QueryParam_with_each"
						}
					]
				},
				"response": {
					"status": 200,
					"body": "{{#each Request.QueryParam}}{{@index}}:{{@key}}:{{this}} {{/each}}",
					"encodedBody": false,
					"templated": true
				}
			},
			{
				"request": {
					"path": [
						{
							"matcher": "exact",
							"value": "/Request.Host"
						}
					]
				},
				"response": {
					"status": 200,
					"body": "{{ Request.Host }}",
					"encodedBody": false,
					"templated": true
				}
			},
			{
				"request": {
					"path": [
						{
							"matcher": "exact",
							"value": "/SetStatusCode"
						}
					]
				},
				"response": {
					"status": 200,
					"body": "{{ setStatusCode 202 }}",
					"encodedBody": false,
					"templated": true
				}
			},
			{
				"request": {
					"path": [
						{
							"matcher": "exact",
							"value": "/global"
						}
					]
				},
				"response": {
					"status": 200,
					"body": "{\"{{ Literals.fieldName }}\":\"{{ Vars.getCityFromJsonBody }}\"}",
					"encodedBody": false,
					"templated": true
				}
			}
		],
		"globalActions": {
			"delays": []
		}, 
		"literals": [
			{
				"name": "fieldName",
				"value": "destination"
		 	}
		],
		"variables": [
			{
				"name": "getCityFromJsonBody",
				"function": "requestBody",
				"arguments": ["jsonpath", "$.city"]
		 	}, 
			{
				"name": "splitRequestPath",
				"function": "split",
				"arguments": ["Request.Path.[1]", ","]
		 	}
		]
	},
	"meta": {
		"schemaVersion": "v5.2"
	}
}`
