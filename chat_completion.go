package zhipu

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/go-resty/resty/v2"
)

const (
	ToolChoiceAuto = "auto"

	FinishReasonStop         = "stop"
	FinishReasonToolCalls    = "tool_calls"
	FinishReasonLength       = "length"
	FinishReasonSensitive    = "sensitive"
	FinishReasonNetworkError = "network_error"

	ToolTypeFunction  = "function"
	ToolTypeWebSearch = "web_search"
	ToolTypeRetrieval = "retrieval"
)

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
	PromptTemplate string `json:"prompt_template"`
}

func (ChatCompletionToolRetrieval) isChatCompletionTool() {}

// ChatCompletionToolWebSearch is the web search for chat completion
type ChatCompletionToolWebSearch struct {
	Enable       bool   `json:"enable"`
	SearchQuery  string `json:"search_query"`
	SearchResult bool   `json:"search_result"`
}

func (ChatCompletionToolWebSearch) isChatCompletionTool() {}

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

// ChatCompletionToolCall is the tool call for chat completion
type ChatCompletionToolCall struct {
	ID       string `json:"id"`
	Type     string `json:"type"`
	Function struct {
		Name      string          `json:"name"`
		Arguments json.RawMessage `json:"arguments"`
	} `json:"function"`
}

// ChatCompletionMessage is the message for chat completion
type ChatCompletionMessage struct {
	Role       string                   `json:"role"`
	Content    string                   `json:"content,omitempty"`
	ToolCalls  []ChatCompletionToolCall `json:"tool_calls,omitempty"`
	ToolCallID string                   `json:"tool_call_id,omitempty"`
}

// ChatCompletionChoice is the choice for chat completion
type ChatCompletionChoice struct {
	Index        int    `json:"index"`
	FinishReason string `json:"finish_reason"`
	// delta is only available in stream mode
	Delta ChatCompletionMessage `json:"delta"`
	// message is only available in non-stream mode
	Message ChatCompletionMessage `json:"message"`
}

// ChatCompletionResponse is the response for chat completion
type ChatCompletionResponse struct {
	ID        string                    `json:"id"`
	Created   int64                     `json:"created"`
	Model     string                    `json:"model"`
	Choices   []ChatCompletionChoice    `json:"choices"`
	Usage     ChatCompletionUsage       `json:"usage"`
	WebSearch []ChatCompletionWebSearch `json:"web_search"`
}

// ChatCompletionStreamHandler is the handler for chat completion stream
type ChatCompletionStreamHandler func(chunk *ChatCompletionResponse) error

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

	messages []ChatCompletionMessage
	tools    []any

	streamHandler ChatCompletionStreamHandler
}

// ChatCompletionRequest is the request for chat completion
func (c *Client) ChatCompletionService(model string) *ChatCompletionService {
	return &ChatCompletionService{
		client: c,
		model:  model,
	}
}

// SetModel set the model of the chat completion
func (s *ChatCompletionService) SetModel(model string) *ChatCompletionService {
	s.model = model
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
func (s *ChatCompletionService) AddMessage(messages ...ChatCompletionMessage) *ChatCompletionService {
	s.messages = append(s.messages, messages...)
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
		}
	}
	return s
}

// Do send the request of the chat completion and return the response
func (s *ChatCompletionService) Do(ctx context.Context) (res ChatCompletionResponse, err error) {
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

	if s.streamHandler == nil {
		var (
			resp     *resty.Response
			apiError APIError
		)
		if resp, err = s.client.R(ctx).SetBody(body).SetResult(&res).SetError(&apiError).Post("chat/completions"); err != nil {
			return
		}
		if resp.IsError() {
			err = apiError
			return
		}
		return
	}

	body["stream"] = true

	//TODO: handle the stream mode
	err = errors.New("stream mode is not implemented yet")

	return
}
