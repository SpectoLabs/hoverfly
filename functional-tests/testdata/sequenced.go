package testdata

var Sequenced = `{
	"data": {
		"pairs": [{
			"request": {
				"path": [{
					"matcher": "exact",
					"value": "/a"
				}],
				"requiresState": {
					"sequence:1": "1"
				}
			},
			"response": {
				"status": 200,
				"body": "response 1a",
				"encodedBody": false,
				"headers": {
					"Content-Type": ["text/plain"],
					"Date": ["date"],
					"Hoverfly": ["Was-Here"]
				},
				"templated": false,
				"transitionsState": {
					"sequence:1": "2"
				}
			}
		}, {
			"request": {
				"path": [{
					"matcher": "exact",
					"value": "/a"
				}],
				"requiresState": {
					"sequence:1": "2"
				}
			},
			"response": {
				"status": 200,
				"body": "response 2a",
				"encodedBody": false,
				"headers": {
					"Content-Type": ["text/plain"],
					"Date": ["date"],
					"Hoverfly": ["Was-Here"]
				},
				"templated": false,
				"transitionsState": {
					"sequence:1": "3"
				}
			}
		}, {
			"request": {
				"path": [{
					"matcher": "exact",
					"value": "/a"
				}],
				"requiresState": {
					"sequence:1": "3"
				}
			},
			"response": {
				"status": 200,
				"body": "response 3a",
				"encodedBody": false,
				"headers": {
					"Content-Type": ["text/plain"],
					"Date": ["date"],
					"Hoverfly": ["Was-Here"]
				},
				"templated": false
			}
		}, {
			"request": {
				"path": [{
					"matcher": "exact",
					"value": "/b"
				}],
				"requiresState": {
					"sequence:2": "1"
				}
			},
			"response": {
				"status": 200,
				"body": "response 1b",
				"encodedBody": false,
				"headers": {
					"Content-Type": ["text/plain"],
					"Date": ["date"],
					"Hoverfly": ["Was-Here"]
				},
				"templated": false,
				"transitionsState": {
					"sequence:2": "2"
				}
			}
		}, {
			"request": {
				"path": [{
					"matcher": "exact",
					"value": "/b"
				}],
				"requiresState": {
					"sequence:2": "2"
				}
			},
			"response": {
				"status": 200,
				"body": "response 2b",
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
