package orderalipay

// import (
// 	//"alipayBase"
// 	"fmt"
// 	//"github.com/ljy2010a/go_alipay"
// 	"io/ioutil"
// 	"net/http"
// 	"net/http/httptest"
// 	"os"
// 	"testing"
// )

// func TestAlipay(t *testing.T) {
// 	initAlipayConfig()
// 	//ts := httptest.NewServer(http.HandleFunc(AlipayWebRequest))
// 	ts := httptest.NewServer(http.HandlerFunc(AlipayWebRequest))
// 	defer ts.Close()
// 	res, err := http.Get(ts.URL)
// 	defer res.Body.Close()
// 	if err != nil {
// 		t.Error("%s", err.Error())
// 	}
// 	got, err := ioutil.ReadAll(res.Body)
// 	if err != nil {
// 		t.Error(err)
// 	}
// 	fmt.Printf("%s\n", got)
// 	h, _ := os.OpenFile("ali.html", os.O_CREATE|os.O_RDWR, 0666)
// 	fmt.Fprintf(h, "%s", got)

// 	// client := &http.Client{}
// 	// request, _ := http.NewRequest("POST", urlStr, body)
// 	// response, _ := client.Do(request)
// 	// if response.StatusCode == 200 {
// 	// 	body, _ := ioutil.ReadAll(response.Body)
// 	// 	bodystr := string(body)
// 	// 	fmt.Println(bodystr)
// 	// }

// }
