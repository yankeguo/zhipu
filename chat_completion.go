package zhipu

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"

	"github.com/go-resty/resty/v2"
)

const (
	RoleSystem    = "system"
	RoleUser      = "user"
	RoleAssistant = "assistant"
	RoleTool      = "tool"

	ToolChoiceAuto = "auto"

	FinishReasonStop         = "stop"
	FinishReasonStopSequence = "stop_sequence"
	FinishReasonToolCalls    = "tool_calls"
	FinishReasonLength       = "length"
	FinishReasonSensitive    = "sensitive"
	FinishReasonNetworkError = "network_error"

	ToolTypeFunction  = "function"
	ToolTypeWebSearch = "web_search"
	ToolTypeRetrieval = "retrieval"

	MultiContentTypeText     = "text"
	MultiContentTypeImageURL = "image_url"

	// New in GLM-4-AllTools
	ToolTypeCodeInterpreter = "code_interpreter"
	ToolTypeDrawingTool     = "drawing_tool"
	ToolTypeWebBrowser      = "web_browser"

	CodeInterpreterSandboxNone = "none"
	CodeInterpreterSandboxAuto = "auto"

	ChatCompletionStatusFailed         = "failed"
	ChatCompletionStatusCompleted      = "completed"
	ChatCompletionStatusRequiresAction = "requires_action"
)

// ChatCompletionTool is the interface for chat completion tool
type ChatCompletionTool interface {
	isChatCompletionTool()
}

// ChatCompletionToolFunction is the function for chat completion
type ChatCompletionToolFunction struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Parameters  any    `json:"parameters"`
}

func (ChatCompletionToolFunction) isChatCompletionTool() {}

// ChatCompletionToolRetrieval is the retrieval for chat completion
type ChatCompletionToolRetrieval struct {
	KnowledgeID    string `json:"knowledge_id"`
	PromptTemplate string `json:"prompt_template,omitempty"`
}

func (ChatCompletionToolRetrieval) isChatCompletionTool() {}

// ChatCompletionToolWebSearch is the web search for chat completion
type ChatCompletionToolWebSearch struct {
	Enable       *bool  `json:"enable,omitempty"`
	SearchQuery  string `json:"search_query,omitempty"`
	SearchResult bool   `json:"search_result,omitempty"`
}

func (ChatCompletionToolWebSearch) isChatCompletionTool() {}

// ChatCompletionToolCodeInterpreter is the code interpreter for chat completion
// only in GLM-4-AllTools
type ChatCompletionToolCodeInterpreter struct {
	Sandbox *string `json:"sandbox,omitempty"`
}

func (ChatCompletionToolCodeInterpreter) isChatCompletionTool() {}

// ChatCompletionToolDrawingTool is the drawing tool for chat completion
// only in GLM-4-AllTools
type ChatCompletionToolDrawingTool struct {
	// no fields
}

func (ChatCompletionToolDrawingTool) isChatCompletionTool() {}

// ChatCompletionToolWebBrowser is the web browser for chat completion
type ChatCompletionToolWebBrowser struct {
	// no fields
}

func (ChatCompletionToolWebBrowser) isChatCompletionTool() {}

// ChatCompletionUsage is the usage for chat completion
type ChatCompletionUsage struct {
	PromptTokens     int64 `json:"prompt_tokens"`
	CompletionTokens int64 `json:"completion_tokens"`
	TotalTokens      int64 `json:"total_tokens"`
}

// ChatCompletionWebSearch is the web search result for chat completion
type ChatCompletionWebSearch struct {
	Icon    string `json:"icon"`
	Title   string `json:"title"`
	Link    string `json:"link"`
	Media   string `json:"media"`
	Content string `json:"content"`
}

// ChatCompletionToolCallFunction is the function for chat completion tool call
type ChatCompletionToolCallFunction struct {
	Name      string          `json:"name"`
	Arguments json.RawMessage `json:"arguments"`
}

// ChatCompletionToolCallCodeInterpreterOutput is the output for chat completion tool call code interpreter
type ChatCompletionToolCallCodeInterpreterOutput struct {
	Type string `json:"type"`
	Logs string `json:"logs"`
	File string `json:"file"`
}

// ChatCompletionToolCallCodeInterpreter is the code interpreter for chat completion tool call
type ChatCompletionToolCallCodeInterpreter struct {
	Input   string                                        `json:"input"`
	Outputs []ChatCompletionToolCallCodeInterpreterOutput `json:"outputs"`
}

