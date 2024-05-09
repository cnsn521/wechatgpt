# wechatgpt
## 支持微信的chatgpt机器人

### **基于[aurora](https://github.com/aurora-develop/aurora) 项目实现chatgpt**

### **基于[wechatbot](https://github.com/danni-cool/wechatbot-webhook)项目实现的微信机器人**

**步骤1：** 先根据wechatbot的文档搭建微信机器人，并完成接收消息的配置，我这边配置的是http://127.0.0.1:3002/msg去接受和处理微信发送过来的消息

**步骤2：** 根据aurora部署好gpt（需要支持chatgpt的服务器或者相应的代理），并配置好域名或者接口信息

**步骤3：** 手动更改内部的接口地址，我这边配置的https://example.com,可以根据实际情况进行更改

剩下就是自己玩玩吧
