# 在初始化golang依赖时使用s5代理

golang安装依赖的时候，会默认从google和github上拉取代码仓库。虽然我们可以用goproxy来指定网页的代理，但是有些时候我们的一部分依赖项并不在公网范围内，安装的时候如果开着goproxy就会安装失败。因此我们需要启用socks5。

```bash
https_proxy="socks5://127.0.0.1:10808" #视个人情况填写
go mod install
go get xxxxx/v2
```

```powershell
$env:https_proxy="socks5://127.0.0.1:10808"
go mod install
go get xxxxx/v2
```
