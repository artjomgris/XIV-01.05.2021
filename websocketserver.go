package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"io"
	"log"
	"net/http"
	"strconv"
)

type Sclient struct {
	Id       int    `json:"id"`
	Ip       string `json:"ip"`
	Name     string `json:"name"`
	Messages []Client
}

type Client struct {
	Id      int    `json:"id"`
	Name    string `json:"name"`
	Message struct {
		Body   string `json:"body"`
		AddrId int    `json:"addr_id"`
	}
}

var sclients []Sclient
var upgrader = websocket.Upgrader{}

func Handler(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "content-type")

	if req.Method == "POST" {
		data, err := io.ReadAll(req.Body)
		req.Body.Close()
		if err != nil {
			return
		}

		log.Printf("%s\n", data)
		io.WriteString(w, "successful post")
	} else if req.Method == "OPTIONS" {
		w.WriteHeader(204)
	} else {
		w.WriteHeader(405)
	}

}

func Registration(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "content-type")

	if req.Method == "POST" {
		data, err := io.ReadAll(req.Body)
		req.Body.Close()
		if err != nil {
			return
		}
		strdata := string(data)
		id := len(sclients) + 1
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
	messageType, p, err := conn.ReadMessage()
	if err != nil {
		panic(err)
		return
	}
	log.Println(string(p)[0:6])
	if string(p)[0:6] == "getmsg" {
		msg := string(p)
		var id int
		id, err = strconv.Atoi(msg[6:])
		fmt.Println(id)
		if err != nil {
			panic(err)
		}
		for _, j := range sclients {
			if j.Messages == nil {
				if err := conn.WriteMessage(messageType, nil); err != nil {
					panic(err)
				}
			} else {
				if j.Id == id {
					fmt.Println(j.Messages)
					jsmsg, err := json.Marshal(j.Messages)
					if err != nil {
						panic(err)
					}
					if err := conn.WriteMessage(messageType, jsmsg); err != nil {
						panic(err)
					}
					j.Messages = nil
				}
			}

		}
	} else {
		log.Println(string(p))
		var message Client
		err = json.Unmarshal(p, &message)
		id := message.Message.AddrId
		for _, j := range sclients {
			if j.Id == id {
				j.Messages = append(j.Messages, message)
			}
		}
	}

}

func main() {
	http.HandleFunc("/", Handler)
	http.HandleFunc("/socket", Socket)
	http.HandleFunc("/register", Registration)

	err := http.ListenAndServe(":8080", nil)
	panic(err)
}
