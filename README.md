# zhipu

[![Go Reference](https://pkg.go.dev/badge/github.com/yankeguo/zhipu.svg)](https://pkg.go.dev/github.com/yankeguo/zhipu)
[![go](https://github.com/yankeguo/zhipu/actions/workflows/go.yml/badge.svg)](https://github.com/yankeguo/zhipu/actions/workflows/go.yml)

A 3rd-Party Golang Client Library for Zhipu AI Platform

## Usage

### Install the package

```bash
go get -u github.com/yankeguo/zhipu
```

### Create a client

```go
// this will use environment variables ZHIPUAI_API_KEY
client, err := zhipu.NewClient()
// or you can specify the API key
client, err = zhipu.NewClient(zhipu.WithAPIKey("your api key"))
```

### Use the client

**ChatCompletion**

```go
service := client.ChatCompletionService("glm-4-flash").
    AddMessage(zhipu.ChatCompletionMessage{
        Role: "user",
        Content: "你好",
    })

res, err := service.Do(context.Background())

if err != nil {
    zhipu.GetAPIErrorCode(err) // get the API error code
} else {
    println(res.Choices[0].Message.Content)
}
```

**ChatCompletion (Stream)**

```go
service := client.ChatCompletionService("glm-4-flash").
    AddMessage(zhipu.ChatCompletionMessage{
        Role: "user",
        Content: "你好",
    }).SetStreamHandler(func(chunk zhipu.ChatCompletionResponse) error {
        println(chunk.Choices[0].Delta.Content)
        return nil
    })

res, err := service.Do(context.Background())

if err != nil {
    zhipu.GetAPIErrorCode(err) // get the API error code
} else {
    // this package will combine the stream chunks and build a final result mimicking the non-streaming API
    println(res.Choices[0].Message.Content)
}
```

**Embedding**

```go
service := client.EmbeddingService("embedding-v2").SetInput("你好呀")
service.Do(context.Background())
```

**Image Generation**

```go
service := client.ImageGenerationService("cogview-3").SetPrompt("一只可爱的小猫咪")
service.Do(context.Background())
```

**Upload File (Retrieval)**

```go
service := client.FileCreateService(zhipu.FilePurposeRetrieval)
service.SetLocalFile(filepath.Join("testdata", "test-file.txt"))
service.SetKnowledgeID("your-knowledge-id")

service.Do(context.Background())
```

**Upload File (Fine-Tune)**

```go
service := client.FileCreateService(zhipu.FilePurposeFineTune)
service.SetLocalFile(filepath.Join("testdata", "test-file.jsonl"))
service.Do(context.Background())
```

> [!NOTE]
>
> More APIs are coming soon.

## Credits

GUO YANKE, MIT License
