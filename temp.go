package main

import "fmt"
import "github.com/conformal/btcec"
//import "crypto/ecdsa"
//import "crypto/rand"
//import "encoding/hex"
import "code.google.com/p/go.crypto/ripemd160";
import "crypto/sha256"
import "os"
//import "bytes"

func main_hold() {
    //ky := base58CheckDecodeKey("5KehCbbxxMsPomgbYqJf2VXKtiD8UKVuaHStjaUyRsZ1X2KjmFZ");
    ky := base58CheckDecodeKey(os.Args[1]);
    //_,ky := base58Decode("5HueCGU8rMjxEXxiPuD5BDku4MkFqeZyd4dZ1jvhTVqvbTLvyTJ");
    //ky,_ := hex.DecodeString("18E14A7B6A307F426A94F8114701E7C8E774E7F9A47E2C2035DB29A206321725");
    private_key, _ := btcec.PrivKeyFromBytes(btcec.S256(), ky);
    //private_key, err := ecdsa.GenerateKey(btcec.S256(), rand.Reader);
    /*private_key, err := hex.DecodeString("18E14A7B6A307F426A94F8114701E7C8E774E7F9A47E2C2035DB29A206321725");
    if err != nil {
        fmt.Println(err);
        return;
    }
    fmt.Println("privateKey: ", hex.EncodeToString(private_key.D.Bytes()));*/
    pub_ky_temp0 := make([]byte, 0, 65);
    pub_ky_temp0 = append(pub_ky_temp0, byte(4));
    pub_ky_temp0 = append(pub_ky_temp0, private_key.PublicKey.X.Bytes()...);
    pub_ky_temp0 = append(pub_ky_temp0, private_key.PublicKey.Y.Bytes()...);

    sha256_hash := sha256.New();
    sha256_hash.Write(pub_ky_temp0);
    pub_ky_temp0_sha256_hash := sha256_hash.Sum(nil);
    ripemd160_hash := ripemd160.New();
    ripemd160_hash.Write(pub_ky_temp0_sha256_hash);
    pub_key_temp1 := make([]byte, 0, 21);
    //pub_key_temp1 = append(pub_key_temp1, byte(0));
    pub_key_temp1 = append(pub_key_temp1, ripemd160_hash.Sum(nil)...);
    /*fmt.Println("Extended RIPEMD64 hash: ", hex.EncodeToString(pub_key_temp1));
    sha256_hash.Reset();
    sha256_hash.Write(pub_key_temp1);
    temp_hash := sha256_hash.Sum(nil);
    sha256_hash.Reset();
    sha256_hash.Write(temp_hash);
    checksum := sha256_hash.Sum(nil);
    fmt.Println("CheckSum :", hex.EncodeToString(checksum));
    public_key := append(pub_key_temp1, checksum[0:4]...);*/
    fmt.Println("publicKey (b58 encoding): ", base58CheckEncodeKey(byte(0), pub_key_temp1));
}

