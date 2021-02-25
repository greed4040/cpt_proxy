package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/go-redis/redis"
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

	fmt.Println("##type of resp:", reflect.TypeOf(resp), "##status code:", resp.StatusCode)
	fmt.Println("##resp:", resp)

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
	if resp.StatusCode == http.StatusInternalServerError {
		bodyBytes, err := ioutil.ReadAll(resp.Body)
		checkErr(err)
		bodyString = string(bodyBytes)
		log.Info(bodyString)
		fmt.Println("##bodystring:", bodyString)
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
	Key     string        `json:"key"`
	Tm      int           `json:"tm"`
}

func btc_get_blockCount(t externalCmd) string {
	local_params := make([]interface{}, 0)
	req_data := req{Json: "1.0", Id: "222", Method: t.Command, Params: local_params}
	resp := master_req(req_data)
	return resp
}
func btc_get_newAddress(t externalCmd) string {
	local_params := make([]interface{}, 0)
	req_data := req{Json: "1.0", Id: "222", Method: t.Command, Params: local_params}
	resp := master_req(req_data)
	return resp
}
func btc_get_raw_mempool(t externalCmd) string {
	local_params := make([]interface{}, 0)
	req_data := req{Json: "1.0", Id: "222", Method: t.Command, Params: local_params}
	resp := master_req(req_data)
	return resp
}
func btc_list_unspent(t externalCmd) string {
	local_params := make([]interface{}, 0)
	req_data := req{Json: "1.0", Id: "222", Method: t.Command, Params: local_params}
	resp := master_req(req_data)
	return resp
}
func btc_dumpprivkey(t externalCmd) string {
	local_params := make([]interface{}, 0)
	fmt.Println(reflect.TypeOf(local_params), len(local_params))
	for a := range t.Params {
		local_params = append(local_params, t.Params[a])
	}
	req_data := req{Json: "1.0", Id: "222", Method: t.Command, Params: local_params}
	resp := master_req(req_data)
	return resp
}
func btc_create_rawtransaction(t externalCmd) string {
	local_params := make([]interface{}, 0)
	fmt.Println("#btc_create_rawtransaction", t.Params)
	fmt.Println("#btc_create_rawtransaction type:", reflect.TypeOf(t.Params))

	for a := range t.Params {
		fmt.Println("%", t.Params[a], reflect.TypeOf(t.Params[a]))
		//if str, ok := t.Params[a].(string); ok {
		local_params = append(local_params, t.Params[a])
		//}
	}

	fmt.Println(local_params, reflect.TypeOf(local_params), len(local_params))
	req_data := req{Json: "1.0", Id: "333", Method: t.Command, Params: local_params}
	resp := master_req(req_data)

	//resp = "##btc_create_rawtransaction##"
	return resp
}
func btc_sign_rawTransactionWithKey(t externalCmd) string {
	local_params := make([]interface{}, 0)
	fmt.Println("#btc_create_rawtransaction", t.Params)
	fmt.Println("#btc_create_rawtransaction type:", reflect.TypeOf(t.Params))

	for a := range t.Params {
		fmt.Println("%", t.Params[a], reflect.TypeOf(t.Params[a]))
		//if str, ok := t.Params[a].(string); ok {
		local_params = append(local_params, t.Params[a])
		//}
	}

	fmt.Println(local_params, reflect.TypeOf(local_params), len(local_params))
	req_data := req{Json: "1.0", Id: "333", Method: t.Command, Params: local_params}
	resp := master_req(req_data)

	//resp = "##btc_create_rawtransaction##"
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
func btc_decode_rawTransaction(t externalCmd) string {
	local_params := make([]interface{}, 0)
	fmt.Println("#btc_decode_rawTransaction", t.Params)
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
func btc_get_txOut(t externalCmd) string {
	local_params := make([]interface{}, 0)
	fmt.Println("#btc_get_txOut:", t.Params, len(t.Params), reflect.TypeOf(t.Params))

	for a := range t.Params {
		fmt.Println("##", t.Params[a], reflect.TypeOf(t.Params[a]))
		switch v := t.Params[a].(type) {
		default:
			fmt.Printf("unexpected type %T", v)
		case string:
			local_params = append(local_params, t.Params[a])
		case float64:
			var aaa = fmt.Sprintf("%v", t.Params[a])
			i, _ := strconv.Atoi(aaa)
			local_params = append(local_params, i)
		}
	}

	req_data := req{Json: "1.0", Id: "222", Method: t.Command, Params: local_params}
	resp := master_req(req_data)
	return resp
}
func btc_get_rawTransaction(t externalCmd) string {
	local_params := make([]interface{}, 0)
	fmt.Println("#btc_get_rawTransaction", t.Params)
	fmt.Println(reflect.TypeOf(t.Params))

	for a := range t.Params {
		fmt.Println("%", t.Params[a], reflect.TypeOf(t.Params[a]))
		//if str, ok := t.Params[a].(string); ok {
		//conv, err := strconv.Atoi(str)
		//if err != nil {
		//	log.Warn(err)
		//}
		//local_params = append(local_params, conv)
		//}
		local_params = append(local_params, fmt.Sprintf("%v", t.Params[a]))
	}

	req_data := req{Json: "1.0", Id: "222", Method: t.Command, Params: local_params}
	resp := master_req(req_data)
	return resp
}
func btc_send_rawTransaction(t externalCmd) string {
	local_params := make([]interface{}, 0)
	fmt.Println("#btc_send_rawTransaction", t.Params)
	fmt.Println(reflect.TypeOf(t.Params))

	for a := range t.Params {
		fmt.Println("%", t.Params[a], reflect.TypeOf(t.Params[a]))
		local_params = append(local_params, fmt.Sprintf("%v", t.Params[a]))
	}

	req_data := req{Json: "1.0", Id: "222", Method: t.Command, Params: local_params}
	resp := master_req(req_data)
	return resp
}
func btc_get_data3(rw http.ResponseWriter, req *http.Request) {
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

func btc_get_data2(rw http.ResponseWriter, req *http.Request) {
	fmt.Println("------------------------------------------")
	fmt.Println("------------processing request------------")
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
	case "getrawmempool":
		{
			resp = "this is answer from go server:getrawmempool"
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
	fmt.Println("----------processing request done---------")
	fmt.Println("------------------------------------------")

}

func randSt(ln int) string {
	var output strings.Builder
	var charSet = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890")
	for i := 0; i < ln; i++ {
		random := rand.Intn(len(charSet))
		randomChar := charSet[random]
		output.WriteRune(randomChar)
	}
	return output.String()
}

func btc_get_data(rw http.ResponseWriter, req *http.Request) {
	t := time.Now()
	rnds := randSt(5)
	fmt.Println("begin", t.Format("2006/01/02 15:04:05.1234"), rnds, reflect.TypeOf(rnds))

	body, _ := ioutil.ReadAll(req.Body)
	fmt.Printf("#btc_get_data recieved_body: %s", body)

	rw.Header().Set("Content-Type", "application/json")
	rw.Header().Set("Access-Control-Allow-Origin", "*")
	rw.Header().Set("Access-Control-Allow-Credentials", "true")
	rw.Header().Set("Access-Control-Allow-Methods", "GET,HEAD,OPTIONS,POST,PUT")
	rw.Header().Set("Access-Control-Allow-Headers", "Access-Control-Allow-Headers, Origin,Accept, X-Requested-With, Content-Type, Access-Control-Request-Method, Access-Control-Request-Headers")

	var resp = "this is answer from go server:getrawmempool"

	js, err := json.Marshal(resp)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
	rw.Write(js)
	fmt.Println("end", rnds)
}

func swich_calls(t externalCmd) string {
	fmt.Println("#swich_calls:", t.Command)
	var resp = ""
	switch t.Command {
	case "getblockcount":
		{
			resp = btc_get_blockCount(t)
		}
	case "getnewaddress":
		{
			resp = btc_get_newAddress(t)
		}
	case "dumpprivkey":
		{
			resp = btc_dumpprivkey(t)
		}
	case "getblockhash":
		{
			resp = btc_get_blockHash(t)
		}
	case "gettxout":
		{
			resp = btc_get_txOut(t)
		}
	case "getrawtransaction":
		{
			resp = btc_get_rawTransaction(t)
		}
	case "decoderawtransaction":
		{
			resp = btc_decode_rawTransaction(t)
		}
	case "getblock":
		{
			resp = btc_get_block(t)
		}
	case "getrawmempool":
		{
			resp = btc_get_raw_mempool(t)
		}
	case "listunspent":
		{
			resp = btc_list_unspent(t)
		}
	case "createrawtransaction":
		{
			resp = btc_create_rawtransaction(t)
		}
	case "signrawtransactionwithkey":
		{
			resp = btc_sign_rawTransactionWithKey(t)
		}
	case "sendrawtransaction":
		{
			resp = btc_send_rawTransaction(t)
		}
	default:
		resp = "Command not supported"
	}

	return resp
}
func btc_getdt(rw http.ResponseWriter, req *http.Request) {
	body, _ := ioutil.ReadAll(req.Body)
	fmt.Printf("#btc_get_data recieved_body: %s\n", body)

	rw.Header().Set("Content-Type", "application/json")
	rw.Header().Set("Access-Control-Allow-Origin", "*")
	rw.Header().Set("Access-Control-Allow-Credentials", "true")
	rw.Header().Set("Access-Control-Allow-Methods", "GET,HEAD,OPTIONS,POST,PUT")
	rw.Header().Set("Access-Control-Allow-Headers", "Access-Control-Allow-Headers, Origin,Accept, X-Requested-With, Content-Type, Access-Control-Request-Method, Access-Control-Request-Headers")

	var t externalCmd
	err := json.Unmarshal(body, &t)
	if err != nil {
		fmt.Println("#btc_get_data Unmarshal error")
		log.Error(err)
	}
	//log.Info(t)
	fmt.Println("#auth key:", t.Key)
	opts := &redis.Options{Addr: "localhost:6379", Password: "", DB: 9}
	client := redis.NewClient(opts)
	a, err := client.HGetAll(t.Key).Result()
	checkErr(err)
	fmt.Println("redis result:", a, len(a))

	answer_map := make(map[string]interface{})
	js, _ := json.Marshal(answer_map)
	if len(a) > 0 {
		fmt.Println("#btc_get_data result", "cmd:", t.Command, "params:", t.Params)

		var t2 externalCmd
		t2.Params = t.Params
		t2.Command = t.Command
		var resp = swich_calls(t2)

		fmt.Println("resp", reflect.TypeOf(resp))
		js, err = json.Marshal(resp)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}
	} else {
		answer_map := `{"result": "Unauthorized", "error": true, "id": "222"}`
		fmt.Println("answer map:", reflect.TypeOf(answer_map), answer_map)
		js, err = json.Marshal(answer_map)
		fmt.Println("answer map:", reflect.TypeOf(js))
	}
	rw.Write(js)
	/**/
}

func main() {
	fmt.Println("Starting.")
	fmt.Println("Server is listening...")
	//http.HandleFunc("/btc_get_test", btc_get_test)
	//http.HandleFunc("/btc_get_data", btc_get_data)
	//http.HandleFunc("/btc_get_data2", btc_get_data2)
	http.HandleFunc("/btc_getdata", btc_getdt)
	err := http.ListenAndServeTLS(":9292", "fullchain.pem", "privkey.pem", nil)
	if err != nil {
		log.Fatal(err)
	}
}
