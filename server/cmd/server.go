package cmd

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path"
	"strings"
	"syscall"
	"time"

	"github.com/casimir/freon/api"
	"github.com/casimir/freon/auth"
	"github.com/casimir/freon/buildinfo"
	"github.com/casimir/freon/control"
	"github.com/casimir/freon/ui"
	"github.com/casimir/freon/wallabagproxy"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const defaultPort = ":8080"

func init() {
	serverCmd.Flags().String("port", defaultPort, "http listening port [env: FREON_PORT]")
	viper.BindPFlag("port", serverCmd.Flags().Lookup("port"))
}

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Start the server",
	Run: func(cmd *cobra.Command, args []string) {
		r := gin.New()
		r.Use(gin.Logger(), gin.Recovery())

		r.GET("/api/info", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"appname": "freon",
				"version": buildinfo.Version,
			})
		})

		// TODO CORS config
		corsMiddleware := cors.Default()

		auth.RegisterRoutes(r.Group("/auth", corsMiddleware))
		api.RegisterRoutes(r.Group("/api", auth.TokenAuth()))
		control.RegisterRoutes(r.Group("/control", corsMiddleware, auth.SessionAuth()))
		wallabagproxy.RegisterRoutes(r.Group("/wallabag", auth.TokenAuth()))

		var hasUI bool
		if ui.FS != nil {
			r.Group("/ui").StaticFS("/", ui.FS)
			hasUI = true
			log.Print("/ui/ enabled: embedded")
		} else if uiPath := getStaticsRoot(); uiPath != "" {
			r.Group("/ui").Static("/", getStaticsRoot())
			hasUI = true
			log.Printf("/ui/ enabled: static files from %s", uiPath)
		} else {
			log.Print("/ui/ disabled: no path provided")
		}
		if hasUI {
			r.GET("/", func(c *gin.Context) {
				c.Redirect(http.StatusMovedPermanently, "/ui/")
			})
		}

		port := viper.GetString("PORT")
		if strings.HasPrefix(port, "tcp://") {
			// this is an env clash with k8s
			// FIXME this is a hack
			port = defaultPort
		} else if !strings.HasPrefix(port, ":") {
			port = ":" + port
		}
		startServer(port, r)
	},
}

func startServer(addr string, r http.Handler) {
	srv := &http.Server{
		Addr:    addr,
		Handler: r,
	}

	go func() {
		log.Printf("listening and serving HTTP on %s\n", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("failed to start: %v\n", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("shutting down...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("failed to shutdown: %v\n", err)
	}
	select {
	case <-ctx.Done():
		log.Println("timed-out after 5 seconds")
	default:
	}
	log.Println("shutdown complete")
}

func getStaticsRoot() string {
	uiFilesPath := os.Getenv("FREON_UI_PATH")
	if uiFilesPath == "" {
		return ""
	}
	if path.IsAbs(uiFilesPath) {
		return uiFilesPath
	}
	cwd, err := os.Getwd()
	if err != nil {
		log.Fatalf("could not get current working directory: %v", err)
	}
	return path.Join(cwd, uiFilesPath)
}
