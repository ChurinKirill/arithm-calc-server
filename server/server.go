package server

import (
	calculator "arithm-calc-server/calculation"
	logger "arithm-calc-server/logger"
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
		logger.ShutdownLogger()
		os.Exit(0)
	}
	result, err := calculator.Calc(expression)
	if err.Type != calculator.Ok {
		panic(err)
	}
	logger.Log(fmt.Sprintf("request completed, the answer is %f", result))
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(Response{
		Body: result})
}

func PanicMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				myerr := err.(calculator.CustomPanic)
				var (
					status int
					msg    string
				)

				if myerr.Type == calculator.InvalidExpression {
					status = http.StatusUnprocessableEntity
					msg = "Expression is not valid"
				} else if myerr.Type == calculator.InnerExpression {
					status = http.StatusInternalServerError
					msg = "Internal server error"
				}
				w.WriteHeader(status)
				logger.Log(myerr.Text.Error())
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
	logger.Log("starting server...")
	http.ListenAndServe(":8080", nil)
}
