package testdata

var Sequenced = `{
	"data": {
		"pairs": [{
			"request": {
				"path": [{
					"matcher": "exact",
					"value": "/"
				}],
				"requiresState": {
					"sequence:0": "1"
				}
			},
			"response": {
				"status": 200,
				"body": "response 1",
				"encodedBody": false,
				"headers": {
					"Content-Type": ["text/plain"],
					"Date": ["date"],
					"Hoverfly": ["Was-Here"]
				},
				"templated": false,
				"transitionsState": {
					"sequence:0": "2"
				}
			}
		}, {
			"request": {
				"path": [{
					"matcher": "exact",
					"value": "/"
				}],
				"requiresState": {
					"sequence:0": "2"
				}
			},
			"response": {
				"status": 200,
				"body": "response 2",
				"encodedBody": false,
				"headers": {
					"Content-Type": ["text/plain"],
					"Date": ["date"],
					"Hoverfly": ["Was-Here"]
				},
				"templated": false,
				"transitionsState": {
					"sequence:0": "3"
				}
			}
		}, {
			"request": {
				"path": [{
					"matcher": "exact",
					"value": "/"
				}],
				"requiresState": {
					"sequence:0": "3"
				}
			},
			"response": {
				"status": 200,
				"body": "response 3",
				"encodedBody": false,
				"headers": {
					"Content-Type": ["text/plain"],
					"Date": ["date"],
					"Hoverfly": ["Was-Here"]
				},
				"templated": false
			}
		}],
		"globalActions": {
			"delays": []
		}
	},
	"meta": {
		"schemaVersion": "v5",
		"hoverflyVersion": "v0.17.0",
		"timeExported": "2018-05-24T14:20:55+01:00"
	}
}`
