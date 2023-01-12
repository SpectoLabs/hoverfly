package matchers

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"

	log "github.com/sirupsen/logrus"
)

var JWT = "jwt"

func JwtMatcher(data interface{}, toMatch string) bool {

	jwt, err := ParseJWT(toMatch)
	if err != nil {
		log.Errorf("Error occurred while fetching jwt token %s", err.Error())
		return false
	}
	return JsonPartialMatch(data, jwt)

}

func ParseJWT(str string) (string, error) {
	toMatchArr := strings.Split(str, ".")
	if len(toMatchArr) != 3 {
		//invalid json token
		return "", fmt.Errorf("invalid json token passed %s", str)
	}
	jwtToken := make(map[string]interface{})
	if header, err := GetDecodedJsonData(toMatchArr[0]); err == nil {
		jwtToken["header"] = header
	}
	if payload, err := GetDecodedJsonData(toMatchArr[1]); err == nil {
		jwtToken["payload"] = payload
	}
	if jsonByteArr, err := json.Marshal(jwtToken); err == nil {
		return string(jsonByteArr), nil
	} else {
		return "", err
	}

}

func GetDecodedJsonData(str string) (interface{}, error) {

	strByteArr, err := base64.RawURLEncoding.DecodeString(str)
	if err != nil {
		log.Errorf("Error occurred while decoding jwt part %s", str)
		return nil, err
	}
	var jsonData interface{}
	if err := json.Unmarshal(strByteArr, &jsonData); err != nil {
		return nil, err
	}
	return jsonData, nil
}
