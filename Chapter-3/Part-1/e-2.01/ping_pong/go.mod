module ping_pong

go 1.25.1

replace common => ../common

require (
	common v0.0.0-0
	github.com/go-chi/chi/v5 v5.2.3
	github.com/go-chi/cors v1.2.2
)

require github.com/joho/godotenv v1.5.1 // indirect
