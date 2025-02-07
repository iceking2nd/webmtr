/*
Copyright © 2025 Daniel Wu<wxc@wxccs.org>

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package cmd

import (
	"context"
	"errors"
	"fmt"
	"github.com/iceking2nd/webmtr/app/models"
	"github.com/iceking2nd/webmtr/app/routes"
	"github.com/iceking2nd/webmtr/global"
	"github.com/sirupsen/logrus"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
	"github.com/szuecs/gin-glog"
	ginlogrus "github.com/toorop/gin-logrus"
)

var (
	cfgFile       string
	logFile       string
	listenAddress string
	listenPort    int
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:     "webmtr",
	Version: global.Version,
	Short:   "A tool that provides REST - style APIs for the MTR",
	/*Long: `A longer description that spans multiple lines and likely contains
	examples and usage of using your application. For example:

	Cobra is a CLI library for Go that empowers applications.
	This application is a tool to generate the needed files
	to quickly create a Cobra application.`,*/
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) {
		if global.LogLevel < 5 {
			gin.SetMode(gin.ReleaseMode)
		}
		apiEngine := gin.New()
		corsConfig := cors.Config{
			AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "HEAD"},
			AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "Authorization"},
			AllowCredentials: true,
			MaxAge:           12 * time.Hour,
			AllowAllOrigins:  true,
		}
		apiEngine.Use(cors.New(corsConfig))
		apiServer := &http.Server{
			Addr:         fmt.Sprintf("%s:%d", listenAddress, listenPort),
			Handler:      apiEngine,
			ReadTimeout:  120 * time.Second,
			WriteTimeout: 120 * time.Second,
		}

		apiEngine.NoRoute(func(c *gin.Context) {
			c.JSON(http.StatusNotFound, gin.H{"error": "404 page not found"})
		})

		root := apiEngine.Group("/")
		root.Use(ginglog.Logger(3 * time.Second))
		root.Use(ginlogrus.Logger(global.Log), gin.Recovery())
		routes.SetupRouter(root)

		ln, err := net.Listen("tcp", fmt.Sprintf("%s:%d", listenAddress, listenPort))
		if err != nil {
			log.Fatalf("create listener error: %s\n", err.Error())
		}
		fmt.Printf("listening on %s ...\n", ln.Addr().String())

		go func() {
			if err := apiServer.Serve(ln); err != nil && !errors.Is(err, http.ErrServerClosed) {
				log.Fatalf("HTTP server listen: %s\n", err.Error())
			}
		}()
		signalChan := make(chan os.Signal)
		signal.Notify(signalChan, os.Interrupt, syscall.SIGHUP, syscall.SIGTERM, syscall.SIGKILL, syscall.SIGINT)
		ticker := time.NewTicker(time.Millisecond)
		for {
			select {
			case sig := <-signalChan:
				log.Println("Get Signal:", sig)
				log.Println("Shutdown Server ...")
				ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
				defer cancel()
				if err := apiServer.Shutdown(ctx); err != nil {
					log.Fatal("Closing web service error: ", err)
				}
				log.Println("Server exiting")
				os.Exit(0)
			case <-ticker.C:
				//do sth
			}
		}

	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.webmtr.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	//rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	rootCmd.Flags().IntVarP(&models.DefParams.COUNT, "count", "c", 5, "set the number of pings sent")
	rootCmd.Flags().DurationVarP(&models.DefParams.TIMEOUT, "timeout", "t", 800*time.Millisecond, "ICMP echo request timeout")
	rootCmd.Flags().DurationVarP(&models.DefParams.INTERVAL, "interval", "i", 100*time.Millisecond, "ICMP echo request interval")
	rootCmd.Flags().DurationVar(&models.DefParams.HOP_SLEEP, "hop-sleep", time.Nanosecond, "wait time between pinging next hop")
	rootCmd.Flags().IntVarP(&models.DefParams.MAX_HOPS, "max-hops", "m", 64, "maximum number of hops")
	rootCmd.Flags().IntVarP(&models.DefParams.MAX_UNKNOWN_HOPS, "max-unknown-hops", "U", 10, "maximum unknown host")
	rootCmd.Flags().IntVar(&models.DefParams.RING_BUFFER_SIZE, "buffer-size", 50, "cached packet buffer size")
	rootCmd.Flags().BoolVarP(&models.DefParams.JsonFmt, "json", "j", false, "output as JSON")
	rootCmd.Flags().BoolVarP(&models.DefParams.PTR_LOOKUP, "no-dns", "n", false, "disable DNS lookup")
	rootCmd.Flags().StringVarP(&models.DefParams.SrcAddr, "address", "a", "0.0.0.0", "bind the outgoing socket to ADDRESS")

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.configs/webhooks-mtr/webhooks-mtr.yaml)")
	rootCmd.PersistentFlags().StringVar(&logFile, "log-file", "", "logging file")
	rootCmd.PersistentFlags().Uint32Var(&global.LogLevel, "log-level", 3, "log level (0 - 6, 3 = warn , 5 = debug)")
	rootCmd.PersistentFlags().StringVarP(&listenAddress, "listen", "l", "127.0.0.1", "listen address (127.0.0.1 as default)")
	rootCmd.PersistentFlags().IntVarP(&listenPort, "port", "p", 0, "listen port (random as default)")

	rootCmd.SetVersionTemplate(fmt.Sprintf(`{{with .Name}}{{printf "%%s version information: " .}}{{end}}
   {{printf "Version:    %%s" .Version}}
   Build Time:		%s
   Git Revision:	%s
   Go version:		%s
   OS/Arch:			%s/%s
`, global.BuildTime, global.GitCommit, runtime.Version(), runtime.GOOS, runtime.GOARCH))

	if models.DefParams.SrcAddr == "" || models.DefParams.SrcAddr == "0.0.0.0" {
		models.DefParams.SrcAddr = GetOutboundIPWithDestination(net.ParseIP("8.8.8.8")).String()
	}
}

func initConfig() {
	global.Log = logrus.New()
	var logWriter io.Writer
	if logFile == "" {
		logWriter = os.Stdout
	} else {
		logFileHandle, err := os.OpenFile(logFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
		if err != nil {
			panic(err.Error())
		}
		logWriter = io.MultiWriter(os.Stdout, logFileHandle)
	}
	global.Log.SetOutput(logWriter)
	global.Log.SetLevel(logrus.Level(global.LogLevel))
}

func GetOutboundIPWithDestination(destination net.IP) net.IP {

	// 尝试建立 UDP 连接
	conn, err := net.Dial("udp", fmt.Sprintf("%s:80", destination.String()))
	if err != nil {
		log.Printf("[ERROR] Failed to dial: %v", err)
		return nil
	}
	defer func() {
		if closeErr := conn.Close(); closeErr != nil {
			log.Printf("[WARNING] Failed to close connection: %v", closeErr)
		}
	}()

	// 获取本地地址
	localAddr := conn.LocalAddr()
	if udpAddr, ok := localAddr.(*net.UDPAddr); ok {
		return udpAddr.IP
	} else {
		log.Printf("[ERROR] Local address is not of type *net.UDPAddr: %T", localAddr)
		return nil
	}
}
