# nkn-node-manager

## 功能描述

1. nkn 节点钱包存储
2. 后台自动检测节点运行情况，如关机超过一定时间则判断钱包为空闲状态
3. 新节点请求 API 获取空闲钱包
4. TODO: 自动给新节点发送10个 nkn 注册费（需要配置钱包私钥）
5. TODO: 前端界面


## 如何使用

1. 安装golang
2. git clone https://github.com/bufrr/nkn-node-manager.git
3. cd nkn-node-manager
4. go run main.go
5. 准备上传钱包

## 如何上传钱包

### 方法一

直接上传节点的钱包文件
curl -F "keystore=@/root/nkn-commercial/services/nkn-ndoe/wallet.json" -F "password=@/root/nkn-commercial/services/nkn-ndoe/wallet.pswd" http://x.x.x.x:30050/walletfile

请注意将 @ 符号后面的绝对路径替换为自己的钱包地址

### 方法二

使用python脚本一次性提交多个钱包
下面的例子中我将大量钱包+密码导入到了wallet.txt文件中，逐行读取并导入

```python
import requests
import json


def send_request():
    # upload wallet
    # POST http://127.0.0.1:30050/wallet
    with open("wallet.txt", "r") as wf:
        lines = wf.readlines()

    for w in lines:
        try:
            response = requests.post(
                url="http://x.x.x.x:30050/wallet",
                headers={
                    "Content-Type": "application/json; charset=utf-8",
                },
                data=json.dumps({
                    "keystore": [your wallet keystore],
                    "password": [your wallet password]
                })
            )
            print('Response HTTP Status Code: {status_code}'.format(
                status_code=response.status_code))
            print('Response HTTP Response Body: {content}'.format(
                content=response.content))
        except requests.exceptions.RequestException:
            print('HTTP Request failed')


if __name__ == '__main__':
    send_request()

```

## 如何获取空闲钱包

1. 安装 jq 命令
2. 通过 api 获取钱包，解析后存入本地

curl -s http://x.x.x.x:30050/wallet/idle > wallet.txt

cat wallet.txt | jq -r .idle.keystore > wallet.json

cat wallet.txt | jq -r .idle.password > wallet.pawd


##注意事项

如果你没有从钱包所在节点的IP上传钱包，会导致无法正常判断钱包状态（上传IP与钱包所在真实节点IP不一致），这种情况下需要自己将所有机器中的钱包删除后重新从 API 获取钱包












TG group: https://t.me/nknnodemanager
