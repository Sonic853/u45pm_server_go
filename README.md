# u45pm_server_go

摊牌。

基于 [Golang login demo](https://github.com/jerryshell/golang-web-login) 制作改进：
 + 双端均以非明文保存密码，服务端以加盐形式保存密码
 + 需要带上通过登录成功后获得的key来发送请求
 + 两个客户端同时在线时，将强制将先登录的客户端踢出
 + 双端均使用正则过滤用户输入的非法字符
 + 房间系统
