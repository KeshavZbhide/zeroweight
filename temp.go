package main

import "fmt"
import "github.com/conformal/btcec"
import "crypto/ecdsa"
import "crypto/rand"
import "crypto/sha256"
import "code.google.com/p/go.crypto/ripemd160"

func publicKeyStructToPublicKeyBytes(key *ecdsa.PublicKey) []byte {
    xylen := len(key.X.Bytes());
    keyBytes := make([]byte, 1+(2*xylen));
    keyBytes[0] = byte(4);
    copied := copy(keyBytes[1:], key.X.Bytes());
    copy(keyBytes[(1+copied):], key.Y.Bytes());
    return keyBytes;
}

func publicKeyHash(key []byte) []byte {
    s256 := sha256.Sum256(key);
    ripemd160Hash := ripemd160.New();
    ripemd160Hash.Write(s256[:]);
    return ripemd160Hash.Sum(nil);
}


func GenKeys() (string, string) {
    privateKeyStrut, err := ecdsa.GenerateKey(btcec.S256(), rand.Reader);
    if err != nil {
        fmt.Println(err);
        return "", "";
    }
    privateKeyWif := base58CheckEncodeKey(0x80, privateKeyStrut.D.Bytes());
    publicKeyBytes := publicKeyStructToPublicKeyBytes(&(privateKeyStrut.PublicKey));
    hash0 := sha256.Sum256(publicKeyBytes);
    ripemd160 := ripemd160.New();
    ripemd160.Write(hash0[:]);
    keyHash := ripemd160.Sum(nil);
    publicKeyWif := base58CheckEncodeKey(0,keyHash);
    return privateKeyWif, publicKeyWif;
}

func main() {
    privateKey, publicKey := GenKeys();
    fmt.Println(privateKey,"<=>",publicKey);
}
