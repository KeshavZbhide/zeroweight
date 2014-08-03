package main

import "fmt"
import "bytes"
import "encoding/hex"
import "crypto/sha256"

var (
    fromAddress string = "1KhefwpMQuMJQnqqTFUn4vSGuVhtm87u2o";
    toAddress string = "1FRBPvxmK6pb58o4rpTD52zB4vWXNvNqLe";
    transferAmount uint64 = 1000000-200000;
    sourceTransaction string = "ecafe20a55a7661a5b8adf6f4adc4911a114437e0bd6f5a3bce47f0e28e3c10d";
    sourceIndex uint32 = 0;
    privateKeyHex string = "";
)


type tx_out struct {
    address string;
    amount uint64;
}

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


func reverseHexString(s string) string {
    length := len(s);
    result := make([]byte, len(s));
    for i := 0; i < length; i+=2 {
        result[i] = s[length-2-i];
        result[i+1] = s[length-1-i];
    }
    return string(result);
}

func Tx() string {
    input := make([]*tx_unspent, 1);
    input[0] = new(tx_unspent);
    input[0].tx_hash = reverseHexString(sourceTransaction);
    input[0].amount = 1000000;
    input[0].tx_output_n = 0;

    output := make([]*tx_out, 1);
    output[0] = new(tx_out);
    output[0].address = toAddress;
    output[0].amount = transferAmount;

    tx := append(makeRawTx(input, makeScriptPubKey(fromAddress), output), uint32Bytes(1)...);
    tx_hash1 := sha256.Sum256(tx);
    tx_hash := sha256.Sum256(tx_hash1[:]);
    return hex.EncodeToString(tx_hash[:]);
}

func main() {
    fmt.Println("rawTXhash ---> ", Tx());
}
