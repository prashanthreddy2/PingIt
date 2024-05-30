package main

import (
	"database/sql"
	"fmt"
	"net/http"
)

func handleUserList(w http.ResponseWriter, r *http.Request) {
	setupResponse(&w)
	var Resp ResponseStruct
	rows, err := DbClient.Query(GetUserList)
	if err != nil {
		if err == sql.ErrNoRows {
			Resp.Status = "Success"
			Resp.Data = []string{}
			writeResponse(w, r, http.StatusOK, Resp, nil)
			return
		}
		logger.Errorf("Error getting users list: [%s]", err)
		Resp.Status = "Failure"
		Resp.Message = fmt.Sprintf("Error getting user list: [%s]", err)
		writeResponse(w, r, http.StatusInternalServerError, nil, Resp)
	}
	var userList []string
	user := r.URL.Query()["user"]
	for rows.Next() {
		var username string
		err := rows.Scan(&username)
		if err != nil {
			logger.Errorf("Error scanning user list: [%s]", err)
		} else {
			if username == user[0] {
				continue
			}
			userList = append(userList, username)
		}
	}
	Resp.Status = "Success"
	Resp.Data = userList
	writeResponse(w, r, http.StatusOK, Resp, nil)
}

func handleFetchMessages(w http.ResponseWriter, r *http.Request) {
	setupResponse(&w)
	var Resp ResponseStruct
	sender := r.URL.Query()["sender"]
	reciever := r.URL.Query()["reciever"]
	if len(sender) == 0 || len(reciever) == 0 || sender[0] == "" || reciever[0] == "" {
		Resp.Message = fmt.Sprint("Missing sender or reciever Id")
		Resp.Status = "Failure"
		writeResponse(w, r, http.StatusBadRequest, nil, Resp)
		return
	}
	type MessageStruct struct {
		Sender    string `json:"sender"`
		Reciever  string `json:"reciever"`
		Message   string `json:"message"`
		Timestamp string `json:"timestamp"`
	}
	Messages := []MessageStruct{}

	rows, err := DbClient.Query(GetMessagesQuery, sender[0], reciever[0])
	if err != nil {
		logger.Errorf("Inside no rows")
		if err == sql.ErrNoRows {
			Resp.Status = "Success"
			logger.Infof("Err No rows")
			Resp.Data = Messages
			writeResponse(w, r, http.StatusOK, Resp, nil)
			return
		} else {
			logger.Errorf("Error querying the db: [%s]", err)
			Resp.Message = fmt.Sprintf("Error querying the db: [%s]", err)
			Resp.Status = "Failure"
			writeResponse(w, r, http.StatusBadRequest, nil, Resp)
			return
		}
	}
	for rows.Next() {
		var sender, reciever, timestamp, message string
		err := rows.Scan(&sender, &reciever, &message, &timestamp)
		if err != nil {
			logger.Errorf("Error scanning values: [%s]", err)
			continue
		}
		Messages = append(Messages, MessageStruct{
			Sender:    sender,
			Reciever:  reciever,
			Timestamp: timestamp,
			Message:   message,
		})
	}
	Resp.Status = "Success"
	Resp.Data = Messages
	writeResponse(w, r, http.StatusOK, Resp, nil)
}

func handleFetchRelations(w http.ResponseWriter, r *http.Request) {
	setupResponse(&w)
	var Resp ResponseStruct
	sender := r.URL.Query()["sender"]
	if len(sender) == 0 || sender[0] == "" {
		Resp.Message = fmt.Sprint("Missing sender params")
		Resp.Status = "Failure"
		writeResponse(w, r, http.StatusBadRequest, nil, Resp)
		return
	}
	type RelationStruct struct {
		RoomName string `json:"roomName"`
	}
	relationsList := []string{}
	rows, err := DbClient.Query(`select distinct receiver from relationdb where sender=$1 union select distinct sender from relationdb where receiver=$1`, sender[0])
	if err != nil {
		logger.Errorf("Error getting relations for user %s: [%s]", sender[0], err)
		Resp.Message = fmt.Sprint("Unable to fetch the relations")
		Resp.Status = "Failure"
		writeResponse(w, r, http.StatusInternalServerError, nil, Resp)
		return
	}
	for rows.Next() {
		var reciever string
		err := rows.Scan(&reciever)
		if err != nil {
			logger.Errorf("Error scanning the relation data: [%s]", err)
			continue
		}
		relationsList = append(relationsList, reciever)
	}
	Resp.Status = "Success"
	Resp.Data = relationsList
	writeResponse(w, r, http.StatusOK, Resp, nil)
}
