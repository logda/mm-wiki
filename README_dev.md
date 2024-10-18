初始化

```
cd install
go build -o install ./
go build -o install ./
./install -port 8090
```

热启动运行

```
air
```

打包后运行

```shell
go build ./
./mm-wiki # 默认读取 conf/mm-wiki.conf
./mm-wiki -conf ./conf/mm-wiki.conf # 指定配置文件
```

## 测试用脚本

通义千问接口

```
curl -X POST https://dashscope.aliyuncs.com/compatible-mode/v1/chat/completions \
-H "Authorization: Bearer YOUR_API_KEY" \
-H "Content-Type: application/json" \
-d '{
    "model": "qwen-plus",
    "messages": [
        {
            "role": "system",
            "content": "You are a helpful assistant."
        },
        {
            "role": "user",
            "content": "你是谁？"
        }
    ]
}'
```

通义千问接口流式

```
curl -X POST https://dashscope.aliyuncs.com/compatible-mode/v1/chat/completions \
-H "Authorization: Bearer YOUR_API_KEY" \
-H "Content-Type: application/json" \
-d '{
    "model": "qwen-plus",
    "messages": [
        {
            "role": "system",
            "content": "You are a helpful assistant."
        },
        {
            "role": "user",
            "content": "你是谁？"
        }
    ],
    "stream":true,
}'
```

文档问答接口验证

```
curl -X POST http://0.0.0.0:8080/api/ai-chat \
     -H "Content-Type: application/json" \
     -d '{"document": "测试文档", "message": "测试消息"}'
```
