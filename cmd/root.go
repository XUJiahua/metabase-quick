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
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/xujiahua/metabase-quick/pkg/sqldb"
	"net/http"
	"net/http/httputil"
	"os"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
)

var cfgFile string
var sqlServerAddr string
var hasHeader bool
var debug bool

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "metabase-quick",
	Short: "visualize local csv via metabase quickly",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		if debug {
			logrus.SetLevel(logrus.DebugLevel)
		}

		if len(args) == 0 {
			fmt.Println("expect dataset filename")
			return
		}

		// metabase server
		router := gin.Default()
		apiGroup := router.Group("/api")
		const contentType = "application/json;charset=utf-8"
		{
			apiGroup.GET("/database", func(c *gin.Context) {
				res := `
[
    {
        "description": null,
        "features": [
            "basic-aggregations",
            "standard-deviation-aggregations",
            "expression-aggregations",
            "foreign-keys",
            "right-join",
            "left-join",
            "native-parameters",
            "nested-queries",
            "expressions",
            "case-sensitivity-string-filter-options",
            "binning",
            "inner-join",
            "advanced-math-expressions"
        ],
        "cache_field_values_schedule": "0 0 9 * * ? *",
        "timezone": "UTC",
        "auto_run_queries": true,
        "metadata_sync_schedule": "0 38 * * * ? *",
        "name": "Sample Dataset",
        "caveats": null,
        "is_full_sync": true,
        "updated_at": "2021-04-21T02:43:59.167",
        "native_permissions": "write",
        "details": {
            "db": "/Users/jiahua/opensource/metabase/resources/sample-dataset.db;USER=GUEST;PASSWORD=guest"
        },
        "is_sample": true,
        "id": 1,
        "is_on_demand": false,
        "options": null,
        "engine": "h2",
        "refingerprint": null,
        "created_at": "2021-04-20T05:02:06.893",
        "points_of_interest": null
    }
]`
				c.Data(200, contentType, []byte(res))
			})
			apiGroup.GET("/database/1/schemas", func(c *gin.Context) {
				res := `
["PUBLIC"]
`
				c.Data(200, contentType, []byte(res))
			})

			apiGroup.GET("/database/1/schema/PUBLIC", func(c *gin.Context) {
				res := `
[
    {
        "description": "This is a confirmed order for a product from a user.",
        "entity_type": "entity/TransactionTable",
        "schema": "PUBLIC",
        "show_in_getting_started": false,
        "name": "ORDERS",
        "caveats": null,
        "updated_at": "2021-04-21T07:38:00.493",
        "entity_name": null,
        "active": true,
        "id": 2,
        "db_id": 1,
        "visibility_type": null,
        "field_order": "database",
        "display_name": "Orders",
        "created_at": "2021-04-20T05:02:07.618",
        "points_of_interest": null
    },
    {
        "description": "This is our product catalog. It includes all products ever sold by the Sample Company.",
        "entity_type": "entity/ProductTable",
        "schema": "PUBLIC",
        "show_in_getting_started": false,
        "name": "PRODUCTS",
        "caveats": null,
        "updated_at": "2021-04-21T07:38:00.551",
        "entity_name": null,
        "active": true,
        "id": 1,
        "db_id": 1,
        "visibility_type": null,
        "field_order": "database",
        "display_name": "Products",
        "created_at": "2021-04-20T05:02:07.606",
        "points_of_interest": null
    }
]
`
				c.Data(200, contentType, []byte(res))
			})

			// TODO: ??? 有必要???
			apiGroup.GET("/table/:id/query_metadata", func(c *gin.Context) {
				id := c.Param("id")
				logrus.Infof(id)
			})
		}

		// reverse proxy to metabase server
		router.NoRoute(ReverseProxy())
		port, _ := os.LookupEnv("PORT")
		router.Run(fmt.Sprintf(":%s", port))

		// sql server
		s, err := sqldb.New(sqlServerAddr)
		if err != nil {
			fmt.Println(err)
			return
		}
		for _, filename := range args {
			err = s.ImportTable(filename, hasHeader)
			if err != nil {
				fmt.Println(err)
				return
			}
		}
		err = s.Start()
		fmt.Println(err)
		return
	},
}

func ReverseProxy() gin.HandlerFunc {
	target := "localhost:3000"

	return func(c *gin.Context) {
		director := func(req *http.Request) {
			req.URL.Scheme = "http"
			req.URL.Host = target
			//req.Header["my-header"] = []string{r.Header.Get("my-header")}
			//// Golang camelcases headers
			//delete(req.Header, "My-Header")
		}
		proxy := &httputil.ReverseProxy{Director: director}
		proxy.ServeHTTP(c.Writer, c.Request)
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

	rootCmd.PersistentFlags().StringVarP(&sqlServerAddr, "sqlServerAddr", "s", "localhost:3306", "the address sql server will listen")
	rootCmd.PersistentFlags().BoolVarP(&hasHeader, "hasHeader", "i", true, "indicate if csv has header row")
	rootCmd.PersistentFlags().BoolVarP(&debug, "debug", "d", false, "debug mode")

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