// ChatCompletionToolCallDrawingToolOutput is the output for chat completion tool call drawing tool
type ChatCompletionToolCallDrawingToolOutput struct {
	Image string `json:"image"`
}

// ChatCompletionToolCallDrawingTool is the drawing tool for chat completion tool call
type ChatCompletionToolCallDrawingTool struct {
	Input   string                                    `json:"input"`
	Outputs []ChatCompletionToolCallDrawingToolOutput `json:"outputs"`
}

// ChatCompletionToolCallWebBrowserOutput is the output for chat completion tool call web browser
type ChatCompletionToolCallWebBrowserOutput struct {
	Title   string `json:"title"`
	Link    string `json:"link"`
	Content string `json:"content"`
}

// ChatCompletionToolCallWebBrowser is the web browser for chat completion tool call
type ChatCompletionToolCallWebBrowser struct {
	Input   string                                   `json:"input"`
	Outputs []ChatCompletionToolCallWebBrowserOutput `json:"outputs"`
}

// ChatCompletionToolCall is the tool call for chat completion
type ChatCompletionToolCall struct {
	ID              string                                 `json:"id"`
	Type            string                                 `json:"type"`
	Function        *ChatCompletionToolCallFunction        `json:"function,omitempty"`
	CodeInterpreter *ChatCompletionToolCallCodeInterpreter `json:"code_interpreter,omitempty"`
	DrawingTool     *ChatCompletionToolCallDrawingTool     `json:"drawing_tool,omitempty"`
	WebBrowser      *ChatCompletionToolCallWebBrowser      `json:"web_browser,omitempty"`
}

type ChatCompletionMessageType interface {
	isChatCompletionMessageType()
}

// ChatCompletionMessage is the message for chat completion
type ChatCompletionMessage struct {
	Role       string                   `json:"role"`
	Content    string                   `json:"content,omitempty"`
	ToolCalls  []ChatCompletionToolCall `json:"tool_calls,omitempty"`
	ToolCallID string                   `json:"tool_call_id,omitempty"`
}

func (ChatCompletionMessage) isChatCompletionMessageType() {}

type ChatCompletionMultiContent struct {
	Type     string   `json:"type"`
	Text     string   `json:"text"`
	ImageURL *URLItem `json:"image_url,omitempty"`
}

// ChatCompletionMultiMessage is the multi message for chat completion
type ChatCompletionMultiMessage struct {
	Role    string                       `json:"role"`
	Content []ChatCompletionMultiContent `json:"content"`
}

func (ChatCompletionMultiMessage) isChatCompletionMessageType() {}

// ChatCompletionMeta is the meta for chat completion
type ChatCompletionMeta struct {
	UserInfo string `json:"user_info"`
	BotInfo  string `json:"bot_info"`
	UserName string `json:"user_name"`
	BotName  string `json:"bot_name"`
}

// ChatCompletionChoice is the choice for chat completion
type ChatCompletionChoice struct {
	Index        int                   `json:"index"`
	FinishReason string                `json:"finish_reason"`
	Delta        ChatCompletionMessage `json:"delta"`   // stream mode
	Message      ChatCompletionMessage `json:"message"` // non-stream mode
}

// ChatCompletionResponse is the response for chat completion
type ChatCompletionResponse struct {
	ID        string                    `json:"id"`
	Created   int64                     `json:"created"`
	Model     string                    `json:"model"`
	Choices   []ChatCompletionChoice    `json:"choices"`
	Usage     ChatCompletionUsage       `json:"usage"`
	WebSearch []ChatCompletionWebSearch `json:"web_search"`
	// Status is the status of the chat completion, only in GLM-4-AllTools
	Status string `json:"status"`
}

// ChatCompletionStreamHandler is the handler for chat completion stream
type ChatCompletionStreamHandler func(chunk ChatCompletionResponse) error

var (
	chatCompletionStreamPrefix = []byte("data:")
	chatCompletionStreamDone   = []byte("[DONE]")
)

