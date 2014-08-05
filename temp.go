package main

import "fmt"
import "github.com/conformal/btcec"
import "code.google.com/p/go.crypto/ripemd160"

func GenKeys() (string, string) {

}

func main() {
    privateKey, publicKey := GenKey();
    fmt.Println(privateKey,"<=>",publicKey);
}
