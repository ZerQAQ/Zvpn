# ZvpnImply
这是Zvpnlib的一个简单实例，将Zvpn封装成了简单的二进制文件

将目录clone到本地后输入go build即可编译

在命令行输入./ZvpnImply -g可以得到同目录下得到example_config文件：
```
{
    "client_address": "0.0.0.0:1080",
    "server_address": "0.0.0.0:1234",
    //对于客户端两个地址都必填，服务端可以只填server_address
    "protocol": "TCP",
    "proxy": "sock5",
    //protocol和proxy都是固定选项
    "algorithm": "RC4",
    "key": "d30948d73a274751db2354ec19a9a7a7439108bdfad0f6511b12a58b",
    //如果algorithm是RC4，那么key应该在8~28字节内
    //如果algorithm是AES，那么key可以是16、24或32字节，分别代表AES-128，AES-192和AES-256
    "role": "server or client"
    //role填server或者client
}
```
以上部分选项可以在命令行中输入
输入./ZvpnImply -help可以得到帮助：
```
Usage of ZvpnImply:
  -c string
        client address
  -f string
        config file path (default "config")
  -g    generate a example config file, save as config
  -r string
        role, server or client
  -s string
        server address
```