package main

import (
	"io/ioutil"
	"log"
	"net/http"
)

type LcRequest struct {
	*http.Request
}
type LcResponse struct {
	http.ResponseWriter
}
type LcHandlerFunc func(*LcRequest, *LcResponse)

func (r *LcRequest) GetParam(name string) interface{} {
	body, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		log.Println("GetParams error:", err)
		return nil
	}
	p := FromJsonBytes[map[string]interface{}](body)
	return p[name]
}
func (r *LcRequest) GetParams() (map[string]interface{}, error) {
	body, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		log.Println("GetParams error:", err)
		return make(map[string]interface{}), err
	}
	return FromJsonBytes[map[string]interface{}](body), nil
}
func (w *LcResponse) WriteJSON(data interface{}) {
	json := ToJsonBytes(data)
	w.WriteFile(json, "application/json")
}
func (w *LcResponse) WriteHTML(data []byte) {
	w.WriteFile(data, "text/html")
}
func (w *LcResponse) WriteFile(data []byte, contentType string) {
	w.Header().Set("Content-Type", contentType)
	w.Write(data)
}
func (w *LcResponse) WriteJS(data []byte) {
	w.WriteFile(data, "application/javascript")
}

func handleFuncWithMethod(method string, pattern string, handler LcHandlerFunc) {
	http.HandleFunc(pattern, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != method {
			w.WriteHeader(http.StatusMethodNotAllowed)
		} else {
			handler(&LcRequest{r}, &LcResponse{w})
		}
	})
}

func Get(pattern string, handler LcHandlerFunc) {
	handleFuncWithMethod(http.MethodGet, pattern, handler)
}
func Post(pattern string, handler LcHandlerFunc) {
	handleFuncWithMethod(http.MethodPost, pattern, handler)
}
func Put(pattern string, handler LcHandlerFunc) {
	handleFuncWithMethod(http.MethodPut, pattern, handler)
}
func Delete(pattern string, handler LcHandlerFunc) {
	handleFuncWithMethod(http.MethodDelete, pattern, handler)
}
func Patch(pattern string, handler LcHandlerFunc) {
	handleFuncWithMethod(http.MethodPatch, pattern, handler)
}
