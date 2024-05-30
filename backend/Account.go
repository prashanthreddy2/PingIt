package main

import (
	"encoding/json"
	"net/http"
)

func setupResponse(w *http.ResponseWriter) {
	(*w).Header().Set("content-type", "application/json")
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	(*w).Header().Set("Access-Control-Allow-Headers", "Content-Type, ACCEPT, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, X-Custom-Header")
}

func writeResponse(w http.ResponseWriter, r *http.Request, code int, data, errResp interface{}) {
	w.WriteHeader(code)
	if errResp != nil && errResp != "" {
		logger.Errorf("%s %s %v", r.Method, r.RequestURI, errResp)
		data = errResp
	} else {
		logger.Infof("%s %s", r.Method, r.RequestURI)
	}
	body, err := json.Marshal(data)
	if err != nil {
		logger.Error("Internal Server Error", http.StatusInternalServerError)
		return
	}
	w.Write(body)
}

type ResponseStruct struct {
	Data    interface{} `json:"data"`
	Message string      `json:"message"`
	Status  string      `json:"status"`
}

func handleAccountCreation(w http.ResponseWriter, r *http.Request) {
	setupResponse(&w)
	var Response ResponseStruct
	defer r.Body.Close()
	querParams := r.URL.Query()
	userName := querParams["username"]
	password := querParams["password"]
	if len(userName) == 0 || len(password) == 0 {
		Response.Status = "Failure"
		Response.Message = "User params missing"
		logger.Error("Username or password missing")
		writeResponse(w, r, http.StatusBadGateway, Response, "Username or password missing")
	}
	_, err := DbClient.Exec(AccountCreationQuery, userName[0], password[0])
	if err != nil {
		logger.Errorf("Error adding data into db: [%s]", err)
		Response.Status = "Failure"
		Response.Message = "Internal Server Error"
		writeResponse(w, r, http.StatusInternalServerError, Response, err)
		return
	}

	logger.Infof("Sucessfully created acount and added into db")
	Response.Status = "Sucesss"
	Response.Message = "Sucessfully created account"
	writeResponse(w, r, http.StatusOK, Response, nil)
}

func handleLogin(w http.ResponseWriter, r *http.Request) {
	setupResponse(&w)
	var Response ResponseStruct
	defer r.Body.Close()
	queryParams := r.URL.Query()
	username := queryParams["username"]
	password := queryParams["password"]
	if len(username) == 0 || len(password) == 0 {
		Response.Status = "Failure"
		Response.Message = "User params missing"
		logger.Error("Username or password missing")
		writeResponse(w, r, http.StatusBadGateway, nil, Response)
		return
	}
	var (
		count int
	)
	err := DbClient.QueryRow(AccountLoginQuery, username[0], password[0]).Scan(&count)
	if err != nil {
		logger.Errorf("Error trying authenticating from db: [%s]", err)
		Response.Status = "Failure"
		Response.Message = "Error authenticating"
		writeResponse(w, r, http.StatusBadGateway, nil, Response)
		return
	}
	if count == 0 {
		Response.Status = "Failure"
		Response.Message = "wrong username or password"
		logger.Error("wrong username or password")
		writeResponse(w, r, http.StatusBadGateway, nil, Response)
		return
	}
	Response.Status = "Sucess"
	Response.Message = "Auth successful"
	Response.Data = struct {
		Sender string `json:"sender"`
	}{Sender: username[0]}
	logger.Infof("sucess responser: [%+v]", Response)
	writeResponse(w, r, http.StatusOK, Response, nil)
}

func handleSignup(w http.ResponseWriter, r *http.Request) {
	setupResponse(&w)
	var (
		Response      ResponseStruct
		userNameCount int
	)
	queryParams := r.URL.Query()
	username := queryParams["username"]
	password := queryParams["password"]
	if len(username) == 0 || len(password) == 0 || username[0] == "" || password[0] == "" {
		Response.Status = "Failure"
		Response.Message = "User params missing"
		logger.Errorf("Username or password missing")
		writeResponse(w, r, http.StatusBadGateway, nil, Response)
		return
	}
	err := DbClient.QueryRow(GetUsernameCount, username[0]).Scan(&userNameCount)
	if err != nil || userNameCount > 0 {
		Response.Status = "Failure"
		if userNameCount > 0 {
			Response.Message = "Username already taken"
			logger.Infof("Username already exists to create a new account")
		} else {
			logger.Errorf("Error getting the username count: [%s]", err)
			Response.Message = "Error creating account"
		}
		writeResponse(w, r, http.StatusInternalServerError, nil, Response)
		return
	}
	_, err = DbClient.Exec(AccountCreationQuery, username[0], password[0])
	if err != nil {
		Response.Status = "Failure"
		Response.Message = "Error creating account"
		writeResponse(w, r, http.StatusInternalServerError, nil, Response)
		logger.Errorf("Error creating account: [%s]", err)
		return
	}
	Response.Status = "Success"
	Response.Message = "Sucessfully created the account"
	writeResponse(w, r, http.StatusOK, Response, nil)
}
