package zeroweight
import "io/ioutil"
import "net/http"
import "encoding/json"
import "errors"

type tx_unspent struct {
    amount uint64;
    tx_hash string;
    tx_output_n uint32;
}

func getUnspent(addr string, amount uint64) ([]*tx_unspent, uint64, error) {
    err_str := "unknown json response from blockchain.info";
    res, err := http.Get("http://blockchain.info/unspent?active="+addr);
    if err != nil {
        return nil, 0, errors.New("unable to make net.http request");
    }
    defer res.Body.Close();
    body, _ := ioutil.ReadAll(res.Body);
    var u interface{};
    json.Unmarshal(body, &u);
    unspent_temp, is_perfect_json := u.(map[string]interface{});
    if !is_perfect_json {
        return nil, 0, errors.New(string(body));
    }
    unspent_temp2, contains := unspent_temp["unspent_outputs"];
    if !contains {
        return nil, 0, errors.New(err_str);
    }
    unspent_main, ok := unspent_temp2.([]interface{});
    if !ok {
        return nil, 0, errors.New(err_str);
    }
    if len(unspent_main) == 0 {
        return nil, 0, nil;
    }
    result := make([]*tx_unspent, len(unspent_main));
    for i := range unspent_main {
        result[i] = new(tx_unspent);
        unspent_tx := unspent_main[i].(map[string]interface{});
        temp_str_, contains := unspent_tx["tx_hash"];
        if !contains {
            return nil, 0, errors.New(err_str);
        }
        result[i].tx_hash = temp_str_.(string);
        temp_, contains2 := unspent_tx["tx_output_n"];
        if !contains2 {
            return nil, 0, errors.New(err_str);
        }
        result[i].tx_output_n = uint32(temp_.(float64));
        temp_, contains2 = unspent_tx["value"];
        if !contains {
            return nil, 0, errors.New(err_str);
        }
        result[i].amount = uint64(temp_.(float64));
    }
    sortTxUnspent(result);
    if amount == 0 {
        return result, 0, nil;
    }
    accumulate := uint64(0);
    for i := range result {
        accumulate += result[i].amount;
        if accumulate >= amount {
            return result[0:i+1], accumulate - amount, nil;
        }
    }
    return nil, 0, nil;
}

func sortTxUnspent(v []*tx_unspent) {
    v_len := len(v);
    for i := range v {
        index := i
        for j := i; j < v_len; j++ {
            if v[index].amount < v[j].amount {
                index = j;
            }
        }
        temp_ := v[i];
        v[i] = v[index];
        v[index] = temp_;
    }
}

