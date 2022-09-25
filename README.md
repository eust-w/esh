## 🎉简介
esh 是一个跨平台的ssh链接管理工具，简单且强大！

## ⚡使用
安装
## 📜从源码安装
二进制文件生成在 out 目录下。

## 📦下载安装
x86-64 linux版本： [esh-linux-amd64](https://github.com/eust-w/esh/releases/download/v/esh-linux-amd64)

arm-64 linux版本: [esh-linux-arm64](https://github.com/eust-w/esh/releases/download/v/esh-linux-arm64)

x86-64 mac版本：[esh-mac-amd64](https://github.com/eust-w/esh/releases/download/v/esh-mac-amd64)

x86-64 windows版本：[esh.exe](https://github.com/eust-w/esh/releases/download/v/esh.exe)

下载后可直接运行，注意!请通过命令行运行！

## 🌱交互
esh	说明
```
  add         add remote ssh                                            
  cluster     use connect to connect remote ssh or run command          
  completion  generate the autocompletion script for the specified shell
  del         del an remote ssh use name                                
  help        Help about any command                                    
  list        list remote ssh                                           
  run         use connect to connect remote ssh or run command          
  set         set global config 
```

## ➕开发
从Home/esh_config.yaml读取信息
2. 密码和用户名ip可以用aes加密，应该至少有两个加解密aes密钥，随机选择(当前时间为随机种子)一个进行加密，根据开头的标识来进行解密判断，有一个root账户能看明文密码，密码为编译时的加盐值
3. 应该有登录补全功能和必须输入密钥才能登录功能
4. 可以像ssh一样能执行远程命令,且支持集群功能