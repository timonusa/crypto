package main

import (
	"bufio"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net"
	"os"
	"strings"
)

func main() {

	//create connection to a server
	//conn, err := net.Dial("tcp", "127.0.0.1:8080")
	conn, err := net.Dial("tcp", "server:8080")
	if err != nil {
		fmt.Println("connection error:", err.Error())
		return
	}
	defer conn.Close()

	//get the response of the server
	go readServerResponses(conn)

	//creating a reader from console and writer for connection
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(conn)

	//read data from console
	for {
		input, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Error reading the input:", err.Error())
			return
		}

		//write message to a server
		writer.WriteString(input)
		writer.Flush()

	}

}

// reading response from server
func readServerResponses(conn net.Conn) {
	reader := bufio.NewReader(conn)

	for {
		response, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("error getting the response:", err.Error())
			return
		}

		fmt.Println(response)
	}
}

// for Proving the Work
func calculatePoW(data string, difficulty int) int {
	nonce := 0
	for {
		hash := sha256.Sum256([]byte(fmt.Sprintf("%s%d", data, nonce)))
		hashStr := hex.EncodeToString(hash[:])

		if strings.HasPrefix(hashStr, strings.Repeat("0", difficulty)) {
			return nonce
		}

		nonce++
	}
}
