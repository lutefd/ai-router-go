# AI Router Service

A Go-based microservice that routes AI requests to multiple LLM providers (OpenAI, Google Gemini, DeepSeek) with authentication, chat history, and streaming support.

## Features

- Multi-provider AI routing (OpenAI, Google Gemini, DeepSeek)
- Real-time streaming responses
- OAuth2 authentication with Google
- JWT-based authorization
- Chat history management
- Health checks and readiness probes
- MongoDB persistence
- Docker support

## Architecture

This project follows clean architecture principles with clear separation of concerns. For detailed explanations of key architectural decisions, see our [Architecture Decision Records](docs/adr/README.md):

- [Using Chi as HTTP Router](docs/adr/0001-use-go-chi-router.md)
- [MongoDB as Primary Database](docs/adr/0002-mongodb-as-database.md)
- [Strategy Pattern for AI Providers](docs/adr/0003-strategy-pattern-for-ai-providers.md)

Project structure:

```
├── cmd/
│ └── ai-router/ # Application entry point
├── internal/
│ ├── config/ # Configuration management
│ ├── database/ # Database connections
│ ├── handler/ # HTTP handlers
│ ├── middleware/ # HTTP middleware
│ ├── models/ # Domain models
│ ├── repository/ # Data access layer
│ ├── service/ # Business logic
│ ├── strategy/ # AI provider strategy pattern
│ └── server/ # Server setup and routing
├── pkg/ # Public packages
├── docs/
│ └── adr/ # Architecture Decision Records
└── docker-compose.yml
```

## Prerequisites

- Go 1.23 or higher
- MongoDB
- Docker (optional)

## Configuration

The service requires the following environment variables:

```env
SERVER_PORT=8080
MONGODB_URI=mongodb://localhost:27017
MONGODB_DATABASE=ai_router
OPENAI_SK=your_openai_key
DEEPSEEK_SK=your_deepseek_key
GEMINI_SK=your_gemini_key
GOOGLE_CLIENT_ID=your_google_client_id
GOOGLE_CLIENT_SECRET=your_google_client_secret
JWT_SECRET=your_jwt_secret
CLIENT_URL=http://localhost:3000
AUTH_REDIRECT_URL=http://localhost:8080/callback
ANDROID_CLIENT_ID=your_android_client_id
```

## Getting Started

1. Clone the repository:

```bash
git clone https://github.com/lutefd/ai-router-go
```

2. Install dependencies:

```bash
go mod download
```

3. Run with Docker:

```bash
docker-compose up
```

Or run locally:

```bash
go run cmd/ai-router/main.go
```

## API Documentation

### Authentication Endpoints

- `GET /api/v1/auth/google/login` - Initiate Google OAuth login
- `GET /api/v1/auth/google/callback` - OAuth callback handler
- `POST /api/v1/auth/google/native/signin` - Native Google sign-in
- `POST /api/v1/auth/google/refresh` - Refresh JWT token

### AI Endpoints

- `POST /api/v1/ai/generate` - Generate AI responses (requires authentication)

### Chat Endpoints

- `POST /api/v1/chats` - Create new chat
- `GET /api/v1/chats/{id}` - Get chat by ID
- `PUT /api/v1/chats/{id}/title` - Update chat title
- `DELETE /api/v1/chats/{id}` - Delete chat

### Health Endpoints

- `GET /healthz` - Liveness probe
- `GET /readiness` - Readiness probe

## Architecture Decisions

See the [Architecture Decision Records](docs/adr/README.md) for detailed explanations of key technical decisions.

## Testing

Run the tests:

```bash
go test ./...
```
