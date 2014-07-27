/*
* READ THE COMMENTS TO UNDERSTAND HOW BIT COIN WORKS!!!
* Raw Bitcoin Transaction Package.
*/

package main

import "bytes"
//import "strings"
import "encoding/hex"
import "fmt"
import "encoding/asn1"
import "github.com/conformal/btcec"
import "math/big"
import "crypto/sha256"

/*
* struct that represents a basics script public key and the amount he/her 
* should recive. 
*/
type txout struct{
    amount uint64;
    script_pub_key []byte;
}

/*--------------------------------------------Utility Functions-------------------------------------*/

const alphabet = "123456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz";
var bigRadix = big.NewInt(58);
var bigZero = big.NewInt(0);

func reverse(s []byte) []byte {
    length := len(s);
    result := make([]byte, len(s));
    for i := 0; i < length; i++ {
        result[i] = s[length-1-i];
    }
    return result;
}

func get_uint32_bytes(x uint32) []byte {
    return []byte{
        byte(x),
        byte(x >> 8),
        byte(x >> 16),
        byte(x >> 24),
    };
}

func varint(x uint32) []byte {
    if x < 253 {
        return []byte{byte(x)};
    } else if x < 65535 {
        return []byte{ byte(253), byte(x), byte(x>>8) }
    } else if x < 4294967295 {
        return []byte{ byte(254), byte(x), byte(x>>8), byte(x>>16), byte(x>>24) }
    } else {
        panic("varint for memmory greater than 4GB");
    }
}

/*---------------------------------------------------------------------------------------*/

func make_script_pub_key(pub_key []byte) []byte{
    var script_pub_key bytes.Buffer;
    script_pub_key.Write([]byte{118, 169, 20}); //OP_DUP , OP_HASH160, PUSHDATA 14
    script_pub_key.Write(pub_key); //publickey decoded base58chech encoding
    script_pub_key.Write([]byte{136, 172}); //OP_EQUALVERIFY , OP_CHECKSIG
    return script_pub_key.Bytes();
}

func make_outputs(outputs []txout) []byte {
    var formated_outputs bytes.Buffer;
    for i := 0; i < len(outputs); i++ {
        satoshi := []byte {
            byte( outputs[i].amount ),
            byte( outputs[i].amount >> 8 ),
            byte( outputs[i].amount >> 16 ),
            byte( outputs[i].amount >> 24 ),
            byte( outputs[i].amount >> 32 ),
            byte( outputs[i].amount >> 40 ),
            byte( outputs[i].amount >> 48 ),
            byte( outputs[i].amount >> 56 )}
        script_len := byte(len(outputs[i].script_pub_key));
        formated_outputs.Write(satoshi);
        formated_outputs.WriteByte(script_len);
        formated_outputs.Write(outputs[i].script_pub_key);
    }
    return formated_outputs.Bytes();
}

func make_raw_transaction(output_transaction_hash string, source_index uint32,
                            script_sig []byte, outputs []txout ) []byte {
    var tx bytes.Buffer;
    formated_outputs := make_outputs(outputs);
    output_transaction_hash_bytes,_ := hex.DecodeString(output_transaction_hash);
    tx.Write([]byte{1, 0, 0, 0}); //4 byte version
    tx.WriteByte(1);//Number of Inputs
    tx.Write(reverse(output_transaction_hash_bytes)); //previouse unspent transaction Hash.
    tx.Write(get_uint32_bytes(source_index)); //source index in unspent transaction.
    tx.WriteByte(byte(len(script_sig)));
    tx.Write(script_sig); //Script Signature
    tx.Write(get_uint32_bytes(4294967295)); //sequence
    tx.WriteByte(byte(len(outputs))); //
    tx.Write(formated_outputs);
    tx.Write([]byte{0, 0, 0, 0});
    return tx.Bytes();
}


/* 
* This Function would make the transaction and return the bytes required to send to the btc
* network. Sending the bytes to a btc node would publish the transaction if valid.
* [private_key] -> wif private key ie base58 check encoded.
* [public_key] -> wif base58checkencoded, sha256/RIPEM hash, ie just the normal public key that 
* one sees on the net, ex -> "1BTCorgHwCg6u2YSAWKgS17qUad6kHmtQW" appearing on 
* "https://bitcoinfoundation.org/donate/"
* [to_public_key] -> recivers address, wif public key
* [amount] -> btc in satoshis, 1 Satoshi ==  0.00000001 BTC.
* refer to this image to avoid confussion on key types 
* https://lh4.googleusercontent.com/-p8yVJXqY7fg/UuLaPjMDtyI/AAAAAAAAWYQ/QoenRIBO1O4/s2048/
* bitcoinkeys.png
*/
func Tx(private_k string, to_public_k string, amount uint64) []byte {
    output_tx_hash := "81b4c832d70cb56ff957589752eb4125a4cab78a25a8fc52d6a09e5bd4404d48";
    source_index := uint32(0);
    private_key, temp_pub_key := btcec.PrivKeyFromBytes(btcec.S256(), base58CheckDecodeKey(private_k));
    public_key := append(append([]byte{4}, temp_pub_key.X.Bytes()...), temp_pub_key.Y.Bytes()...);

    fmt.Println("--->\n",hex.EncodeToString(public_key),"\n<------\n");
    script_pub_key := make_script_pub_key(public_key);
    outputs := make([]txout, 1);
    outputs[0] = txout{amount, make_script_pub_key([]byte(base58CheckDecodeKey(to_public_k)))};
    raw_tx := append(
        make_raw_transaction(output_tx_hash, source_index, script_pub_key, outputs),
        get_uint32_bytes(16777216)...);
    //fmt.Println(hex.EncodeToString(raw_tx));
    s256_0 := sha256.Sum256(raw_tx);
    s256 := sha256.Sum256(s256_0[:]);
    temp_sig,_ := private_key.Sign(s256[:]);
    der_enc_temp,_ := asn1.Marshal(*temp_sig);
    //After appending the sighash all bytes, we get the final signature.
    signature := append(der_enc_temp, byte(1));
    //Now Build script sig bytes.
    var script_sig bytes.Buffer;
    script_sig.Write(varint(uint32(len(signature))));
    script_sig.Write(signature);
    script_sig.Write(varint(uint32(len(public_key))));
    script_sig.Write(public_key);
    signed_tx := make_raw_transaction(output_tx_hash, source_index, script_sig.Bytes(), outputs);
    return signed_tx;
}

func get_unspent_tx(public_key string) map[uint32]string {
    return make(map[uint32]string, 2);
}

func main() {
    tx := Tx("5HusYj2b2x4nroApgfvaSfKYZhRbKFH41bVyPooymbC6KfgSXdD",
                "1KKKK6N21XKo48zWKuQKXdvSsCf95ibHFa", uint64(91234));
    fmt.Println(hex.EncodeToString(tx));
    //get_unspent_tx_go();
}

