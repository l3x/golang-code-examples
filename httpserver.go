package main

import (
	"fmt"
	"log"
	"errors"
	"net/http"
	"io/ioutil"
	"html/template"
	"regexp"
	"encoding/json"
	"github.com/l3x/jsoncfgo"
	"github.com/go-goodies/go_utils"
)

var Dir string
var Users jsoncfgo.Obj
var AppContext *go_utils.Singleton

func HtmlFileHandler(response http.ResponseWriter, request *http.Request, filename string){
	response.Header().Set("Content-type", "text/html")
	webpage, err := ioutil.ReadFile(Dir + filename)  // read whole the file
	if err != nil {
		http.Error(response, fmt.Sprintf("%s file error %v", filename, err), 500)
	}
	fmt.Fprint(response, string(webpage));
}

func HelpHandler(response http.ResponseWriter, request *http.Request){
	HtmlFileHandler(response, request, "/help.html")
}

func AjaxHandler(response http.ResponseWriter, request *http.Request){
	HtmlFileHandler(response, request, "/ajax.html")
}

func printCookies(response http.ResponseWriter, request *http.Request) {
	cookieNameForUsername := AppContext.Data["CookieNameForUsername"].(string)
	fmt.Println("COOKIES:")
	for _, cookie := range request.Cookies() {
		fmt.Printf("%v: %v\n", cookie.Name, cookie.Value)
		if cookie.Name == cookieNameForUsername {
			SetUsernameCookie(response, cookie.Value)
		}
	}; fmt.Println("")
}

func UserHandler(response http.ResponseWriter, request *http.Request){
	response.Header().Set("Content-type", "application/json")
	// json data to send to client
	data := map[string]string { "api" : "user", "name" : "" }
	userApiURL := regexp.MustCompile(`^/user/(\w+)$`)
	usernameMatches := userApiURL.FindStringSubmatch(request.URL.Path)
	// regex matches example: ["/user/joesample", "joesample"]
	if len(usernameMatches) > 0 {
		printCookies(response, request)
		var userName string
		userName = usernameMatches[1]  // ex: joesample
		userObj := AppContext.Data[userName]
		fmt.Printf("userObj: %v\n", userObj)
		if userObj == nil {
			msg := fmt.Sprintf("Invalid username (%s)", userName)
			panic(errors.New(msg))
		} else {
			// Send JSON to the client
			thisUser := userObj.(jsoncfgo.Obj)
			fmt.Printf("thisUser: %v\n", thisUser)
			data["name"] = thisUser["firstname"].(string) + " " + thisUser["lastname"].(string)
		}
		json_bytes, _ := json.Marshal(data)
		fmt.Printf("json_bytes: %s\n", string(json_bytes[:]))
		fmt.Fprintf(response, "%s\n", json_bytes)

	} else {
		http.Error(response, "404 page not found", 404)
	}
}

func SetUsernameCookie(response http.ResponseWriter, userName string){
	// Add cookie to response
	cookieName := AppContext.Data["CookieNameForUsername"].(string)
	cookie := http.Cookie{Name: cookieName, Value: userName}
	http.SetCookie(response, &cookie)
}

func DebugFormHandler(response http.ResponseWriter, request *http.Request){

	printCookies(response, request)

	err := request.ParseForm()  // Parse URL and POST data into request.Form
	if err != nil {
		http.Error(response, fmt.Sprintf("error parsing url %v", err), 500)
	}

	// Set cookie and MIME type in the HTTP headers.
	fmt.Printf("request.Form: %v\n", request.Form)
	if request.Form["username"] != nil {
		cookieVal := request.Form["username"][0]
		fmt.Printf("cookieVal: %s\n", cookieVal)
		SetUsernameCookie(response, cookieVal)
	}; fmt.Println("")

	templateHandler(response, request)
	response.Header().Set("Content-type", "text/plain")

	// Send debug diagnostics to client
	fmt.Fprintf(response, "<table>")
	fmt.Fprintf(response, "<tr><td><strong>request.Method    </strong></td><td>'%v'</td></tr>", request.Method)
	fmt.Fprintf(response, "<tr><td><strong>request.RequestURI</strong></td><td>'%v'</td></tr>", request.RequestURI)
	fmt.Fprintf(response, "<tr><td><strong>request.URL.Path  </strong></td><td>'%v'</td></tr>", request.URL.Path)
	fmt.Fprintf(response, "<tr><td><strong>request.Form      </strong></td><td>'%v'</td></tr>", request.Form)
	fmt.Fprintf(response, "<tr><td><strong>request.Cookies() </strong></td><td>'%v'</td></tr>", request.Cookies())
	fmt.Fprintf(response, "</table>")

}

