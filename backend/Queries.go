package main

var (
	AccountCreationQuery = `insert into acctdb(username, password) values($1,$2)`
	DataEntryQuery       = `insert into datadb(sender,receiver,message,timestamp) values($1,$2,$3,$4)`
	AccountLoginQuery    = `select count(*) from acctdb where username=$1 and password=$2`
	GetUserList          = `select username from acctdb`
	GetMessagesQuery     = `select sender, receiver, message, timestamp from datadb where (sender=$1 and receiver=$2) or (sender=$2 and receiver=$1) order by timestamp asc`
	GetUsernameCount     = `select count(username) from acctdb where username=$1`
	DeleteDuplicates     = `DELETE FROM relationdb d WHERE EXISTS ( SELECT 1 FROM relationdb t WHERE t.ctid < d.ctid AND t.* = d.* );`
)
