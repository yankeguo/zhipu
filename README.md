# zhipu

[![Go Reference](https://pkg.go.dev/badge/github.com/yankeguo/zhipu.svg)](https://pkg.go.dev/github.com/yankeguo/zhipu)
[![go](https://github.com/yankeguo/zhipu/actions/workflows/go.yml/badge.svg)](https://github.com/yankeguo/zhipu/actions/workflows/go.yml)
[![codecov](https://codecov.io/gh/yankeguo/zhipu/graph/badge.svg?token=O08DOWX2TU)](https://codecov.io/gh/yankeguo/zhipu)

[中文文档](README.zh.md)

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
service := client.ChatCompletion("glm-4-flash").
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
service := client.ChatCompletion("glm-4-flash").
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

**ChatCompletion (Stream with GLM-4-AllTools)**

```go
// CodeInterpreter
s := client.ChatCompletion("GLM-4-AllTools")
s.AddMessage(zhipu.ChatCompletionMultiMessage{
    Role: "user",
    Content: []zhipu.ChatCompletionMultiContent{
        {
            Type: "text",
            Text: "计算[5,10,20,700,99,310,978,100]的平均值和方差。",
        },
    },
})
s.AddTool(zhipu.ChatCompletionToolCodeInterpreter{
    Sandbox: zhipu.Ptr(CodeInterpreterSandboxAuto),
})
s.SetStreamHandler(func(chunk zhipu.ChatCompletionResponse) error {
    for _, c := range chunk.Choices {
        for _, tc := range c.Delta.ToolCalls {
            if tc.Type == ToolTypeCodeInterpreter && tc.CodeInterpreter != nil {
                if tc.CodeInterpreter.Input != "" {
                    // DO SOMETHING
                }
                if len(tc.CodeInterpreter.Outputs) > 0 {
                    // DO SOMETHING
                }
            }
        }
    }
    return nil
})

// WebBrowser
// CAUTION: NOT 'WebSearch'
s := client.ChatCompletion("GLM-4-AllTools")
s.AddMessage(zhipu.ChatCompletionMultiMessage{
    Role: "user",
    Content: []zhipu.ChatCompletionMultiContent{
        {
            Type: "text",
            Text: "搜索下本周深圳天气如何",
        },
    },
})
s.AddTool(zhipu.ChatCompletionToolWebBrowser{})
s.SetStreamHandler(func(chunk zhipu.ChatCompletionResponse) error {
    for _, c := range chunk.Choices {
        for _, tc := range c.Delta.ToolCalls {
            if tc.Type == ToolTypeWebBrowser && tc.WebBrowser != nil {
                if tc.WebBrowser.Input != "" {
                    // DO SOMETHING
                }
                if len(tc.WebBrowser.Outputs) > 0 {
                    // DO SOMETHING
                }
            }
        }
    }
    return nil
})
s.Do(context.Background())

// DrawingTool
s := client.ChatCompletion("GLM-4-AllTools")
s.AddMessage(zhipu.ChatCompletionMultiMessage{
    Role: "user",
    Content: []zhipu.ChatCompletionMultiContent{
        {
            Type: "text",
            Text: "画一个正弦函数图像",
        },
    },
})
s.AddTool(zhipu.ChatCompletionToolDrawingTool{})
s.SetStreamHandler(func(chunk zhipu.ChatCompletionResponse) error {
    for _, c := range chunk.Choices {
        for _, tc := range c.Delta.ToolCalls {
            if tc.Type == ToolTypeDrawingTool && tc.DrawingTool != nil {
                if tc.DrawingTool.Input != "" {
                    // DO SOMETHING
                }
                if len(tc.DrawingTool.Outputs) > 0 {
                    // DO SOMETHING
                }
            }
        }
    }
    return nil
})
s.Do(context.Background())
```

**Embedding**

```go
service := client.Embedding("embedding-v2").SetInput("你好呀")
service.Do(context.Background())
```

**Image Generation**

```go
service := client.ImageGeneration("cogview-3").SetPrompt("一只可爱的小猫咪")
service.Do(context.Background())
```

**Video Generation**

```go
service := client.VideoGeneration("cogvideox").SetPrompt("一只可爱的小猫咪")
resp, err := service.Do(context.Background())

for {
    result, err := client.AsyncResult(resp.ID).Do(context.Background())

    if result.TaskStatus == zhipu.VideoGenerationTaskStatusSuccess {
        _ = result.VideoResult[0].URL
        _ = result.VideoResult[0].CoverImageURL
        break
    }

    if result.TaskStatus != zhipu.VideoGenerationTaskStatusProcessing {
        break
    }

    time.Sleep(5 * time.Second)
}
```

**Upload File (Retrieval)**

```go
service := client.FileCreate(zhipu.FilePurposeRetrieval)
service.SetLocalFile(filepath.Join("testdata", "test-file.txt"))
service.SetKnowledgeID("your-knowledge-id")

service.Do(context.Background())
```

**Upload File (Fine-Tune)**

```go
service := client.FileCreate(zhipu.FilePurposeFineTune)
service.SetLocalFile(filepath.Join("testdata", "test-file.jsonl"))
service.Do(context.Background())
```

**Batch Create**

```go
service := client.BatchCreate().
  SetInputFileID("fileid").
  SetCompletionWindow(zhipu.BatchCompletionWindow24h).
  SetEndpoint(BatchEndpointV4ChatCompletions)
service.Do(context.Background())
```

**Knowledge Base**

```go
client.KnowledgeCreate("")
client.KnowledgeEdit("")
```

**Fine Tune**

```go
client.FineTuneCreate("")
```

### Batch Support

**Batch File Writer**

```go
f, err := os.OpenFile("batch.jsonl", os.O_CREATE|os.O_WRONLY, 0644)

bw := zhipu.NewBatchFileWriter(f)

bw.Add("action_1", client.ChatCompletion("glm-4-flash").
    AddMessage(zhipu.ChatCompletionMessage{
        Role: "user",
        Content: "你好",
    }))
bw.Add("action_2", client.Embedding("embedding-v2").SetInput("你好呀"))
bw.Add("action_3", client.ImageGeneration("cogview-3").SetPrompt("一只可爱的小猫咪"))
```

**Batch Result Reader**

```go
br := zhipu.NewBatchResultReader[zhipu.ChatCompletionResponse](r)

for {
    var res zhipu.BatchResult[zhipu.ChatCompletionResponse]
    err := br.Read(&res)
    if err != nil {
        break
    }
}
```

## Donation

**This project is a personal open-source project maintained by GUO YANKE. The following donation channels are not related to Zhipu AI.**

**本项目是个人维护的开源项目，以下赞助渠道与智谱AI官方无关。**

Executing unit tests will actually call the ChatGLM API and consume my quota. Please donate and thank you for your support!

执行单元测试会真实调用GLM接口，消耗我充值的额度，开发不易，请微信扫码捐赠，感谢您的支持！

<img src="./wechat-donation.png" width="180"/>

## Credits

GUO YANKE, MIT License
