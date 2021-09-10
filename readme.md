1. 从/etc/esh_config.yaml读取信息，windows在 userAPP/roaming下不存在则创建
2. 密码和用户名ip可以用aes加密，应该至少有两个加解密aes密钥，随机选择(当前时间为随机种子)一个进行加密，根据开头的标识来进行解密判断，有一个root账户能看明文密码，密码为编译时的加盐值
3. 应该有登录补全功能和必须输入密钥才能登录功能
4. 可以像ssh一样能执行远程命令
5. esh clean命令清除所有有关esh的操作记录`history |grep history|awk '{print $1}'` 然后用`history -d`删除