// chatCompletionReduceResponse reduce the chunk to the response
func chatCompletionReduceResponse(out *ChatCompletionResponse, chunk ChatCompletionResponse) {
	if len(out.Choices) == 0 {
		out.Choices = append(out.Choices, ChatCompletionChoice{})
	}

	// basic
	out.ID = chunk.ID
	out.Created = chunk.Created
	out.Model = chunk.Model

	// choices
	if len(chunk.Choices) != 0 {
		oc := &out.Choices[0]
		cc := chunk.Choices[0]

		oc.Index = cc.Index
		if cc.Delta.Role != "" {
			oc.Message.Role = cc.Delta.Role
		}
		oc.Message.Content += cc.Delta.Content
		oc.Message.ToolCalls = append(oc.Message.ToolCalls, cc.Delta.ToolCalls...)
		if cc.FinishReason != "" {
			oc.FinishReason = cc.FinishReason
		}
	}

	// usage
	if chunk.Usage.CompletionTokens != 0 {
		out.Usage.CompletionTokens = chunk.Usage.CompletionTokens
	}
	if chunk.Usage.PromptTokens != 0 {
		out.Usage.PromptTokens = chunk.Usage.PromptTokens
	}
	if chunk.Usage.TotalTokens != 0 {
		out.Usage.TotalTokens = chunk.Usage.TotalTokens
	}

	// web search
	out.WebSearch = append(out.WebSearch, chunk.WebSearch...)
}

// chatCompletionDecodeStream decode the sse stream of chat completion
func chatCompletionDecodeStream(r io.Reader, fn func(chunk ChatCompletionResponse) error) (err error) {
	br := bufio.NewReader(r)

	for {
		var line []byte

		if line, err = br.ReadBytes('\n'); err != nil {
			if errors.Is(err, io.EOF) {
				err = nil
			}
			break
		}

		line = bytes.TrimSpace(line)

		if len(line) == 0 {
			continue
		}

		if !bytes.HasPrefix(line, chatCompletionStreamPrefix) {
			continue
		}

		data := bytes.TrimSpace(line[len(chatCompletionStreamPrefix):])

		if bytes.Equal(data, chatCompletionStreamDone) {
			break
		}

		if len(data) == 0 {
			continue
		}

		var chunk ChatCompletionResponse
		if err = json.Unmarshal(data, &chunk); err != nil {
			return
		}
		if err = fn(chunk); err != nil {
			return
		}
	}

	return
}

// ChatCompletionStreamService is the service for chat completion stream
type ChatCompletionService struct {
	client *Client

	model       string
	requestID   *string
	doSample    *bool
	temperature *float64
	topP        *float64
	maxTokens   *int
	stop        []string
	toolChoice  *string
	userID      *string
	meta        *ChatCompletionMeta

	messages []any
	tools    []any

	streamHandler ChatCompletionStreamHandler
}

var (
	_ BatchSupport = &ChatCompletionService{}
)

// NewChatCompletionService creates a new ChatCompletionService.
func NewChatCompletionService(client *Client) *ChatCompletionService {
	return &ChatCompletionService{
		client: client,
	}
}

func (s *ChatCompletionService) BatchMethod() string {
	return "POST"
}

func (s *ChatCompletionService) BatchURL() string {
	return BatchEndpointV4ChatCompletions
}

func (s *ChatCompletionService) BatchBody() any {
	return s.buildBody()
}

// SetModel set the model of the chat completion
func (s *ChatCompletionService) SetModel(model string) *ChatCompletionService {
	s.model = model
	return s
}

// SetMeta set the meta of the chat completion, optional
func (s *ChatCompletionService) SetMeta(meta ChatCompletionMeta) *ChatCompletionService {
	s.meta = &meta
	return s
}

// SetRequestID set the request id of the chat completion, optional
func (s *ChatCompletionService) SetRequestID(requestID string) *ChatCompletionService {
	s.requestID = &requestID
	return s
}

// SetTemperature set the temperature of the chat completion, optional
func (s *ChatCompletionService) SetDoSample(doSample bool) *ChatCompletionService {
	s.doSample = &doSample
	return s
}

// SetTemperature set the temperature of the chat completion, optional
func (s *ChatCompletionService) SetTemperature(temperature float64) *ChatCompletionService {
	s.temperature = &temperature
	return s
}

// SetTopP set the top p of the chat completion, optional
func (s *ChatCompletionService) SetTopP(topP float64) *ChatCompletionService {
	s.topP = &topP
	return s
}

// SetMaxTokens set the max tokens of the chat completion, optional
func (s *ChatCompletionService) SetMaxTokens(maxTokens int) *ChatCompletionService {
	s.maxTokens = &maxTokens
	return s
}

// SetStop set the stop of the chat completion, optional
func (s *ChatCompletionService) SetStop(stop ...string) *ChatCompletionService {
	s.stop = stop
	return s
}

