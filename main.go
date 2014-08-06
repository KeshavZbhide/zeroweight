package main 
import "os"
import "fmt"
import "bytes"
import "crypto/aes"
import "os/user"
import "strconv"
import "errors"
import "encoding/hex"
import "io/ioutil"

/*- 
- Basic usage 
-
- $zeroweight createWallet [password] :
- this creates a file in user's home directory "zeroweight.wal" please backup this file
- in your dropbox/gdrive account, since it contains your private keys. Losing this file 
- would cause you to lose all your bitcoins. BE VERY CAREFULL.
- 
- $zeroweight send [toAddress] [BTC amount] [password]:
- sends butcoin from your wallet to [toAddress]. only users with the right password can
- can make transactions. Transactions are brodcasted to the bitcoin network via 
- https://blockchain.info/pushtx.
-
- $zeroweight balance [password]:
- shows the current blance and your public key (ie your wallet's address).
-*/

func main() {
    printUsage := func() {
        fmt.Println("# use => $zeroweight createWallet [encryptionKey|password]");
        fmt.Println("#     => $zeroweight send [toAddress] [BTC amount] [password]");
        fmt.Println("#     => $zeroweight balance [password]");
        return;
    }
    if len(os.Args) < 3 {
        printUsage();
        return;
    }
    switch os.Args[1] {
        /*- creates wallet -*/
        case "createWallet":
            file, err := encryptAndBuildWallet(GenRandPrivateKey(), os.Args[2]);
            if err != nil {
                fmt.Println("# error =>", err.Error());
                return;
            }
            fmt.Println("# success => wallet built.")
            fmt.Println("#         => please backup", file, "with dropbox|gdrive");
        /*- builds and broadcasts transaction -*/
        case "send":
            if len(os.Args) < 5 {
                printUsage();
                return;
            }
            walletPrivateKeyWif, err := decryptAndGetPrivateKey(os.Args[4]);
            if err != nil {
                fmt.Println("# error =>", err);
                return;
            }
            amount,err := strconv.ParseFloat(os.Args[3], 64)
            if (err != nil) || (amount == 0) {
                fmt.Println("# error => unable to parse amount enterd");
                return;
            }
            b, err := Balance(GetPublicKey(walletPrivateKeyWif));
            if err != nil {
                fmt.Println("# error =>", err.Error());
                return;
            }
            if (b < amount) {
                fmt.Println("# error => your wallet does not have enoughf balance");
                fmt.Println("#       => execute:$ zeroweight balance [password] to",
                            "check balance");
                return;
            }
            tx, err := Tx(walletPrivateKeyWif, os.Args[2], amount);
            if err != nil {
                fmt.Println("# error =>", err.Error());
                return;
            }
            res := SubmitTransaction(tx);
            fmt.Println("# status =>", res);
        /*- prints balance -*/
        case "balance":
            walletPrivateKeyWif, err := decryptAndGetPrivateKey(os.Args[2]);
            if err != nil {
                fmt.Println("# error =>", err.Error());
                return;
            }
            walletPublicKeyWif := GetPublicKey(walletPrivateKeyWif);
            fmt.Println("# public address =>", walletPublicKeyWif);
            balance, err := Balance(walletPublicKeyWif);
            if err != nil {
                fmt.Println("# error =>", err.Error());
                return;
            }
            fmt.Println("# balance => ", balance);
        default:
            printUsage();
            return;
    }
}

func pathExist(path string) (bool, os.FileInfo) {
    info, err :=  os.Stat(path)
    if err != nil {
        if os.IsNotExist(err) {
            return false, info;
        }
        panic("panic => cannot determine if path exist");
    }
    return true, info ;
}

func encryptAndBuildWallet(privateKey string, password string) (string, error) {
    var userDir string;
    if user,err := user.Current(); err != nil {
        return "", errors.New("unable to lookup user directory");
    } else {
        userDir = user.HomeDir;
    }
    if exist,_ := pathExist(userDir+"/zeroweight.wal"); exist {
        return "", errors.New("wallet already exists");
    }
    walletFileContent := "key{"+privateKey+"}";
    if len(password) < 6 {
        return "", errors.New("password|encryption key should be atleast 6 characters");
    }
    encryptionKey := make([]byte, 16);
    copy(encryptionKey, []byte(password));
    aesCipher, err := aes.NewCipher(encryptionKey);
    if err != nil {
        return "", err;
    }
    cBlockLen := aesCipher.BlockSize();
    toEncryptLen := len(walletFileContent) +
        (cBlockLen - (len(walletFileContent) % cBlockLen));
    toEncrypt := make([]byte, toEncryptLen);
    copy(toEncrypt, walletFileContent);
    for i := 0; i < toEncryptLen; i += cBlockLen {
        slice := toEncrypt[i:(i+cBlockLen)];
        aesCipher.Encrypt(slice, slice);
    }
    wallet := hex.EncodeToString(toEncrypt);
    err = ioutil.WriteFile(userDir+"/zeroweight.wal", []byte(wallet), 0644);
    return userDir+"/zeroweight.wal", err;
}

func decryptAndGetPrivateKey(pass string) (string, error) {
    var userDir string;
    if user, err := user.Current(); err != nil {
        return "", errors.New("unable to lookup user directory");
    } else {
        userDir = user.HomeDir;
    }
    if exist,_ := pathExist(userDir+"/zeroweight.wal"); !exist {
        return "", errors.New("no wallet created, exec $zeroweight createWallet");
    }
    file, err := ioutil.ReadFile(userDir+"/zeroweight.wal");
    if err != nil {
        return "", err;
    }
    wallet, err := hex.DecodeString(string(file));
    if err != nil {
        return "", err;
    }
    walletLen := len(wallet);
    encryptionKey := make([]byte, 16);
    copy(encryptionKey, []byte(pass));
    aesCipher, err := aes.NewCipher(encryptionKey);
    if err != nil {
        return "", err;
    }
    cBlockLen := aesCipher.BlockSize();
    for i := 0; i < walletLen; i += cBlockLen {
        slice := wallet[i:(i+cBlockLen)];
        aesCipher.Decrypt(slice, slice);
    }
    if string(wallet[0:4]) != "key{" {
        return "", errors.New("wrong password|encryptionKey");
    }
    last := bytes.IndexByte(wallet, '}');
    if last == -1 {
        return "", errors.New("corrupt wallet file, try again with right password");
    }
    key := string(wallet[4:last]);
    return key, nil;
}

