package server

import (
	"encoding/json"
	"fmt"
	"money/agent/dao"
	"net/http"
)

func setCOROS(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Add("Access-Control-Allow-Headers", "Content-Type,AccessToken,X-CSRF-Token, Authorization, Token")
	w.Header().Add("Access-Control-Allow-Credentials", "true")
	w.Header().Add("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	w.Header().Set("content-type", "application/json;charset=UTF-8")
	w.WriteHeader(http.StatusOK)
}

func BscLatestTransaction(w http.ResponseWriter, r *http.Request) {
	setCOROS(w)

	data, _ := dao.GetBscLatestTransactionFromDB()
	ret, _ := json.Marshal(data)

	fmt.Fprintf(w, string(ret))
}

func BscHotTransaction(w http.ResponseWriter, r *http.Request) {
	setCOROS(w)

	data, _ := dao.GetBscHotTransactionFromDB()
	ret, _ := json.Marshal(data)

	fmt.Fprintf(w, string(ret))
}

func BscNewCoin(w http.ResponseWriter, r *http.Request) {
	setCOROS(w)

	data, _ := dao.GetNewBscCoin()
	ret, _ := json.Marshal(data)

	fmt.Fprintf(w, string(ret))
}

func Serve() {
	http.HandleFunc("/bsc/latest", BscLatestTransaction)
	http.HandleFunc("/bsc/hot", BscHotTransaction)
	http.HandleFunc("/bsc/new", BscNewCoin)
	http.ListenAndServe(":8080", nil)
}
