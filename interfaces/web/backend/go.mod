module github.com/syntropy-cc/cooperative-grid/interfaces/web/backend

go 1.21

require (
	github.com/gin-gonic/gin v1.9.1
	github.com/golang-jwt/jwt/v5 v5.0.0
	github.com/gorilla/websocket v1.5.0
	github.com/99designs/gqlgen v0.17.36
	github.com/vektah/gqlparser/v2 v2.5.8
	github.com/syntropy-cc/cooperative-grid/core v0.0.0
	github.com/sirupsen/logrus v1.9.3
)

replace github.com/syntropy-cc/cooperative-grid/core => ../../../core
