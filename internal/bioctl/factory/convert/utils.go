package convert

import (
	"encoding/json"
	"errors"
	"net/url"
	"reflect"
	"strings"

	"github.com/go-resty/resty/v2"
	"github.com/spf13/cast"
)

func AssignToHttpRequest(req interface{}, restyRequest *resty.Request) {
	reqV := reflect.ValueOf(req).Elem()
	reqT := reqV.Type()
	var pathParams map[string]string
	var queryParams url.Values
	var body map[string]interface{}

	for i := 0; i < reqV.NumField(); i++ {
		pathTag := reqT.Field(i).Tag.Get("path")
		if pathTag != "" {
			if pathParams == nil {
				pathParams = make(map[string]string)
			}
			pathParams[pathTag] = cast.ToString(reqV.Field(i).Interface())
		}
		queryTag := reqT.Field(i).Tag.Get("query")
		if queryTag != "" {
			if queryParams == nil {
				queryParams = make(url.Values)
			}
			if strings.Contains(queryTag, "omitempty") {
				if reqV.Field(i).IsZero() {
					continue
				}
				queryTag = strings.Split(queryTag, ",")[0]
			}
			queryParams[queryTag] = cast.ToStringSlice(reqV.Field(i).Interface())
		}
		bodyTag := reqT.Field(i).Tag.Get("json")
		if bodyTag != "" {
			if body == nil {
				body = make(map[string]interface{})
			}
			if strings.Contains(bodyTag, "omitempty") {
				if reqV.Field(i).IsZero() {
					continue
				}
				bodyTag = strings.Split(bodyTag, ",")[0]
			}
			body[bodyTag] = reqV.Field(i).Interface()
		}
	}
	restyRequest.SetPathParams(pathParams).SetQueryParamsFromValues(queryParams).SetBody(body)
}

func AssignFromHttpResponse(restyResp *resty.Response, resp interface{}) error {
	if restyResp.RawResponse != nil && restyResp.RawResponse.StatusCode >= 400 {
		return errors.New(string(restyResp.Body()))
	}
	err := json.Unmarshal(restyResp.Body(), resp)
	if err != nil {
		return errors.New(string(restyResp.Body()))
	}
	return nil
}

func RawBodyFromHttpResponse(restyResp *resty.Response) ([]byte, error) {
	if restyResp.RawResponse != nil && restyResp.RawResponse.StatusCode >= 400 {
		return nil, errors.New(string(restyResp.Body()))
	}
	return restyResp.Body(), nil
}
