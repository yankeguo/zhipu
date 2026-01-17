# Project Context

## Purpose
A third-party Golang client library for the Zhipu AI Platform (智谱AI). This SDK provides a clean, idiomatic Go interface to interact with Zhipu AI's REST API, covering capabilities such as chat completions (LLMs), embeddings, image generation, video generation, file management, knowledge bases, fine-tuning, and batch processing.

## Tech Stack
- **Language**: Go 1.24.0
- **HTTP Client**: go-resty/resty v2 (for REST API calls)
- **Authentication**: golang-jwt/jwt v5 (for JWT-based API authentication)
- **Testing**: stretchr/testify (for unit and integration tests)
- **Coverage**: Codecov integration for test coverage tracking
- **CI/CD**: GitHub Actions (workflow: go.yml)

## Project Conventions

### Code Style
- **Standard Go conventions**: Follow effective Go guidelines and standard formatting with `gofmt`
- **Naming**:
  - Service types end with `Service` (e.g., `ChatCompletionService`, `EmbeddingService`)
  - Factory methods on client use descriptive names matching the API capability (e.g., `client.ChatCompletion()`, `client.VideoGeneration()`)
  - Builder pattern: setter methods prefixed with `Set` and adder methods with `Add` (e.g., `SetModel()`, `AddMessage()`)
  - Constants use descriptive names with capability prefix (e.g., `RoleUser`, `FinishReasonStop`)
- **Error handling**: Errors are returned explicitly; API errors wrapped with custom error types (`APIError`)
- **Context-aware**: All API calls accept `context.Context` for cancellation and timeout control
- **Options pattern**: Client configuration uses functional options (`ClientOption` type)

### Architecture Patterns
- **Builder/Fluent API**: Service creation and configuration use chainable methods
  - Example: `client.ChatCompletion("model").AddMessage(...).SetStreamHandler(...).Do(ctx)`
- **Service-oriented design**: Each API capability is encapsulated in its own service type
- **Client-based factory**: Central `Client` type provides factory methods for all services
- **JWT authentication**: Client generates JWT tokens automatically using API key components (keyID.keySecret)
- **Streaming support**: Stream responses handled via callback handlers (e.g., `SetStreamHandler`)
- **Batch processing**: Specialized batch writer and reader for JSONL batch files

### Testing Strategy
- **Integration tests**: Tests use real API calls (require `ZHIPUAI_API_KEY` environment variable)
- **Test framework**: stretchr/testify with `require` assertions
- **Test coverage**: Tracked via Codecov
- **Test naming**: `Test<ServiceName>` pattern (e.g., `TestChatCompletionService`)
- **Test data**: Stored in `testdata/` directory for file upload tests
- **No mocking**: Tests verify actual API behavior (not unit tests with mocks)

### Git Workflow
- **Commit conventions**: Use conventional commits with scope (e.g., `feat(chat): add streaming support`, `fix(client): handle API errors`)
- **Badges**: Maintain Go reference, build status, and coverage badges in README
- **Documentation**: Bilingual documentation (English README.md, Chinese README.zh.md)

## Domain Context
- **Zhipu AI Platform**: Chinese AI platform providing large language models (GLM series), embeddings, image/video generation
- **API Authentication**: Uses JWT tokens signed with HMAC-SHA256, requiring API key split into ID and secret
- **Model naming**: 
  - Chat: `glm-4-flash`, `charglm-3`, `GLM-4-AllTools`
  - Embedding: `embedding-v2`
  - Image: `cogview-3`
  - Video: `cogvideox`
- **Tool support**: GLM-4-AllTools supports code interpreter, web browser, and drawing tools
- **Async operations**: Video generation returns task ID; requires polling via `AsyncResult` service
- **File purposes**: Files categorized by purpose (`retrieval` for knowledge bases, `fine-tune` for training)

## Important Constraints
- **API key format**: Must be in format `keyID.keySecret` (validated during client creation)
- **Environment variables**: 
  - `ZHIPUAI_API_KEY`: API key (required if not passed explicitly)
  - `ZHIPUAI_BASE_URL`: Custom base URL (defaults to `https://open.bigmodel.cn/api/paas/v4`)
  - `ZHIPUAI_DEBUG`: Enable debug mode for HTTP request/response logging
- **JWT expiration**: Tokens expire after 7 days (automatically regenerated per request)
- **Streaming caveat**: Package combines stream chunks into final result mimicking non-streaming API
- **Batch format**: Batch files use JSONL format with custom ID field

## External Dependencies
- **Zhipu AI REST API**: `https://open.bigmodel.cn/api/paas/v4` (configurable)
- **go-resty**: Wrapper around Go's net/http for fluent HTTP client API
- **golang-jwt**: JWT generation and signing for API authentication
- **GitHub Actions**: CI/CD pipeline for automated testing
- **Codecov**: Test coverage reporting and tracking
- **pkg.go.dev**: Go package documentation hosting
