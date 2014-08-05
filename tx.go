package main

import "fmt"
import "bytes"
import "encoding/hex"
import "crypto/sha256"
import "strconv"
import "os"
/*- dependancies, install with => $go get -*/
import "github.com/conformal/btcec"
import "code.google.com/p/go.crypto/ripemd160"

var (
    miningFees uint64 = 50000;
);

type tx_out struct {
    address string;
    amount uint64;
}

/*- >>>>>>>>>>>>>> Utility-Functions <<<<<<<<<<<<<<<<<< -*/

func uint32Bytes(n uint32) []byte {
    return []byte{
        byte(n),
        byte(n >> 8),
        byte(n >> 16),
        byte(n >> 24) };
}

func formatVarInt(a interface{}) []byte{
    formatVarUint32 := func (n uint32) []byte {
        if n < 0xfd {
            return []byte{ byte(n) };
        } else if n < 0xffff {
            return []byte{ byte(0xfd), byte(n), byte(n >> 8) };
        } else if n < 0xffffffff {
            return []byte{
                byte(0xfe),
                byte(n),
                byte(n >> 8),
                byte(n >> 16),
                byte(n >> 24) };
        }
        return nil;
    }
    if v,ok := a.(uint32); ok {
        return formatVarUint32(v);
    } else if v, ok := a.(uint64); ok {
        return formatVarUint32(uint32(v));
    } else if v, ok := a.(int); ok {
        return formatVarUint32(uint32(v));
    }
    return nil;
}

func publicKeyStructToPublicKeyBytes(key *btcec.PublicKey) []byte {
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
/*-------------------------------------------------------*/


func formatInputs(input []*tx_unspent, script []byte) []byte {
    var formatedInput bytes.Buffer;
    for _, v := range input {
        tx_hash,_ := hex.DecodeString(v.tx_hash);
        formatedInput.Write(tx_hash);
        formatedInput.Write(uint32Bytes(v.tx_output_n));
        formatedInput.Write(formatVarInt(len(script)));
        formatedInput.Write(script);
        formatedInput.Write(uint32Bytes(0xffffffff));   //sequence
    }
    return formatedInput.Bytes();
}

func formatOutputs(output []*tx_out) []byte {
    var formatedOutputs bytes.Buffer;
    for _, v := range output {
        scriptPubKey := makeScriptPubKey(v.address);
        formatedOutputs.Write(uint32Bytes(uint32(v.amount)));
        formatedOutputs.Write(uint32Bytes(uint32((v.amount) >> 32)));
        formatedOutputs.Write(formatVarInt(len(scriptPubKey)));
        formatedOutputs.Write(scriptPubKey);
    }
    return formatedOutputs.Bytes();
}

func makeScriptPubKey(addr string) []byte {
    pubkey_hash := base58CheckDecodeKey(addr);
    scriptPubKey := make([]byte, 0, len(pubkey_hash)+5);
    scriptPubKey = append(scriptPubKey, 0x76, 0xa9, 0x14);
    scriptPubKey = append(scriptPubKey, pubkey_hash...);
    scriptPubKey = append(scriptPubKey, 0x88, 0xac);
    return scriptPubKey;
}

func makeRawTx(inputs []*tx_unspent, script []byte, outputs []*tx_out) []byte {
    formatedInputs := formatInputs(inputs, script);
    formatedOutputs := formatOutputs(outputs);
    var tx bytes.Buffer;
    tx.Write([]byte{1, 0, 0, 0});           //4 Byte version
    tx.Write(formatVarInt(len(inputs)));    //# of inputs
    tx.Write(formatedInputs);               //formated inputs
    tx.Write(formatVarInt(len(outputs)));   //# of outputs
    tx.Write(formatedOutputs);              //formated outputs
    tx.Write([]byte{0, 0, 0, 0});           //block lock time
    return tx.Bytes();
}

/*- 
- Exported: 
-*/
func Tx(fromPrivateKeyWif string, toAddress string, amount float64) (string, error) {
    var signedScript bytes.Buffer;
    transferAmount := uint64(amount * 100000000);
    /*- all required key types -*/
    fromPrivateKeyBytes := base58CheckDecodeKey(fromPrivateKeyWif);
    fromPrivateKeyStruct, fromPublicKeyStruct := btcec.PrivKeyFromBytes(
                                                    btcec.S256(), fromPrivateKeyBytes);
    fromPublicKeyBytes := publicKeyStructToPublicKeyBytes(fromPublicKeyStruct);
    fromPublicKeyBase58 := base58CheckEncodeKey(byte(0), publicKeyHash(fromPublicKeyBytes));
    /*- build inputs and outputs -*/
    input, change, err := getUnspent(fromPublicKeyBase58, transferAmount);
    if err != nil {
        return err.Error(), err;
    }
    output := make([]*tx_out, 1, 2);
    output[0] = new(tx_out);
    output[0].address = toAddress;
    output[0].amount = transferAmount;
    /*- add the mining fees -*/
    if change > miningFees {
        outputChange := new(tx_out);
        outputChange.address = fromPublicKeyBase58;
        outputChange.amount = change - miningFees;
        output = append(output, outputChange);
    } else {
        output[0].amount -= miningFees;
    }
    /*- build hash of raw transaction for signing -*/
    tx := append(makeRawTx(input, makeScriptPubKey(fromPublicKeyBase58), output), uint32Bytes(1)...);
    tx_hash0 := sha256.Sum256(tx);
    tx_hash := sha256.Sum256(tx_hash0[:]);
    tempSig,_ := fromPrivateKeyStruct.Sign(tx_hash[:]);
    signature := append(tempSig.Serialize(), byte(1));
    /*- build script_sig -*/
    signedScript.Write(formatVarInt(len(signature)));
    signedScript.Write(signature);
    signedScript.Write(formatVarInt(len(fromPublicKeyBytes)));
    signedScript.Write(fromPublicKeyBytes);
    /*- return the hex-encoded signed transaction -*/
    return hex.EncodeToString(makeRawTx(input, signedScript.Bytes(), output));
}

/*- 
- Exported:
-*/
func Balance(publicKey string) float64, error {
    var satoshis uint64 = 0;
    unspent, _, err := getUnspent(publicKey);
    if err != nil {
        return 0.0, err;
    }
    for _,v := range unspent {
        satoshis += v.amount;
    }
    return float64(satoshi) * float64(0.00000001);
}

/*-
- Exported
-*/
func GenRandPrivateKey() string {

}

/*-
- Exported
*/
func GetPublicKey(privateKeyWif string) string {
}
