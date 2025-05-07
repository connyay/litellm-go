package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/connyay/litellm-go/internal/config"
	"github.com/connyay/litellm-go/internal/provider"
	"github.com/connyay/litellm-go/internal/router"
	openai "github.com/sashabaranov/go-openai"
	"golang.org/x/time/rate"
)

func main() {
	var cfgPath string
	var addr string
	flag.StringVar(&cfgPath, "config", "config.yaml", "path to config file")
	flag.StringVar(&addr, "addr", ":4000", "listen address")
	flag.Parse()

	cfg, err := config.Load(cfgPath)
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	rtr := router.New()

	// build providers from config
	for _, mc := range cfg.ModelList {
		key := os.Getenv(mc.APIKeyEnv)
		switch mc.Provider {
		case "openai":
			p := provider.NewOpenAIProvider(mc.ModelName, mc.APIBase, key, "", false)
			rtr.Register(mc.ModelName, p)
		case "azure":
			p := provider.NewOpenAIProvider(mc.ModelName, mc.APIBase, key, mc.APIVersion, true)
			rtr.Register(mc.ModelName, p)
		case "bedrock":
			p, err := provider.NewBedrockProvider(mc.ModelName, mc.Deployment)
			if err != nil {
				log.Printf("failed to init bedrock provider for %s: %v", mc.ModelName, err)
				continue
			}
			rtr.Register(mc.ModelName, p)
		default:
			log.Printf("unknown provider %s for model %s - skipping", mc.Provider, mc.ModelName)
		}
	}

	if rtr.Len() == 0 {
		log.Fatalf("no valid providers configured")
	}

	gin.SetMode(gin.ReleaseMode)
	engine := gin.New()
	engine.Use(gin.Logger())
	engine.Use(gin.Recovery())

	// rate limit middleware
	if cfg.RateLimit != nil && cfg.RateLimit.RequestsPerMinute > 0 {
		lim := rate.NewLimiter(rate.Every(time.Minute/time.Duration(cfg.RateLimit.RequestsPerMinute)), cfg.RateLimit.RequestsPerMinute)
		engine.Use(func(c *gin.Context) {
			if !lim.Allow() {
				c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{"error": "rate limit exceeded"})
				return
			}
			c.Next()
		})
	}

	engine.GET("/healthz", func(c *gin.Context) { c.String(200, "ok") })

	engine.POST("/v1/chat/completions", func(c *gin.Context) {
		var req openai.ChatCompletionRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		prov, ok := rtr.Get(req.Model)
		if !ok {
			c.JSON(http.StatusBadRequest, gin.H{"error": "model not found"})
			return
		}

		resp, err := prov.ChatCompletion(c.Request.Context(), req)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, resp)
	})

	log.Printf("server listening on %s", addr)
	if err := engine.Run(addr); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
