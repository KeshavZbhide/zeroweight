package main
import "os"
import "fmt"

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
        fmt.Println("usage => $zeroweight createWallet [encryptionKey|password]");
        fmt.Println("      => $zeroweight send [toAddress] [BTC amount] [password]");
        fmt.Println("      => $zeroweight balance [password]");
        return;
    }
    if len(os.Args) < 3 {
        printUsage();
        return;
    }
    switch os.Args[1] {
        /*- creates wallet -*/
        case "createWallet":
            walletPrivateKeyWif := "key=" + GenRandPrivateKey();
            if len(os.Args[2]) < 6 {
               fmt.Println("your encryptionKey/password should be atleast 6 characters");
               return;
            }
            content := encryptString(walletPrivateKeyBase58, os.Args[2]);
            walletUser, err := user.Current();
            if err != nil {
                fmt.Println("error => unable to get user directory");
                return;
            }
            err := ioutil.WriteFile(walletUser.HomeDir+"/zeroweight.wal", content, 0644);
            if err != nil {
                fmt.Println("error => unable to write file");
                return;
            }
        /*- builds and broadcasts transaction -*/
        case "send":
            if len(os.Args) < 5 {
                printUsage();
                return;
            }
            walletPrivateKeyWif, err := decryptAndGetPrivateKey(os.Args[4]);
            if err != nil {
                fmt.Println("error =>", err);
                return;
            }
            amount,err := strconv.ParseFloat(os.Args[3], 64)
            if (err != nil) || (amount == 0) {
                fmt.Println("error => unable to parse amount enterd");
                return;
            }
            if Balance(walletPrivateKeyWif) < amount {
                fmt.Println("error => your wallet does not have enoughf balance");
                fmt.Println("      => execute:$ zeroweight balance [password] to check balance");
                return;
            }
            tx, err := Tx(walletPrivateKeyWif, os.Args[2], amount);
            if err != nil {
                fmt.Println("error =>",tx);
                return;
            }
            res, err = SubmitTransaction(tx);
            if err != nil {
                fmt.Println("error =>", err.Error());
                return;
            }
            fmt.Println("success =>", res);
        /*- prints balance -*/
        case "balance":
            walletPrivateKeyWif, err := decryptAndGetPrivateKey(os.Args[2]);
            if err != nil {
                fmt.Println("error =>", err.Error());
                return;
            }
            walletPublicKeyWif := GetPublicKey(walletPrivateKeyWif);
            balance, err := Balance(wallerPublicKeyWif);
            if err != nil {
                fmt.Println("error =>", err.Error());
                return;
            }
            fmt.Println("public address =>", walletPublicKeyWif);
            fmt.Println("balance => ", balance);
    }
}
