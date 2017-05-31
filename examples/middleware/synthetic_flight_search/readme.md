# Synthesizing https://qpx-express-demo.itasoftware.com/

## Usage:

Start hoverfly in synthesize mode and point to this middleware:
```
hoverfly -synthesize -middleware ./synthetic_flight_search.py
```

Execute client.py:
```    
export HTTPS_PROXY=http://localhost:8500/
python client.py --origin "AMS" --destination "NYC" --date "2017-15-15"
```

Try different destinations/origins:
```
python client.py --origin "AMS" --destination "LLL" --date "2017-15-15" | grep LLL
``` 

## More info

Example request:

```javascript
{
  "request": {
    "slice": [
      {
        "origin": "LGW",
        "destination": "AMS",
        "date": "2016-02-11",
        "maxStops": 0
      }
    ],
    "passengers": {
      "adultCount": 1,
      "infantInLapCount": 0,
      "infantInSeatCount": 0,
      "childCount": 0,
      "seniorCount": 0
    },
    "solutions": 2,
    "refundable": false
  }
}
```

Example response:

```javascript
{
 "kind": "qpxExpress#tripsSearch",
 "trips": {
  "kind": "qpxexpress#tripOptions",
  "requestId": "SYzLMFMFPCrebUp5H0NaGL",
  "data": {
   "kind": "qpxexpress#data",
   "airport": [
    {
     "kind": "qpxexpress#airportData",
     "code": "AMS",
     "city": "AMS",
     "name": "Amsterdam Schiphol Airport"
    },
    {
     "kind": "qpxexpress#airportData",
     "code": "LGW",
     "city": "LON",
     "name": "London Gatwick"
    }
   ],
   "city": [
    {
     "kind": "qpxexpress#cityData",
     "code": "AMS",
     "name": "Amsterdam"
    },
    {
     "kind": "qpxexpress#cityData",
     "code": "LON",
     "name": "London"
    }
   ],
   "aircraft": [
    {
     "kind": "qpxexpress#aircraftData",
     "code": "319",
     "name": "Airbus A319"
    },
    {
     "kind": "qpxexpress#aircraftData",
     "code": "320",
     "name": "Airbus A320"
    }
   ],
   "tax": [
    {
     "kind": "qpxexpress#taxData",
     "id": "GB_001",
     "name": "United Kingdom Air Passengers Duty"
    },
    {
     "kind": "qpxexpress#taxData",
     "id": "UB",
     "name": "United Kingdom Passenger Service Charge"
    }
   ],
   "carrier": [
    {
     "kind": "qpxexpress#carrierData",
     "code": "BA",
     "name": "British Airways p.l.c."
    }
   ]
  },
  "tripOption": [
   {
    "kind": "qpxexpress#tripOption",
    "saleTotal": "GBP47.27",
    "id": "OAcAQw8rr9MNhwQoBntUKJ001",
    "slice": [
     {
      "kind": "qpxexpress#sliceInfo",
      "duration": 75,
      "segment": [
       {
        "kind": "qpxexpress#segmentInfo",
        "duration": 75,
        "flight": {
         "carrier": "BA",
         "number": "2762"
        },
        "id": "GStLakphRYJX3LbK",
        "cabin": "COACH",
        "bookingCode": "O",
        "bookingCodeCount": 1,
        "marriedSegmentGroup": "0",
        "leg": [
         {
          "kind": "qpxexpress#legInfo",
          "id": "LgJHYCVgG0AiE1PH",
          "aircraft": "320",
          "arrivalTime": "2016-02-11T18:05+01:00",
          "departureTime": "2016-02-11T15:50+00:00",
          "origin": "LGW",
          "destination": "AMS",
          "originTerminal": "N",
          "duration": 75,
          "mileage": 226,
          "meal": "Snack or Brunch"
         }
        ]
       }
      ]
     }
    ],
    "pricing": [
     {
      "kind": "qpxexpress#pricingInfo",
      "fare": [
       {
        "kind": "qpxexpress#fareInfo",
        "id": "A855zsItBCELBykaeqeBDQb5hPZQIOtkOZ8uDq0lD5VU",
        "carrier": "BA",
        "origin": "LON",
        "destination": "AMS",
        "basisCode": "OV1KO"
       }
      ],
      "segmentPricing": [
       {
        "kind": "qpxexpress#segmentPricing",
        "fareId": "A855zsItBCELBykaeqeBDQb5hPZQIOtkOZ8uDq0lD5VU",
        "segmentId": "GStLakphRYJX3LbK"
       }
      ],
      "baseFareTotal": "GBP22.00",
      "saleFareTotal": "GBP22.00",
      "saleTaxTotal": "GBP25.27",
      "saleTotal": "GBP47.27",
      "passengers": {
       "kind": "qpxexpress#passengerCounts",
       "adultCount": 1
      },
      "tax": [
       {
        "kind": "qpxexpress#taxInfo",
        "id": "UB",
        "chargeType": "GOVERNMENT",
        "code": "UB",
        "country": "GB",
        "salePrice": "GBP12.27"
       },
       {
        "kind": "qpxexpress#taxInfo",
        "id": "GB_001",
        "chargeType": "GOVERNMENT",
        "code": "GB",
        "country": "GB",
        "salePrice": "GBP13.00"
       }
      ],
      "fareCalculation": "LON BA AMS 33.71OV1KO NUC 33.71 END ROE 0.652504 FARE GBP 22.00 XT 13.00GB 12.27UB",
      "latestTicketingTime": "2016-01-11T23:59-05:00",
      "ptc": "ADT"
     }
    ]
   },
   {
    "kind": "qpxexpress#tripOption",
    "saleTotal": "GBP62.27",
    "id": "OAcAQw8rr9MNhwQoBntUKJ002",
    "slice": [
     {
      "kind": "qpxexpress#sliceInfo",
      "duration": 80,
      "segment": [
       {
        "kind": "qpxexpress#segmentInfo",
        "duration": 80,
        "flight": {
         "carrier": "BA",
         "number": "2758"
        },
        "id": "GW8rUjsDA234DdHV",
        "cabin": "COACH",
        "bookingCode": "Q",
        "bookingCodeCount": 9,
        "marriedSegmentGroup": "0",
        "leg": [
         {
          "kind": "qpxexpress#legInfo",
          "id": "Lp08eKxnXnyWfJo4",
          "aircraft": "319",
          "arrivalTime": "2016-02-11T10:05+01:00",
          "departureTime": "2016-02-11T07:45+00:00",
          "origin": "LGW",
          "destination": "AMS",
          "originTerminal": "N",
          "duration": 80,
          "mileage": 226,
          "meal": "Snack or Brunch"
         }
        ]
       }
      ]
     }
    ],
    "pricing": [
     {
      "kind": "qpxexpress#pricingInfo",
      "fare": [
       {
        "kind": "qpxexpress#fareInfo",
        "id": "AslXz8S1h3mMcnYUQ/v0Zt0p9Es2hj8U0We0xFAU1qDE",
        "carrier": "BA",
        "origin": "LON",
        "destination": "AMS",
        "basisCode": "QV1KO"
       }
      ],
      "segmentPricing": [
       {
        "kind": "qpxexpress#segmentPricing",
        "fareId": "AslXz8S1h3mMcnYUQ/v0Zt0p9Es2hj8U0We0xFAU1qDE",
        "segmentId": "GW8rUjsDA234DdHV"
       }
      ],
      "baseFareTotal": "GBP37.00",
      "saleFareTotal": "GBP37.00",
      "saleTaxTotal": "GBP25.27",
      "saleTotal": "GBP62.27",
      "passengers": {
       "kind": "qpxexpress#passengerCounts",
       "adultCount": 1
      },
      "tax": [
       {
        "kind": "qpxexpress#taxInfo",
        "id": "UB",
        "chargeType": "GOVERNMENT",
        "code": "UB",
        "country": "GB",
        "salePrice": "GBP12.27"
       },
       {
        "kind": "qpxexpress#taxInfo",
        "id": "GB_001",
        "chargeType": "GOVERNMENT",
        "code": "GB",
        "country": "GB",
        "salePrice": "GBP13.00"
       }
      ],
      "fareCalculation": "LON BA AMS 56.70QV1KO NUC 56.70 END ROE 0.652504 FARE GBP 37.00 XT 13.00GB 12.27UB",
      "latestTicketingTime": "2016-01-11T23:59-05:00",
      "ptc": "ADT"
     }
    ]
   }
  ]
 }
}
```