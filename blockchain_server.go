package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"reflect"
	"strconv"
	"strings"

	"github.com/labstack/gommon/log"
)

func checkErr(err error) {
	if err != nil {
		fmt.Println(err)
	}
}

type req struct {
	Json   string        `json:"json"`
	Id     string        `json:"id"`
	Method string        `json:"method"`
	Params []interface{} `json:"params"`
}

func master_req(inp_req req) string {
	url := "http://176.126.167.107:8332"
	client := &http.Client{}

	fmt.Println("#this is request struct:", inp_req, reflect.TypeOf(inp_req))
	mJson, marchal_err := json.Marshal(inp_req)
	fmt.Println("#this is marshalled json:", string(mJson), len(string(mJson)), marchal_err)
	fmt.Println("---")

	contentReader := bytes.NewReader(mJson)
	req, _ := http.NewRequest("POST", url, contentReader)
	req.Header.Set("Content-Type", "text/plain")
	req.SetBasicAuth("bitcoin", "gyqtxm7oELcqAlO3XBVtoghcSNXgVLCCGHKDL6wstOE")
	resp, reqErr := client.Do(req)
	checkErr(reqErr)
	defer resp.Body.Close()

	fmt.Println(reflect.TypeOf(resp), resp.StatusCode)
	fmt.Println(resp)

	var bodyString string
	if resp.StatusCode == http.StatusOK {
		bodyBytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatal(err)
		}
		bodyString = string(bodyBytes)
		log.Info(bodyString)
		fmt.Println(bodyString)
	}
	return bodyString
}

func basic_test() string {
	url := "http://176.126.167.107:8332"
	client := &http.Client{}

	my_params := make([]interface{}, 0)
	req_data := req{Json: "1.0", Id: "222", Method: "getdifficulty", Params: my_params}

	fmt.Println("#this is request struct:", req_data, reflect.TypeOf(req_data))
	mJson, marchal_err := json.Marshal(req_data)
	fmt.Println("#this is marshalled json:", string(mJson), len(string(mJson)), marchal_err)
	fmt.Println("---")

	contentReader := bytes.NewReader(mJson)
	req, _ := http.NewRequest("POST", url, contentReader)
	req.Header.Set("Content-Type", "text/plain")
	req.SetBasicAuth("bitcoin", "gyqtxm7oELcqAlO3XBVtoghcSNXgVLCCGHKDL6wstOE")
	resp, reqErr := client.Do(req)
	checkErr(reqErr)
	defer resp.Body.Close()

	fmt.Println(reflect.TypeOf(resp), resp.StatusCode)
	fmt.Println(resp)

	var bodyString string
	if resp.StatusCode == http.StatusOK {
		bodyBytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatal(err)
		}
		bodyString = string(bodyBytes)
		log.Info(bodyString)
		fmt.Println(bodyString)
	}
	return bodyString
}

func test_inp_param_method(inp_params string, inp_method string) {
	url := "http://176.126.167.107:8332"
	client := &http.Client{}

	my_params := make([]interface{}, 0)
	req_data := req{Json: "1.0", Id: "222", Method: inp_method, Params: my_params}

	mJson, _ := json.Marshal(req_data)
	mJsonS := string(mJson)
	new_aaa := fmt.Sprintf(`"params":[%s]`, inp_params)
	mJsonS = strings.Replace(mJsonS, `"params":""`, new_aaa, 1)
	fmt.Println("#this is replaced json:", mJsonS, "#", new_aaa)
	fmt.Println(inp_params)
	fmt.Println("---")

	contentReader := bytes.NewReader([]byte(mJsonS))
	req, _ := http.NewRequest("POST", url, contentReader)
	req.Header.Set("Content-Type", "text/plain")
	req.SetBasicAuth("bitcoin", "gyqtxm7oELcqAlO3XBVtoghcSNXgVLCCGHKDL6wstOE")
	resp, reqErr := client.Do(req)
	checkErr(reqErr)
	defer resp.Body.Close()

	fmt.Println(reflect.TypeOf(resp), resp.StatusCode)
	fmt.Println(resp)

	if resp.StatusCode == http.StatusOK {
		bodyBytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatal(err)
		}
		bodyString := string(bodyBytes)
		log.Info(bodyString)
		fmt.Println(bodyString)
	}
}

type the_data struct {
	Main []string
	Opt  string
}

func gettxoutproof() {
	my_method := "gettxoutproof"
	d1 := "f5710135e557b002d803a1e742b7b24170679c63120eb84d4e14dad8836bf796"
	d2 := "00000000000000000007a714a0a1271dc08b690c7a1d94e57967a5cbb82ce160"

	aaa, _ := json.Marshal([]string{d1})
	//fmt.Println(string(aaa))
	a := fmt.Sprintf(`%s, "%s"`, string(aaa), d2)
	fmt.Println("#this is request:", a, reflect.TypeOf(a))
	fmt.Println("#method:", my_method)

	//fmt.Println("#aaa:", a, "#", my_method)
	test_inp_param_method(a, my_method)
}

func extract_http_vars(params []string, v *http.Request) map[string]string {
	m := make(map[string]string)
	fmt.Println("#extract_http_vars:", v.URL.Query())
	for b := 0; b < len(params); b++ {
		tmp, ok := v.URL.Query()[params[b]]
		if !ok || len(tmp[0]) < 1 {
			log.Warn("some params are missing:", params[b])
		}
		m[params[b]] = tmp[0]
	}
	return m
}

func extract_post_vars(params []string, v *http.Request) map[string]string {
	m := make(map[string]string)
	for b := 0; b < len(params); b++ {
		m[params[b]] = v.Form.Get(params[b])
	}
	return m
}

