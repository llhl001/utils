package iaas

import (
	"encoding/json"
	"fmt"
	"reflect"
	"sort"
	"strings"
	"time"
	"utils/conv"
	vhttp "utils/net/http"
	verror "utils/os/error"
	qcutil "utils/qingcloud"
	"utils/util"
)

// Send 发送请求到Iaas
// conf 包含配置：console_key_id,console_secrect_key,console_uri,host,port
func Send(method string, params map[string]interface{}, conf map[string]interface{}) (interface{}, error) {
	_method := strings.ToLower(method)
	urlParams, _, data, err := Signature(_method, conf["console_uri"].(string), conf["console_key_id"].(string), conf["console_secrect_key"].(string), params)
	if err != nil {
		return nil, err
	}

	headers := map[string]string{}
	if _method == "post" || _method == "put" {
		headers["Content-Type"] = "'application/x-www-form-urlencoded'"
		headers["Accept"] = "text/plain"
		headers["Connection"] = "Keep-Alive"
		headers["Content-Length"] = string(len(data))
	}

	var url string = fmt.Sprintf("http://%s:%s%s", conf["host"].(string), conf["port"].(string), conf["console_uri"].(string))

	var resp interface{}
	if _method == "get" {
		vhttp.Get2(url+"?"+urlParams, &resp, headers)
	} else if _method == "post" {
		vhttp.Post2(url+"?"+urlParams, data, &resp, headers)
	} else if _method == "put" {
		vhttp.Put(url+"?"+urlParams, data, &resp, headers)
	} else if _method == "delete" {
		vhttp.Delete(url+"?"+urlParams, &resp, headers)
	}

	return resp, nil
}

func Signature(method, uri, ak, sk string, params map[string]interface{}) (string, string, string, error) {
	_method := strings.ToLower(method)
	// _params := url.Values{}
	_params := map[string]interface{}{}

	var _data string = ""
	if _method == "get" || _method == "delete" {
		_params = params
	} else {
		bData, err := json.Marshal(params)
		if err != nil {
			return "", "", "", verror.New("parameter parsing error")
		}
		_data = string(bData)
	}

	// time_stamp := time.Now() //time.Now().UTC().Format(time.RFC3339)
	time_stamp := time.Now()
	_params["time_stamp"] = util.TimeToString(time_stamp, "ISO 8601")                  // TimeToString(time_stamp, "ISO 8601")
	_params["expires"] = util.TimeToString(time_stamp.Add(10*time.Second), "ISO 8601") // time.Now().Add(time.Hour).Format("2006-01-02T15:04:05Z")
	_params["signature_version"] = "1"
	_params["signature_method"] = "HmacSHA256"
	_params["access_key_id"] = ak

	keys := []string{}
	for key := range _params {
		keys = append(keys, key)
	}

	sort.Strings(keys)
	parts := []string{}
	for _, key := range keys {
		v := _params[key]
		if v != nil {
			_v := ""
			switch reflect.TypeOf(v).String() {
			case "string":
				_v = qcutil.QueryEscape(v.(string))
				parts = append(parts, key+"="+_v)
			case "[]interface {}":
				for i, val := range v.([]interface{}) {
					_v = qcutil.QueryEscape(conv.String(val))
					parts = append(parts, key+"."+conv.String(i+1)+"="+_v)
				}
			case "[]string":
				for i, val := range v.([]string) {
					_v = qcutil.QueryEscape(val)
					parts = append(parts, key+"."+conv.String(i+1)+"="+_v)
				}
			default:
				_v = qcutil.QueryEscape(conv.String(v))
				parts = append(parts, key+"="+_v)
			}
		}
	}
	urlParams := strings.Join(parts, "&")
	signature := qcutil.Get_iaas_authorization(sk, _method, uri, urlParams)
	urlParams = urlParams + "&signature=" + signature

	return urlParams, signature, _data, nil
}
