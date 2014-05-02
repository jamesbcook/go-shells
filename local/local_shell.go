package main

import (
        "bufio"
        "fmt"
        "os"
        "os/exec"
        "strings"
        "sync"
)

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
        if head == "cd" {
                os.Chdir(parts[0])
        } else {
                out, err := exec.Command(head, parts...).Output()
                if err != nil {
                        fmt.Printf("%s\n", err)
                }
                fmt.Printf("%s", out)
                wg.Done()
        }
}
