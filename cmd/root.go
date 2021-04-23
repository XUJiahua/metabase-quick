/*
Copyright © 2021 xujiahua <littleguner@gmail.com>

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
	"fmt"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/phayes/freeport"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/xujiahua/metabase-quick/pkg/metabase"
	"github.com/xujiahua/metabase-quick/pkg/sqlclient"
	"github.com/xujiahua/metabase-quick/pkg/sqldb"
	"os"
	"os/exec"
	"os/signal"
	"runtime"
	"syscall"
)

var cfgFile string
var webServerPort int
var sqlServerPort int
var hasHeader bool
var verbose bool
var dev bool

const defaultDBName = "default"
const defaultDBUser = "root"
const defaultDBPass = ""

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "metabase-quick",
	Short: "visualize local csv via metabase quickly",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		if verbose {
			logrus.SetLevel(logrus.DebugLevel)
		} else {
			logrus.SetLevel(logrus.InfoLevel)
		}

		if len(args) == 0 {
			fmt.Println("expect dataset filename")
			return
		}

		var err error
		if sqlServerPort == 0 {
			sqlServerPort, err = freeport.GetFreePort()
			handleErr(err)
		}
		sqlServerAddr := fmt.Sprintf("127.0.0.1:%d", sqlServerPort)
		if webServerPort == 0 {
			webServerPort, err = freeport.GetFreePort()
			handleErr(err)
		}
		webServerAddr := fmt.Sprintf("127.0.0.1:%d", webServerPort)

		shutdownSignal := make(chan struct{}, 1)
		done := make(chan struct{}, 2)

		// start built in sql server
		sqlServer, err := sqldb.New(sqlServerAddr, defaultDBUser, defaultDBPass, defaultDBName)
		handleErr(err)

		err = sqlServer.Import(args, hasHeader)
		handleErr(err)

		go func() {
			fmt.Println("You can visit database engine via below command:")
			fmt.Printf("mysql --host=127.0.0.1 --port=%d %s -u %s\n", sqlServerPort, defaultDBName, defaultDBUser)
			fmt.Println()

			go func() {
				<-shutdownSignal
				logrus.Debug("shutting down sql server...")
				_ = sqlServer.Close()
			}()

			logrus.Debug("starting sql server")
			handleErr(sqlServer.Start())
			logrus.Debug("sql server is closed")
			done <- struct{}{}
		}()

		// create sql client
		client, err := sqlclient.New(sqlServerAddr, defaultDBUser, defaultDBPass, defaultDBName)
		handleErr(err)

		// start metabase mock server
		webServer, err := metabase.New(client, webServerAddr, dev)
		handleErr(err)

		go func() {
			fmt.Printf("open http://localhost:%d\n", webServerPort)
			fmt.Println()

			go func() {
				<-shutdownSignal
				logrus.Debug("shutting down web server...")
				_ = webServer.Close()
			}()

			logrus.Debug("starting web server")
			handleErr(webServer.Start())
			logrus.Debug("web server is closed")
			done <- struct{}{}
		}()

		openBrowser(fmt.Sprintf("http://localhost:%d", webServerPort))

		gracefulShutdown(shutdownSignal, done)
	},
}

func gracefulShutdown(shutdownSignal chan struct{}, done chan struct{}) {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	sig := <-sigs
	fmt.Println()
	fmt.Println(sig)
	close(shutdownSignal)

	// wait web server finish
	<-done
	<-done
}

// https://gist.github.com/hyg/9c4afcd91fe24316cbf0
func openBrowser(url string) {
	var err error

	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	default:
		err = fmt.Errorf("unsupported platform")
	}
	if err != nil {
		logrus.Fatal(err)
	}

}

func handleErr(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(2)
	}
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().IntVarP(&sqlServerPort, "sqlServerPort", "s", 0, "will use free open port if you don't specify")
	rootCmd.PersistentFlags().IntVarP(&webServerPort, "webServerPort", "w", 0, "will use free open port if you don't specify")
	rootCmd.PersistentFlags().BoolVarP(&hasHeader, "hasHeader", "i", true, "indicate if csv has header row")
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "show verbose logs")
	rootCmd.PersistentFlags().BoolVarP(&dev, "dev", "d", false, "dev mode, reverse proxy metabase:3000")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	//rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".metabase-quick" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".metabase-quick")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}