func DebugQueryHandler(response http.ResponseWriter, request *http.Request){

	// Set cookie and MIME type in the HTTP headers.
	response.Header().Set("Content-type", "text/plain")

	// Parse URL and POST data into the request.Form
	err := request.ParseForm()
	if err != nil {
		http.Error(response, fmt.Sprintf("error parsing url %v", err), 500)
	}

	// Send debug diagnostics to client
	fmt.Fprintf(response, " request.Method     '%v'\n", request.Method)
	fmt.Fprintf(response, " request.RequestURI '%v'\n", request.RequestURI)
	fmt.Fprintf(response, " request.URL.Path   '%v'\n", request.URL.Path)
	fmt.Fprintf(response, " request.Form       '%v'\n", request.Form)
	fmt.Fprintf(response, " request.Cookies()  '%v'\n", request.Cookies())
}

func templateHandler(w http.ResponseWriter, r *http.Request) {
	t, _ := template.New("form.html").Parse(form)
	t.Execute(w, "")
}

func formHandler(w http.ResponseWriter, r *http.Request) {
	log.Println(r.Form)
	templateHandler(w, r)
}

var form = `
<h1>Debug Info (POST form)</h1>
<form method="POST" action="" name="frmTest">
<div>
    <label for="username">User Name</label>
    <input id="username" name="username" placeholder="joesample, alicesmith, or bobbrown" required="" type="text"
size="50">
</div>
<div><input type="submit" value="Submit"></div>
</form>

</form>
`

func errorHandler(f func(http.ResponseWriter, *http.Request) error) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("errorHandler...")
		err := f(w, r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			log.Printf("handling %q: %v", r.RequestURI, err)
		}
	}
}

func doThis() error { return nil }
func doThat() error { return errors.New("ERROR - doThat") }

func wrappedHandler(w http.ResponseWriter, r *http.Request) error {
	log.Println("betterHandler...")
	if err := doThis(); err != nil {
		return fmt.Errorf("doing this: %v", err)
	}

	if err := doThat(); err != nil {
		return fmt.Errorf("doing that: %v", err)
	}
	return nil
}


func main() {
	cfg := jsoncfgo.Load("/Users/lex/dev/go/data/webserver/webserver-config.json")

	host := cfg.OptionalString("host", "localhost")
	fmt.Printf("host: %v\n", host)

	port := cfg.OptionalInt("port", 8080)
	fmt.Printf("port: %v\n", port)

	Dir = cfg.OptionalString("dir", "www/")
	fmt.Printf("web_dir: %v\n", Dir)

	redirect_code := cfg.OptionalInt("redirect_code", 307)
	fmt.Printf("redirect_code: %v\n\n", redirect_code)

	mux := http.NewServeMux()

	fileServer := http.Dir(Dir)
	fileHandler := http.FileServer(fileServer)
	mux.Handle("/", fileHandler)

	rdh := http.RedirectHandler("http://example.org", redirect_code)
	mux.Handle("/redirect", rdh)
	mux.Handle("/notFound", http.NotFoundHandler())

	mux.Handle("/help", http.HandlerFunc( HelpHandler ))

	mux.Handle("/debugForm", http.HandlerFunc( DebugFormHandler ))
	mux.Handle("/debugQuery", http.HandlerFunc( DebugQueryHandler ))

	mux.Handle("/user/", http.HandlerFunc( UserHandler ))
	mux.Handle("/ajax", http.HandlerFunc( AjaxHandler ))

	mux.Handle("/adapter", errorHandler(wrappedHandler))

	log.Printf("Running on port %d\n", port)

	addr := fmt.Sprintf("%s:%d", host, port)

	Users := jsoncfgo.Load("/Users/lex/dev/go/data/webserver/users.json")

	joesample := Users.OptionalObject("joesample")
	fmt.Printf("joesample: %v\n", joesample)
	fmt.Printf("joesample['firstname']: %v\n", joesample["firstname"])
	fmt.Printf("joesample['lastname']: %v\n\n", joesample["lastname"])

	alicesmith := Users.OptionalObject("alicesmith")
	fmt.Printf("alicesmith: %v\n", alicesmith)
	fmt.Printf("alicesmith['firstname']: %v\n", alicesmith["firstname"])
	fmt.Printf("alicesmith['lastname']: %v\n\n", alicesmith["lastname"])

	bobbrown := Users.OptionalObject("bobbrown")
	fmt.Printf("bobbrown: %v\n", bobbrown)
	fmt.Printf("bobbrown['firstname']: %v\n", bobbrown["firstname"])
	fmt.Printf("bobbrown['lastname']: %v\n\n", bobbrown["lastname"])

	AppContext = go_oops.NewSingleton()
	AppContext.Data["CookieNameForUsername"] = "testapp-username"
	AppContext.Data["joesample"] = joesample
	AppContext.Data["alicesmith"] = alicesmith
	AppContext.Data["bobbrown"] = bobbrown
	fmt.Printf("AppContext: %v\n", AppContext)
	fmt.Printf("AppContext.Data[\"joesample\"]: %v\n", AppContext.Data["joesample"])
	fmt.Printf("AppContext.Data[\"alicesmith\"]: %v\n", AppContext.Data["alicesmith"])
	fmt.Printf("AppContext.Data[\"bobbrown\"]: %v\n\n", AppContext.Data["bobbrown"])

	err := http.ListenAndServe(addr, mux)
	fmt.Println(err.Error())
}
