package main

import (
        "bufio"
        "fmt"
        "os"
        "os/exec"
        "runtime"
        "strings"
        "sync"
)

const os_version string = runtime.GOOS

func main() {
        shell()
}

func shell() {
        prompt := "shell> "
        var wg sync.WaitGroup
        for {
                fmt.Printf("%s", prompt)
                reader := bufio.NewReader(os.Stdin)
                input, _ := reader.ReadString('\n')
                input = strings.Replace(input, "\n", "", -1)
                if input == "" {
                        continue
                } else if input == "exit" {
                        break
                }
                wg.Add(1)
                exec_cmd(input, &wg)
        }
}

func exec_cmd(cmd string, wg *sync.WaitGroup) {
        parts := strings.Fields(cmd)
        head := parts[0]
        parts = parts[1:len(parts)]
        switch head {
        case "exit":
                os.Exit(0)
        case "cd":
                os.Chdir(parts[0])
                fmt.Println("dir chaged")
        default:
                //out, err := exec.Command(head, parts...).Output()
                if os_version != "windows" {
                        out, err := exec.Command("sh", "-c", cmd).Output()
                        if err != nil {
                                fmt.Println("command not found")
                        }
                        wg.Done()
                        fmt.Println(string(out))
                } else {
                        out, err := exec.Command("cmd.exe", "/c", cmd).Output()
                        if err != nil {
                                fmt.Println("command not found")
                        }
                        wg.Done()
                        fmt.Println(string(out))
                }
        }
}