// SetToolChoice set the tool choice of the chat completion, optional
func (s *ChatCompletionService) SetToolChoice(toolChoice string) *ChatCompletionService {
	s.toolChoice = &toolChoice
	return s
}

// SetUserID set the user id of the chat completion, optional
func (s *ChatCompletionService) SetUserID(userID string) *ChatCompletionService {
	s.userID = &userID
	return s
}

// SetStreamHandler set the stream handler of the chat completion, optional
// this will enable the stream mode
func (s *ChatCompletionService) SetStreamHandler(handler ChatCompletionStreamHandler) *ChatCompletionService {
	s.streamHandler = handler
	return s
}

// AddMessage add the message to the chat completion
func (s *ChatCompletionService) AddMessage(messages ...ChatCompletionMessageType) *ChatCompletionService {
	for _, message := range messages {
		s.messages = append(s.messages, message)
	}
	return s
}

// AddFunction add the function to the chat completion
func (s *ChatCompletionService) AddTool(tools ...ChatCompletionTool) *ChatCompletionService {
	for _, tool := range tools {
		switch tool := tool.(type) {
		case ChatCompletionToolFunction:
			s.tools = append(s.tools, map[string]any{
				"type":           ToolTypeFunction,
				ToolTypeFunction: tool,
			})
		case ChatCompletionToolRetrieval:
			s.tools = append(s.tools, map[string]any{
				"type":            ToolTypeRetrieval,
				ToolTypeRetrieval: tool,
			})
		case ChatCompletionToolWebSearch:
			s.tools = append(s.tools, map[string]any{
				"type":            ToolTypeWebSearch,
				ToolTypeWebSearch: tool,
			})
		case ChatCompletionToolCodeInterpreter:
			s.tools = append(s.tools, map[string]any{
				"type":                  ToolTypeCodeInterpreter,
				ToolTypeCodeInterpreter: tool,
			})
		case ChatCompletionToolDrawingTool:
			s.tools = append(s.tools, map[string]any{
				"type":              ToolTypeDrawingTool,
				ToolTypeDrawingTool: tool,
			})
		case ChatCompletionToolWebBrowser:
			s.tools = append(s.tools, map[string]any{
				"type":             ToolTypeWebBrowser,
				ToolTypeWebBrowser: tool,
			})
		}
	}
	return s
}

func (s *ChatCompletionService) buildBody() M {
	body := map[string]any{
		"model":    s.model,
		"messages": s.messages,
	}
	if s.requestID != nil {
		body["request_id"] = *s.requestID
	}
	if s.doSample != nil {
		body["do_sample"] = *s.doSample
	}
	if s.temperature != nil {
		body["temperature"] = *s.temperature
	}
	if s.topP != nil {
		body["top_p"] = *s.topP
	}
	if s.maxTokens != nil {
		body["max_tokens"] = *s.maxTokens
	}
	if len(s.stop) != 0 {
		body["stop"] = s.stop
	}
	if len(s.tools) != 0 {
		body["tools"] = s.tools
	}
	if s.toolChoice != nil {
		body["tool_choice"] = *s.toolChoice
	}
	if s.userID != nil {
		body["user_id"] = *s.userID
	}
	if s.meta != nil {
		body["meta"] = s.meta
	}
	return body
}

// Do send the request of the chat completion and return the response
func (s *ChatCompletionService) Do(ctx context.Context) (res ChatCompletionResponse, err error) {
	body := s.buildBody()

	streamHandler := s.streamHandler

	if streamHandler == nil {
		var (
			resp     *resty.Response
			apiError APIErrorResponse
		)
		if resp, err = s.client.request(ctx).SetBody(body).SetResult(&res).SetError(&apiError).Post("chat/completions"); err != nil {
			return
		}
		if resp.IsError() {
			err = apiError
			return
		}
		return
	}

	// stream mode

	body["stream"] = true

	var resp *resty.Response

	if resp, err = s.client.request(ctx).SetBody(body).SetDoNotParseResponse(true).Post("chat/completions"); err != nil {
		return
	}
	defer resp.RawBody().Close()

	if resp.IsError() {
		err = errors.New(resp.Status())
		return
	}

	var choice ChatCompletionChoice

	if err = chatCompletionDecodeStream(resp.RawBody(), func(chunk ChatCompletionResponse) error {
		// reduce the chunk to the response
		chatCompletionReduceResponse(&res, chunk)
		// invoke the stream handler
		return streamHandler(chunk)
	}); err != nil {
		return
	}

	res.Choices = append(res.Choices, choice)

	return
}
