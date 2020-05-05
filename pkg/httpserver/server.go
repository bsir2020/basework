package httpserver

import (
	"github.com/bsir2020/basework/pkg/filter"
	"github.com/gin-gonic/gin"
	"net/http"
)

type Server struct {
	engine   *gin.Engine
	pathGet  map[string]gin.HandlerFunc
	pathPost map[string]gin.HandlerFunc
}

func New() *Server {
	e := gin.Default()
	return &Server{
		engine: e,
	}
}

func (s *Server) SetGetRouter(route string, handle func(*gin.Context)) {
	s.pathGet[route] = handle
}

func (s *Server) SetPostRouter(route string, handle func(*gin.Context)) {
	s.pathPost[route] = handle
}

func (s *Server) assem() {
	filter := filter.Filter{}

	authorized := s.engine.Group("/")
	authorized.Use(filter.Checkauth())
	{
		for key, handle := range s.pathGet {
			s.engine.GET(key, handle)
		}

		for key, handle := range s.pathPost {
			s.engine.POST(key, handle)
		}
	}
}

func (s *Server) Run(ip string, port int64) {
	s.assem()
	s.engine.Use(s.cross())
	s.engine.Use(gin.Logger())
	s.engine.Use(gin.Recovery())
	s.engine.NoRoute(s.noResponse)
	s.engine.Run(ip + ":" + string(port))
}

func (s *Server) cross() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method

		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Headers", "Content-Type,AccessToken,X-CSRF-Token, Authorization, Token")
		c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, PATCH, DELETE")
		c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Content-Type")
		c.Header("Access-Control-Allow-Credentials", "true")

		if method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
		}

		// 处理请求
		c.Next()
	}
}

func (s *Server) noResponse(c *gin.Context) {
	c.JSON(http.StatusNotFound, gin.H{
		"status": 404,
		"error":  "404, page not exists!",
	})
}
