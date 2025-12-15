## API

### Handler template

```go
type Handler struct {}

func (h *Handler) Handle(c *gin.Context) {}

func (h *Handler) Method() string {
	return http.MethodGet
}

func (h *Handler) Path() string {
	return "/handler"
}

func (h *Handler) Middleware() []string {
	return nil
}

func New() *Handler {
	return &Handler{}
}
```

### Handler package naming

```
path+method+handler
```

Например

```
user_me_statistics_get_handler
```

### Migration

Call from backend root dir

```bash
migrate create -ext sql -dir src/migrations create_users_table
```
