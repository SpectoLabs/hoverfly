package testdata

var Issue607 = `{
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
