package main

import (
	"handler"
	"server"
)

func main() {
	// Instantiate server
	dbURL := "admin/adminpassword@database-1.c3yn9dlnua9q.us-east-2.rds.amazonaws.com/ORCL"
	srvAddr := "ec2-3-128-153-42.us-east-2.compute.amazonaws.com:1323"
	// srvAddr := "localhost:1323"
	h := handler.NewHandler(dbURL)
	e := server.NewServer(h)

	// Initiate parallel server control
	go server.TerminalControl(e, h, srvAddr)

	// Start server
	e.Logger.Fatal(e.Start(srvAddr))
}
