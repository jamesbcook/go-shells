package main

import (
        "bufio"
        "flag"
        "fmt"
        "log"
        "net"
        "os"
        "strings"
)

func main() {
        port := flag.String("port", "4444", "Port to listen on")
        flag.Parse()
        serv, err := net.Listen("tcp", ":"+*port)
        if err != nil {
                log.Fatal(err)
        }
        defer serv.Close()
        fmt.Println("Listening on port: ", *port)
        conn, err := serv.Accept()
        if err != nil {
                fmt.Println("Error: ", err)
        }
        defer conn.Close()
        prompt := "shell> "
        for {
                fmt.Printf("%s", prompt)
                reader := bufio.NewReader(os.Stdin)
                input, _ := reader.ReadString('\n')
                input = strings.Replace(input, "\n", "", -1)
                if input == "" {
                        fmt.Printf("%s", prompt)
                } else if input == "exit" {
                        conn.Write([]byte("exit\n"))
                        os.Exit(0)
                }
                conn.Write([]byte(input + "\n"))
                recv(conn)
                println("")
        }
}

func recv(conn net.Conn) {
        reply := make([]byte, 1024)
        length, err := conn.Read(reply)
        if err != nil {
                println(err)
                os.Exit(1)
        }
        fmt.Print(string(reply))
        if length == 1024 {
                recv(conn)
        }
}
