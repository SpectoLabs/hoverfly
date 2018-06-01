package testdata

var Issue607 = `{
	"data": {
		"pairs": [
			{
				"request": {
					"path": [
						{
							"matcher": "exact",
							"value": "/billing/v1/servicequotes/123456"
						}
					],
					"method": [
						{
							"matcher": "exact",
							"value": "GET"
						}
					],
					"destination": [
						{
							"matcher": "exact",
							"value": "domain.com"
						}
					],
					"scheme": [
						{
							"matcher": "exact",
							"value": "https"
						}
					],
					"deprecatedQuery": [
						{
							"matcher": "exact",
							"value": "saleschannel=RETAIL"
						}
					],
					"body": [
						{
							"matcher": "json",
							"value": ""
						}
					],
					"headers": {
						"Accept": [
							{
								"matcher": "exact",
								"value": "application/json"
							}
						],
						"Activityid": [
							{
								"matcher": "exact",
								"value": "ChangeMSISDN_CR_PushtoBill(Get)-200"
							}
						],
						"Applicationid": [
							{
								"matcher": "exact",
								"value": "ACUI"
							}
						],
						"Authorization": [
							{
								"matcher": "exact",
								"value": "Bearer token"
							}
						],
						"Cache-Control": [
							{
								"matcher": "exact",
								"value": "no-cache"
							}
						],
						"Channelid": [
							{
								"matcher": "exact",
								"value": "RETAIL"
							}
						],
						"Content-Type": [
							{
								"matcher": "exact",
								"value": "application/json"
							}
						],
						"Interactionid": [
							{
								"matcher": "exact",
								"value": "123456787"
							}
						],
						"Senderid": [
							{
								"matcher": "exact",
								"value": "ACUI"
							}
						],
						"User-Agent": [
							{
								"matcher": "exact",
								"value": "curl/7.54.0"
							}
						],
						"Workflowid": [
							{
								"matcher": "exact",
								"value": "CHANGEMSISDN"
							}
						]
					}
				},
				"response": {
					"status": 200,
					"body": "",
					"encodedBody": false,
					"headers": {
						"Connection": [
							"keep-alive"
						],
						"Content-Type": [
							"application/json"
						],
						"Date": [
							"Tue, 13 Jun 2017 17:54:51 GMT"
						],
						"Hoverfly": [
							"Was-Here"
						],
						"Transfer-Encoding": [
							"chunked"
						]
					},
					"templated": false
				}
			}
		],
		"globalActions": {
			"delays": []
		}
	},
	"meta": {
		"schemaVersion": "v5",
		"hoverflyVersion": "v0.17.0",
		"timeExported": "2018-05-03T15:07:33+01:00"
	}
}`

var V3Issue607 = `{
	"data": {
		"pairs": [
			{
				"response": {
					"status": 200,
					"body": "",
					"encodedBody": false,
					"headers": {
						"Connection": [
							"keep-alive"
						],
						"Content-Type": [
							"application/json"
						],
						"Date": [
							"Tue, 13 Jun 2017 17:54:51 GMT"
						],
						"Hoverfly": [
							"Was-Here"
						],
						"Transfer-Encoding": [
							"chunked"
						]
					}
				},
				"request": {
					"path": {
						"exactMatch": "/billing/v1/servicequotes/123456"
					},
					"method": {
						"exactMatch": "GET"
					},
					"destination": {
						"exactMatch": "domain.com"
					},
					"scheme": {
						"exactMatch": "https"
					},
					"query": {
						"exactMatch": "saleschannel=RETAIL"
					},
					"body": {
						"jsonMatch": ""
					},
					"headers": {
						"Accept": [
							"application/json"
						],
						"Activityid": [
							"ChangeMSISDN_CR_PushtoBill(Get)-200"
						],
						"Applicationid": [
							"ACUI"
						],
						"Authorization": [
							"Bearer token"
						],
						"Cache-Control": [
							"no-cache"
						],
						"Channelid": [
							"RETAIL"
						],
						"Content-Type": [
							"application/json"
						],
						"Interactionid": [
							"123456787"
						],
						"Senderid": [
							"ACUI"
						],
						"User-Agent": [
							"curl/7.54.0"
						],
						"Workflowid": [
							"CHANGEMSISDN"
						]
					}
				}
			}
		],
		"globalActions": {
			"delays": []
		}
	},
	"meta": {
		"schemaVersion": "v3",
		"hoverflyVersion": "v0.12.0",
		"timeExported": "2017-06-13T10:55:12-07:00"
	}
}`
