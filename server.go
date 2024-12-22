package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

type Response struct {
	Body float64 `json:"result"`
}

type Error struct {
	Body string `json:"error"`
}

func MainHandler(w http.ResponseWriter, r *http.Request) {
	expression := r.URL.Query().Get("expression")
	if expression == "shutdown" {
		ShutdownLogger()
		os.Exit(0)
	}
	result, err := Calc(expression)
	if err.Type != Ok {
		panic(err)
	}
	Log(fmt.Sprintf("request completed, the answer is %f", result))
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(Response{
		Body: result})
}

func PanicMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				myerr := err.(CustomPanic)
				var (
					status int
					msg    string
				)

				if myerr.Type == InvalidExpression {
					status = http.StatusUnprocessableEntity
					msg = "Expression is not valid"
				} else if myerr.Type == InnerExpression {
					status = http.StatusInternalServerError
					msg = "Internal server error"
				}
				w.WriteHeader(status)
				Log(myerr.Text.Error())
				// json.NewEncoder(w).Encode(Error{
				// 	Body: myerr.Text.Error()})
				json.NewEncoder(w).Encode(Error{
					Body: msg,
				})
			}
		}()
		w.Header().Set("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})

}

func StartServer() {
	http.HandleFunc("/apl/v1/calculate", PanicMiddleware(MainHandler))
	Log("starting server...")
	http.ListenAndServe(":8080", nil)
}
