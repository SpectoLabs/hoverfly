package testdata

var FormDataMatch = `{
	"data": {
		"pairs": [
			{
				"request": {
					"method": [{
						"matcher": "exact", 
						"value": "POST"
					}],
					"body": [
						{
							"matcher": "form",
							"value": {
								"grant_type": [{
									"matcher": "exact", 
									"value": "authorization_code"
								}]
							}
						}
					]
				},
				"response": {
					"status": 200,
					"body": "form data matches",
					"encodedBody": false,
					"templated": false
				}
			},
			{
				"request": {
					"method": [{
						"matcher": "exact", 
						"value": "POST"
					}],
					"body": [
						{
							"matcher": "form",
							"value": {
								"grant_type": [{
									"matcher": "exact", 
									"value": "authorization_code"
								}],
								"client_assertion": [{
									"matcher": "exact", 
									"value": "fake-client-assertion"
								}],
								"code": [{
									"matcher": "exact", 
									"value": "fake-auth-code-1"
								}]
							}
						}
					]
				},
				"response": {
					"status": 200,
					"body": "all form data matches",
					"encodedBody": false,
					"templated": false
				}
			},
			{
				"request": {
					"method": [{
						"matcher": "exact", 
						"value": "POST"
					}],
					"body": [
						{
							"matcher": "form",
							"value": {
								"code": [{
									"matcher": "exact", 
									"value": "fake-auth-code-2"
								}],
								"client_assertion": [{
									"matcher": "jwt", 
									"value": "{\"header\":{\"alg\":\"HS256\"},\"payload\":{\"sub\":\"1234567890\",\"name\":\"John Doe\"}}"
								}]
							}
						}
					]
				},
				"response": {
					"status": 200,
					"body": "jwt in form data matches",
					"encodedBody": false,
					"templated": false
				}
			},
			{
				"request": {
					"method": [{
						"matcher": "exact", 
						"value": "POST"
					}],
					"body": [
						{
							"matcher": "form",
							"value": {
								"code": [{
									"matcher": "exact", 
									"value": "fake-auth-code-3"
								}],
								"client_assertion": [{
									"matcher": "jwt", 
									"value": "{\"header\":{\"alg\":\"HS256\"},\"payload\":{}}",
									"doMatch": {
										"matcher": "jsonpath", 
										"value": "$.payload",
										"doMatch": {
											"matcher": "jsonpath", 
											"value": "$.name",
											"doMatch": {
												"matcher": "exact", 
												"value": "John Doe"
											}
										}
									}
								}]
							}
						}
					]
				},
				"response": {
					"status": 200,
					"body": "jwt in form data matches with chaining",
					"encodedBody": false,
					"templated": false
				}
			}
		],
		"globalActions": {
			"delays": []
		}
	},
	"meta": {
		"schemaVersion": "v5.2",
		"hoverflyVersion": "v0.16.0",
		"timeExported": "2018-05-03T15:19:58+01:00"
	}
}`
