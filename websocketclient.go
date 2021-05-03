
package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
)

type Client struct {
	Id int `json:"id"`
	Name string `json:"name"`
	Message struct {
		Body string `json:"body"`
		AddrId int `json:"addr_id"`
	}
}
var client Client

func main() {
	c, _, err := websocket.DefaultDialer.Dial("ws://localhost:8080/socket", nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	defer c.Close()
	fmt.Println("Please enter your name:")
	fmt.Scan(&client.Name)
	client.Id = register(client.Name)
	fmt.Println(client.Id)

	fmt.Println("Type /info to get information about users aviable!")
	go readMessage(c)
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Println("Type ID Message")
		text, _ := reader.ReadString('\n')

		if text[0:5] == "/info" {
			/*var users []Client
			users := getUsers()
			fmt.Println(string(users))*/
			fmt.Println(111)
		} else {
			i := strings.Index(text, " ")
			client.Message.AddrId, err = strconv.Atoi(text[0:i+1])
			if err == nil {
				client.Message.Body = text[i:]
				msg, err := json.Marshal(client)
				if err != nil {
					panic(err)
				}
				err = c.WriteMessage(websocket.TextMessage, msg)
				if err != nil {
					panic(err)
				}
			} else {
				fmt.Println("Please enter your message as [ID MESSAGE]")
			}

		}
	}
}

func readMessage(c *websocket.Conn) {
	for {
		_, message, err := c.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			return
		}
		log.Printf("Message: %s", message)
	}
}

func getUsers() []byte {
	client := &http.Client{}
	req, err := http.NewRequest("GET", "http://127.0.0.1:8080/users", nil)
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}

	rec, err := io.ReadAll(resp.Body)
	defer resp.Body.Close()
	return rec
}

func register(name string) int {
	client := &http.Client{}
	var body bytes.Buffer
	body.Write([]byte(name))

	req, err := http.NewRequest("GET", "http://127.0.0.1:8080/register", &body)
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}

	rec, err := io.ReadAll(resp.Body)
	defer resp.Body.Close()
	ret, err := strconv.Atoi(string(rec))
	if err != nil {
		panic(err)
	}
	return ret
}