package main

import "bytes"
import "crypto/sha256"
import "strings"

var pszBase58 string = "123456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz";

func check_if_bytes_eq(s []byte, v []byte) bool {
    if len(s) != len(v) {
        return false;
    }
    for i := 0; i < len(s); i++ {
        if s[i] != v[i] {
            return false;
        }
    }
    return true;
}

func base58Encode(b []byte) string {
    zeros := 0;
    //count the leading zeros.
    for i := 0;  b[i] == 0; i++ {
        zeros++;
    }
    //Allocate Enoughf space in big-endian base58 representation
    b58 := make([]byte, (len(b)-zeros) * 138 / 100 + 1);
    //Process The Bytes
    for i := zeros; i < len(b); i++ {
        var carry int = int(b[i]);
        //Apply b58 = b58 * 256 + ch;
        for j := len(b58)-1; j > -1; j-- {
            carry += 256 * int(b58[j]);
            b58[j] = byte(carry % 58);
            carry /= 58;
        }
        //asert(carry == 0)
    }
    //skip leading zeros in base58 result
    index := 0
    for ; (index < len(b58)) && (b58[index] == 0); {
        index++;
    }
    var result bytes.Buffer;
    for i := 0; i < zeros; i++ {
        result.Write([]byte{byte('1')});
    }
    for ; index < len(b58); index++ {
        result.WriteByte(pszBase58[b58[(index)]]);
    }
    return result.String();
}

func base58CheckEncode(b []byte) string {
    var for_b58enc bytes.Buffer;
    hash1 := sha256.Sum256(b);
    hash2 := sha256.Sum256(hash1[:]);
    for_b58enc.Write(b);
    for_b58enc.Write(hash2[0:4]);
    return base58Encode(for_b58enc.Bytes());
}

func is_space(ch byte) bool {
    if((ch == ' ') || (ch == '\t') || (ch == '\n')) {
        return true;
    } else if ((ch == '\v') || (ch == '\f') || (ch == '\r')) {
        return true;
    }
    return false;
}

func base58Decode(s string) (bool, []byte) {
    str := s[:];
    index := 0;
    //skip white spaces
    for ; index < len(str) && (is_space(str[index])) ; index++ { }
    //skip and count leading ones
    var zeros int = 0;
    for ; (index < len(str)) && (str[index] == '1'); index++ {
        zeros++;
    }
    // Allocate enough space in big-endian base256 representation.
    b256 := make([]byte, len(str) * 733 / 1000 + 1);
    //process the characters...
    for ;(index < len(str)) && (!is_space(str[index])); index++ {
        //Decode base58 chracters...
        var carry int = strings.IndexByte(pszBase58, str[index]);
        if carry  == -1 {
            return false, nil;
        }
        //ch := pszBase58[carry];
        for i := len(b256)-1; i > -1; i-- {
            carry += 58 * int(b256[i]);
            b256[i] = byte(carry % 256);
            carry /= 256;
        }
        //assert(carry == 0);
    }
    //Remove trailing spaces
    for ;(index < len(str)) && (is_space(str[index])); index++ {}
    if index != len(str) {
        return false, nil;
    }
    //skip leading zeros in b256
    b256_index := 0
    for ;(b256_index < len(b256)) && (b256[b256_index] == 0); b256_index++ { }

    //copy result into output byte slice 
    output := make([]byte, zeros + (len(b256)-b256_index));
    output_index := 0;
    for ; output_index < zeros; output_index++ {
        output[output_index] = byte(0);
    }
    for ; b256_index < len(b256); b256_index++ {
        output[output_index] = b256[b256_index];
        output_index++;
    }
    return true, output;
}

func base58CheckDecode(s string) []byte {
    success, result := base58Decode(s);
    if (!success) || (len(result) < 4) {
        return nil;
    }
    s256_0 := sha256.Sum256(result[0:len(result)-4]);
    s256 := sha256.Sum256(s256_0[:]);
    if !check_if_bytes_eq(s256[0:4], result[len(result)-4:len(result)]) {
        return nil;
    }
    return result[0:len(result)-4];
}

func base58CheckEncodeKey(version byte, b []byte) string {
    return base58CheckEncode(append([]byte{version}, b...));
}

func base58CheckDecodeKey(s string) []byte {
    return base58CheckDecode(s)[1:];
}


