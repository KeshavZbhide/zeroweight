package main
import "fmt"
import "io/ioutil"
import "net/http"
import "encoding/json"
import "os"
import "errors"
import "strconv"

type tx_unspent struct {
    amount uint64;
    tx_hash string;
    tx_output_n int;
}

func get_unspent(addr string) ([]tx_unspent, error) {
    var ret []tx_unspent;
    res, _ := http.Get("http://blockchain.info/unspent?active="+addr);
    defer res.Body.Close();
    body, _ := ioutil.ReadAll(res.Body);
    var u interface{};
    json.Unmarshal(body, &u);
    unspent_temp, is_perfect_json := u.(map[string]interface{});
    if !is_perfect_json {
        fmt.Println(string(body));
        return nil, errors.New("Cannot Parse JSON");
    }
    unspent_main, contains := unspent_temp["unspent_outputs"];
    if contains {
        unspent := unspent_main.([]interface{});
        ret = make([]tx_unspent, len(unspent));
        for i := 0; i < len(unspent); i++ {
            unspent_tx := unspent[i].(map[string]interface{});
            for k, v := range unspent_tx {
                if k == "tx_hash" {
                    ret[i].tx_hash = v.(string);
                }
                if k == "tx_output_n" {
                    ret[i].tx_output_n = int(v.(float64));
                }
                if k == "value" {
                    ret[i].amount = uint64(v.(float64));
                }
            }
        }
    } else {
        return nil, errors.New("Unknown Json File");
    }
    return ret, nil;
}

func get_unspent_2(addr string, amount uint64) ([]*tx_unspent, uint64, error) {
    err_str := "unknown json response from blockchain.info";
    res, _ := http.Get("http://blockchain.info/unspent?active="+addr);
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
    unspent_main := unspent_temp2.([]interface{});
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
        result[i].tx_output_n = int(temp_.(float64));
        temp_, contains2 = unspent_tx["value"];
        if !contains {
            return nil, 0, errors.New(err_str);
        }
        result[i].amount = uint64(temp_.(float64));
    }
    sort_tx_unspent2(result);
    accumulate := uint64(0);
    for i := range result {
        accumulate += result[i].amount;
        if accumulate >= amount {
            return result[0:i+1], accumulate - amount, nil;
        }
    }
    return nil, 0, nil;
}

func sort_tx_unspent2(v []*tx_unspent) {
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

func sort_tx_unspent (unspent *[]tx_unspent) {
    for i := 0; i < len(*unspent); i++ {
        index := i
        for j := i; j < len(*unspent); j++ {
            if (*unspent)[index].amount < (*unspent)[j].amount {
                index = j;
            }
        }
        temp := (*unspent)[i];
        (*unspent)[i] = (*unspent)[index];
        (*unspent)[index] = temp;
    }
}

func get_unspent_tx_go() {
    if len(os.Args) < 3 {
        fmt.Println("Please Enter bitcoin address and the send amount in satoshi");
        return;
    }
    all, err := get_unspent(os.Args[1]);
    if err != nil {
        fmt.Println("Unable to Parse JSON file");
        fmt.Println(err);
        return;
    }
    sort_tx_unspent(&all);
    total_balance := uint64(0);
    for i := 0; i < len(all); i++ {
        total_balance += all[i].amount;
        fmt.Println("Index --> "+
            strconv.FormatUint(uint64(all[i].tx_output_n), 10)+
            " tx_hash --> "+all[i].tx_hash+" --> "+strconv.FormatUint(all[i].amount, 10));
    }
    fmt.Println("Total Balance --> ", total_balance);
    satoshi,_ := strconv.ParseUint(os.Args[2], 10, 64);
    use_this_tx := make(map[string]int);
    if (satoshi > total_balance) && (satoshi != 0) {
        fmt.Println("Insuffisent funds in account");
    } else {
        fmt.Println("\nAmount to Be Retrived == >", satoshi);
        enterd := uint64(0);
        for i := 0; i < len(all); i++ {
            use_this_tx[all[i].tx_hash] = all[i].tx_output_n;
            enterd += all[i].amount;
            if enterd > satoshi {
                break;
            }
        }
        fmt.Println("-----------------------------------------------------------");
        for k, v := range use_this_tx {
            fmt.Println(k+" --> ",v);
        }
        fmt.Println("-----------------------------------------------------------");
        fmt.Println("Amount back as change == >", enterd - satoshi);
    }
}

