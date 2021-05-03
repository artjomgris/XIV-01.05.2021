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
	"time"
)

type Client struct {
	Id      int    `json:"id"`
	Name    string `json:"name"`
	Message struct {
		Body   string `json:"body"`
		AddrId int    `json:"addr_id"`
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
		if len(text) >= 5 {
			if text[0:5] == "/info" {
				/*var users []Client
				users = getUsers()
				fmt.Println(string(users))*/
				fmt.Println(111)
			} else {
				i := strings.Index(text, " ")
				if i == -1 {
					fmt.Println("Please enter your message as [ID MESSAGE] OR /info")
				} else {
					client.Message.AddrId, err = strconv.Atoi(text[0:i])
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
						fmt.Println("Please enter your message as [ID MESSAGE] OR /info")
					}
				}

			}
		} else {
			fmt.Println("Please enter your message as [ID MESSAGE] OR /info")
		}

	}
}

func readMessage(c *websocket.Conn) {
	for {
		log.Println("READ")
		err := c.WriteMessage(websocket.TextMessage, []byte("getmsg"+strconv.Itoa(client.Id)))
		_, message, err := c.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			return
		}
		log.Printf("Message: %s", message)
		time.Sleep(1 * time.Second)
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

	req, err := http.NewRequest("POST", "http://127.0.0.1:8080/register", &body)
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
