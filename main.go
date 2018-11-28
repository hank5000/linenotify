// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"io/ioutil"


	"html/template"
)

var clientID string
var clientSecret string
var callbackURL string
var token string

func main() {
	http.HandleFunc("/callback", callbackHandler)
	http.HandleFunc("/notify", notifyHandler)
	http.HandleFunc("/auth", authHandler)
	http.HandleFunc("/connect",connectHandler)
	http.HandleFunc("/connecting",connectingHandler)

	clientID = os.Getenv("ClientID")
	clientSecret = os.Getenv("ClientSecret")
	callbackURL = os.Getenv("CallbackURL")
	port := os.Getenv("PORT")
	fmt.Printf("ENV port:%s, cid:%s csecret:%s\n", port, clientID, clientSecret)
	addr := fmt.Sprintf(":%s", port)
	http.ListenAndServe(addr, nil)
}

func notifyHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm() // Populates request.Form
	msg := r.Form.Get("msg")
	user_token := r.Form.Get("token")
	fmt.Printf("Get msg=%s\n", msg)
	b, errr := ioutil.ReadAll(r.Body)
	if errr != nil {
		panic(errr)
	}
	fmt.Printf("Body=%s\n",b)

	data := url.Values{}
	data.Add("message", msg)

	byt, err := apiCall("POST", apiNotify, data, user_token)
	fmt.Println("ret:", string(byt), " err:", err)

	res := newTokenResponse(byt)
	fmt.Println("result:", res)
	token = res.AccessToken
	w.Write(byt)
}

func callbackHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm() // Populates request.Form
	code := r.Form.Get("code")
	state := r.Form.Get("state")
	fmt.Printf("Get code=%s, state=%s \n", code, state)

	data := url.Values{}
	data.Add("grant_type", "authorization_code")
	data.Add("code", code)
	data.Add("redirect_uri", callbackURL)
	data.Add("client_id", clientID)
	data.Add("client_secret", clientSecret)

	byt, err := apiCall("POST", apiToken, data, "")
	fmt.Println("ret:", string(byt), " err:", err)

	res := newTokenResponse(byt)
	fmt.Println("result:", res)
	token = res.AccessToken
	w.Write(byt)



	aaa := url.Values{}
	// aaa.Add("message", "hello, welcome to test!: token:"+token)
	
	aaa.Add("message", "https://hankwutest-linenotify.herokuapp.com/connect?token="+token)

	cccc, dddd := apiCall("POST", apiNotify, aaa, token)
	fmt.Println("ret:", string(cccc), " err:", dddd)
}
func authHandler(w http.ResponseWriter, r *http.Request) {
	check := func(err error) {
		if err != nil {
			log.Fatal(err)
		}
	}
	t, err := template.New("webpage").Parse(authTmpl)
	check(err)
	noItems := struct {
		ClientID    string
		CallbackURL string
	}{
		ClientID:    clientID,
		CallbackURL: callbackURL,
	}

	err = t.Execute(w, noItems)
	check(err)
}


func connectHandler(w http.ResponseWriter, r *http.Request) {
	check := func(err error) {
		if err != nil {
			log.Fatal(err)
		}
	}

	r.ParseForm() // Populates request.Form
	token := r.Form.Get("token")
	fmt.Printf("connect Get token=%s\n", token)


	t, err := template.New("webpage").Parse(formTmpl)
	check(err)
	noItems := struct {
		TOKEN    string
	}{
		TOKEN:    token,
	}

	err = t.Execute(w, noItems)
	check(err)
}



func connectingHandler(w http.ResponseWriter, r *http.Request) {
	check := func(err error) {
		if err != nil {
			log.Fatal(err)
		}
	}

	r.ParseForm() // Populates request.Form
	token := r.Form.Get("token")
	code := r.Form.Get("code")

	fmt.Printf("connecting token=%s, code=%s\n", token, code)


	t, err := template.New("webpage").Parse(connectedTmpl)
	check(err)
	noItems := struct {
		TOKEN    string
		CODE	 string
	}{
		TOKEN:    token,
		CODE: code,
	}

	err = t.Execute(w, noItems)
	check(err)


	aaa := url.Values{}
	// aaa.Add("message", "hello, welcome to test!: token:"+token)
	aaa.Add("message", "您註冊了您的服務碼為:"+code+", 感謝您")

	cccc, dddd := apiCall("POST", apiNotify, aaa, token)
	fmt.Println("ret:", string(cccc), " err:", dddd)
}