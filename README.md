# zhipu

A 3rd-Party Golang Client Library for Zhipu AI Platform

## Usage

1. Install the package

```bash
go get -u github.com/yankeguo/zhipu
```

2. Create a client

```go
// this will use environment variables ZHIPUAI_API_KEY
client := zhipu.NewClient()
// or you can specify the API key
client = zhipu.NewClient(zhipu.WithAPIKey("your api key"))
```

3. Use the client

**ChatCompletion**

```go
service := client.ChatCompletionService("glm-4-flash").
    AddMessage(zhipu.ChatCompletionMessage{
        Role: "user",
        Content: "你好",
    })
res, err := service.Do(context.Background())
println(res.Choices[0].Message.Content)
```

**ChatCompletion(Stream)**

```go
service := client.ChatCompletionService("glm-4-flash").
    AddMessage(zhipu.ChatCompletionMessage{
        Role: "user",
        Content: "你好",
    }).SetStreamHandler(func(chunk zhipu.ChatCompletionStreamResponse) error {
        println(chunk.Choices[0].Delta.Content)
        return nil
    })
res, err := service.Do(context.Background())
println(res.Choices[0].Message.Content)
```
> [!NOTE]
>
> More APIs are coming soon.

## Credits

GUO YANKE, MIT License
