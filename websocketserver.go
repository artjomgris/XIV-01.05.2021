package main

import (
	"fmt"
	"net/http"
	"io"
	"log"
	"github.com/gorilla/websocket"
)

type Sclient struct {
	Id int `json:"id"`
	Ip string `json:"ip"`
	Name string `json:"name"`
}
var sclients []Sclient

var upgrader = websocket.Upgrader{}

func Handler(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin","*")
	w.Header().Set("Access-Control-Allow-Methods","POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "content-type")
	
	if req.Method == "POST" {
		data, err := io.ReadAll(req.Body)
		req.Body.Close()
		if err != nil {return }
		
		log.Printf("%s\n", data)
		io.WriteString(w, "successful post")
	} else if req.Method == "OPTIONS" {
		w.WriteHeader(204)
	} else {
		w.WriteHeader(405)
	}
	
}

func Registration(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin","*")
	w.Header().Set("Access-Control-Allow-Methods","POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "content-type")

	if req.Method == "POST" {
		data, err := io.ReadAll(req.Body)
		req.Body.Close()
		if err != nil {return }
		strdata := string(data)
		for _, j := range sclients {
			if j.Name == strdata {
				io.WriteString(w, "666")
				break
			}
		}
		id := len(sclients)+1
		sclients = append(sclients, Sclient{
			Id:   id,
			Ip:   req.RemoteAddr,
			Name: strdata,
		})
		io.WriteString(w, fmt.Sprint(id))
		fmt.Println(id)
	} else if req.Method == "OPTIONS" {
		w.WriteHeader(204)
	} else {
		w.WriteHeader(405)
	}

}

func Socket(w http.ResponseWriter, req *http.Request) {
	conn, err := upgrader.Upgrade(w, req, nil)
	if err != nil {
		panic(err)
		return
	}
	for {
		messageType, p, err := conn.ReadMessage()
		if err != nil {
			panic(err)
			return
		}
		if err := conn.WriteMessage(messageType, p); err != nil {
			panic(err)
			return
		}
		log.Println(string(p))
	}
}

func main() {
	http.HandleFunc("/", Handler)
	http.HandleFunc("/socket", Socket)
	http.HandleFunc("/register", Registration)
	
	err := http.ListenAndServe(":8080", nil)
	panic(err)
}