func btc_get_test(rw http.ResponseWriter, req *http.Request) {
	fmt.Println("#btc_get req body:", req.Body)

	err := req.ParseForm()
	if err != nil {
		log.Warn(err)
		return
	}
	for k, v := range req.Header {
		fmt.Println("#headers:", k, v)
	}
	command := req.PostForm.Get("command")
	params := req.PostForm.Get("params")

	fmt.Println("#btc_get cmd:", command)
	fmt.Println("#btc_get par:", params)

	rw.Header().Set("Content-Type", "application/json")
	rw.Header().Set("Access-Control-Allow-Origin", "*")
	rw.Header().Set("Access-Control-Allow-Credentials", "true")
	rw.Header().Set("Access-Control-Allow-Methods", "GET,HEAD,OPTIONS,POST,PUT")
	rw.Header().Set("Access-Control-Allow-Headers", "Access-Control-Allow-Headers, Origin,Accept, X-Requested-With, Content-Type, Access-Control-Request-Method, Access-Control-Request-Headers")

	js, err := json.Marshal("hello")
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
	rw.Write(js)
}

type externalCmd struct {
	Command string        `json:"command"`
	Params  []interface{} `json:"params"`
}

func btc_get_blockCount(t externalCmd) string {
	local_params := make([]interface{}, 0)
	req_data := req{Json: "1.0", Id: "222", Method: t.Command, Params: local_params}
	resp := master_req(req_data)
	return resp
}
func btc_get_bestBlockHash(t externalCmd) string {
	local_params := make([]interface{}, 0)
	req_data := req{Json: "1.0", Id: "222", Method: t.Command, Params: local_params}
	resp := master_req(req_data)
	return resp
}
func btc_get_blockHash(t externalCmd) string {
	local_params := make([]interface{}, 0)
	fmt.Println("#btc_get_blockHash", t.Params)
	fmt.Println(reflect.TypeOf(t.Params))

	for a := range t.Params {
		fmt.Println("%", t.Params[a], reflect.TypeOf(t.Params[a]))
		if str, ok := t.Params[a].(string); ok {
			conv, err := strconv.Atoi(str)
			if err != nil {
				log.Warn(err)
			}
			local_params = append(local_params, conv)
		}
	}

	req_data := req{Json: "1.0", Id: "222", Method: t.Command, Params: local_params}
	resp := master_req(req_data)
	return resp
}
func btc_get_blockHeader(t externalCmd) string {
	local_params := make([]interface{}, 0)
	fmt.Println("#btc_get_blockHeader", t.Params)
	fmt.Println(reflect.TypeOf(t.Params))

	for a := range t.Params {
		fmt.Println("%", t.Params[a], reflect.TypeOf(t.Params[a]))
		if str, ok := t.Params[a].(string); ok {
			local_params = append(local_params, str)
		}
	}

	req_data := req{Json: "1.0", Id: "222", Method: t.Command, Params: local_params}
	resp := master_req(req_data)
	return resp
}
func btc_get_block(t externalCmd) string {
	local_params := make([]interface{}, 0)
	fmt.Println("#btc_get_block", t.Params)
	fmt.Println(reflect.TypeOf(t.Params))

	for a := range t.Params {
		fmt.Println("%", t.Params[a], reflect.TypeOf(t.Params[a]))
		if str, ok := t.Params[a].(string); ok {
			local_params = append(local_params, str)
		}
	}

	req_data := req{Json: "1.0", Id: "222", Method: t.Command, Params: local_params}
	resp := master_req(req_data)
	return resp
}
func btc_get_data(rw http.ResponseWriter, req *http.Request) {
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		fmt.Println("#btc_get_data ioutil.ReadAll error")
		log.Error(err)
	}
	//log.Info("#btc_get2", string(body))
	fmt.Printf("#btc_get_data recieved_body: %s", body)

	var t externalCmd
	err = json.Unmarshal(body, &t)
	if err != nil {
		fmt.Println("#btc_get_data Unmarshal error")
		log.Error(err)
	}
	//log.Info(t)
	fmt.Println("#btc_get_data result", "cmd:", t.Command, "params:", t.Params)

	var resp string // := basic_test()
	switch t.Command {
	case "getblockcount":
		{
			resp = btc_get_blockCount(t)
		}
	case "getbestblockhash":
		{
			resp = btc_get_bestBlockHash(t)
		}
	case "getblockhash":
		{
			resp = btc_get_blockHash(t)
		}
	case "getblockheader":
		{
			resp = btc_get_blockHeader(t)
		}
	case "getblock":
		{
			resp = btc_get_block(t)
		}
	}

	rw.Header().Set("Content-Type", "application/json")
	rw.Header().Set("Access-Control-Allow-Origin", "*")
	rw.Header().Set("Access-Control-Allow-Credentials", "true")
	rw.Header().Set("Access-Control-Allow-Methods", "GET,HEAD,OPTIONS,POST,PUT")
	rw.Header().Set("Access-Control-Allow-Headers", "Access-Control-Allow-Headers, Origin,Accept, X-Requested-With, Content-Type, Access-Control-Request-Method, Access-Control-Request-Headers")

	js, err := json.Marshal(resp)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
	rw.Write(js)
}

func main() {
	fmt.Println("Starting.")
	fmt.Println("Server is listening...")
	http.HandleFunc("/btc_get_test", btc_get_test)
	http.HandleFunc("/btc_get_data", btc_get_data)

	err := http.ListenAndServeTLS(":9292", "fullchain.pem", "privkey.pem", nil)
	if err != nil {
		log.Fatal(err)
	}
}
