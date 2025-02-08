# 1. Use Chi as the HTTP Router

## Status

Accepted

## Context

We needed a lightweight, flexible HTTP router for our Go service that supports:

- Middleware chaining for cross-cutting concerns
- URL parameter extraction
- Route grouping for API versioning
- High performance for streaming responses
- Easy testing capabilities
- Compatibility with standard library

We evaluated several options including:

- Gin
- Echo
- Gorilla Mux
- Standard net/http
- Chi

## Decision

We chose to use Chi (github.com/go-chi/chi) as our HTTP router for the following reasons:

1. **Middleware Architecture**

   - Supports both global and route-specific middleware
   - Clean middleware chaining syntax
   - Easy to implement custom middleware

2. **Routing Features**

   - URL parameters with type-safe extraction
   - Flexible pattern matching
   - Route groups for clean API organization

3. **Standard Library Compatibility**

   - Uses standard `http.Handler` interface
   - No custom context types
   - Easy integration with existing Go packages

4. **Performance**
   - Minimal allocations
   - Fast routing with radix tree
   - Efficient middleware chain

## Consequences

### Positive

- Lightweight and fast with minimal dependencies
- Middleware support with an elegant API
- Compatible with net/http interface
- Good community support and active maintenance
- Easy to test using standard Go testing tools
- Built-in support for URL parameters and routing groups
- Clean routing syntax for REST APIs
- Easy to implement custom middleware
- Good documentation and examples

### Negative

- Less feature-rich compared to frameworks like Gin or Echo
- Requires more manual setup for certain features
- No built-in template rendering
- No automatic OpenAPI documentation generation
- Manual validation implementation required
- No built-in binding of request bodies to structs

### Mitigations

1. For request/response binding:

   ```go
   func (h *Handler) CreateChat(w http.ResponseWriter, r *http.Request) {
       var chat models.Chat
       if err := json.NewDecoder(r.Body).Decode(&chat); err != nil {
           http.Error(w, "Invalid request body", http.StatusBadRequest)
           return
       }
   }
   ```

2. For validation:

   ```go
   func validateChat(chat *models.Chat) error {
       if chat.Title == "" {
           return fmt.Errorf("chat title is required")
       }
       return nil
   }
   ```

3. For middleware:
   ```go
   func (m *AuthMiddleware) RequireAuth(next http.Handler) http.Handler {
       return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
           // Authentication logic
           ctx := context.WithValue(r.Context(), UserContextKey, claims)
           next.ServeHTTP(w, r.WithContext(ctx))
       })
   }
   ```
