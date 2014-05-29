package main

import (
        "bufio"
        "code.google.com/p/go.crypto/nacl/box"
        "crypto/rand"
        "encoding/base64"
        "flag"
        "fmt"
        "log"
        "net"
        "os"
        "strings"
)

type KeyHolder struct {
        publicKey  *[32]byte
        privateKey *[32]byte
}

func main() {
        port := flag.String("port", "4444", "Port to listen on")
        flag.Parse()
        serv, err := net.Listen("tcp", ":"+*port)
        if err != nil {
                log.Fatal(err)
        }
        defer serv.Close()
        fmt.Println("Listening on port:", *port)
        peerKeys := &KeyHolder{}
        publicKey, privateKey, _ := setup()
        myKeys := &KeyHolder{publicKey, privateKey}
        var nonce [24]byte
        rand.Reader.Read(nonce[:])
        conn, err := serv.Accept()
        if err != nil {
                fmt.Println("Error: ", err)
        }
        defer conn.Close()
        keyExchange(myKeys, peerKeys, nonce, conn)
        prompt := "shell> "
        for {
                fmt.Printf("%s", prompt)
                reader := bufio.NewReader(os.Stdin)
                input, _ := reader.ReadString('\n')
                input = strings.Replace(input, "\n", "", -1)
                enc_input := encrypt(
                        *peerKeys.publicKey,
                        *myKeys.privateKey,
                        nonce,
                        []byte(input))
                if input == "" {
                        continue
                } else if input == "exit" {
                        conn.Write([]byte(enc_input + "\n"))
                        os.Exit(0)
                }
                conn.Write([]byte(enc_input + "\n"))
                recv(*myKeys.privateKey, *peerKeys.publicKey, conn, nonce)
                println("")
        }
}

func recv(mKey, pKey [32]byte, conn net.Conn, nonce [24]byte) {
        message, err := bufio.NewReader(conn).ReadString('\n')
        a := []byte(message)
        i := len(a) - 1
        a = a[:i+copy(a[i:], a[i+1:])]
        if err != nil {
                fmt.Println("Error: ", err)
        }
        dec_mss, ok := decrypt(pKey, mKey, string(a), nonce)
        if ok != true {
                fmt.Println("Decrypt failed")
        }
        fmt.Printf("%s", dec_mss)
}

func setup() (private, public *[32]byte, err error) {
        return box.GenerateKey(rand.Reader)
}

func keyExchange(myKeys, peerKeys *KeyHolder, nonce [24]byte, conn net.Conn) {
        fmt.Println("Exchanging Keys")
        // Send Keys
        var myPub []byte
        for x := range *myKeys.publicKey {
                myPub = append(myPub, myKeys.publicKey[x])
        }
        conn.Write(myPub)
        // Recieve Keys
        reply := make([]byte, 32)
        _, err := conn.Read(reply)
        if err != nil {
                log.Fatal(err)
        }
        var convert [32]byte
        copy(convert[:], reply)
        peerKeys.publicKey = &convert
        // Send Nonce
        var no []byte
        for x := range nonce {
                no = append(no, nonce[x])
        }
        conn.Write(no)
        fmt.Println("Done")
}

func encrypt(publicKey, privateKey [32]byte, nonce [24]byte, message []byte) string {
        enc := box.Seal(nil, message, &nonce, &publicKey, &privateKey)
        return base64.StdEncoding.EncodeToString(enc)
}

func decrypt(publicKey, privateKey [32]byte, message string, nonce [24]byte) (open []byte, status bool) {
        dec, _ := base64.StdEncoding.DecodeString(message)
        return box.Open(nil, dec, &nonce, &publicKey, &privateKey)
}
