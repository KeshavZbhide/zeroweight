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

