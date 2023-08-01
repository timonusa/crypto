package main

import (
	"bufio"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"strings"
)

const difficulty = 4  // difficult for PoW
const word = "naruto" // base word

func main() {
	//create a Listener
	listener, err := net.Listen("tcp", "127.0.0.1:8080")
	if err != nil {
		fmt.Println("listening error:", err.Error())
		return
	}
	defer listener.Close()

	fmt.Println("server is working...")

	//endless loop for handling connections
	for {

		//try to create one connection
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Ошибка при принятии подключения:", err.Error())
			return
		}

		//if its ok go for hadling the connection in own thread
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	fmt.Println("Клиент подключен:", conn.RemoteAddr())

	//create objects for reading and writing
	reader := bufio.NewReader(conn)
	writer := bufio.NewWriter(conn)

	//write welcome message with w - object
	writer.WriteString("Connection created \n")
	writer.Flush()

	//endless loop for getting and sending messages with a client
	for {

		//send message for client
		writer.WriteString("Enter a code\n") // не выводит в клиенте а значи сервер криво отдает
		writer.Flush()

		//read answer from client
		input, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Ошибка при чтении данных:", err.Error())
			return
		}
		data := strings.TrimSpace(input)
		fmt.Println(data)

		//figuring out if he had done work
		isRight := nonceIsRight("naruto", 4, data)

		//if hashes are not the same
		if !isRight {
			writer.WriteString(fmt.Sprintf("Проверка не пройдена\n"))
			writer.Flush()
			continue
		}

		//tell him he is a good boy
		writer.WriteString(fmt.Sprintf("Проверка пройдена !!!!!\n"))
		writer.Flush()

		//show the truth
		writer.WriteString(fmt.Sprintf("get the wisdom....rand\n"))
		writer.WriteString(getQuote() + "\n")
		writer.Flush() // Flush the writer to send the response immediately
	}
}

// for checking the PoW
func nonceIsRight(word string, difficulty int, nonce string) bool {

	hash := sha256.Sum256([]byte(fmt.Sprintf("%s%s", word, nonce)))
	hashStr := hex.EncodeToString(hash[:])

	if strings.HasPrefix(hashStr, strings.Repeat("0", difficulty)) {
		return true
	}

	return false
}

func getQuote() string {
	url := "https://zenquotes.io/?api=quotes"
	response, err := http.Get(url)
	if err != nil {
		fmt.Println("error qoutes reading:", err.Error())
		return ""
	}
	defer response.Body.Close()
	data, _ := io.ReadAll(response.Body)
	var quotes []interface{}
	err = json.Unmarshal(data, &quotes)
	if err != nil {
		fmt.Println("error quotes parsing:", err.Error())
		return ""
	}
	firstQuote := quotes[0].(map[string]interface{})
	return firstQuote["q"].(string)
}
