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
        case "createWallet":
            walletPrivateKeyBase58 := "key=" + genRandomPrivateKey();
            if len(os.Args[2]) < 6 {
               fmt.Println("your encryptionKey/password should be atleast 6 characters");
               return;
            }
            content := encryptString(walletPrivateKeyBase58, os.Args[2]);
            walletUser, err := user.Current();
            if err != nil {
                fmt.Println("unable to get user directory");
                return;
            }
            err := ioutil.WriteFile(walletUser.HomeDir+"/zeroweight.wal", content, 0644);
            if err != nil {
                fmt.Println("unable to write file");
                return;
            }
        case "send":
        case "balance":
    }
    /*amount, err := strconv.ParseFloat(os.Args[3], 64);
    if err != nil {
        fmt.Println("unable to parse [BTC amount] specified");
        return;
    }
    fmt.Println("tx => ");
    fmt.Println(Tx(os.Args[1], os.Args[2], amount));*/
}
