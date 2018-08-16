// Copyright © 2018 NAME HERE <EMAIL ADDRESS>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/gin-gonic/gin"
	"html/template"
	"time"
	"bytes"
	"strconv"
	"github.com/ThreesomeInc/health_deploy/deploy"
	"net/http"
)

// upCmd represents the up command
var upCmd = &cobra.Command{
	Use:   "up",
	Short: "Bring up the deploy app",
	Long:  `Bring up the deploy app, need to specify the config files using --config=XXXX flag.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Params are: ", "port->", port, ",user->", user, ",repo->", repo, ",deployRoot->", deployRoot)
		startHttpServer(&port, user, repo, deployRoot)
	},
}

func formatAsDate(t time.Time) string {
	year, month, day := t.Date()
	return fmt.Sprintf("%d%02d/%02d", year, month, day)
}

func startHttpServer(port *int, user, repo, deployRoot string) {
	router := gin.Default()
	router.SetFuncMap(template.FuncMap{
		"formatAsDate": formatAsDate,
	})
	//https://github.com/delphinus/gin-assets-sample --see this sample for static bundles.
	router.GET("/deploy", func(c *gin.Context) {
		query := c.Request.URL.Query()
		if query["release"] == nil || len(query["release"]) == 0 || len(query["release"][0]) == 0 {
			c.JSON(http.StatusOK, deploy.Result{Code: 1, Success: false, Message: "release param is not available"})
		}
		response, _ := deploy.Deploy(user, repo, query["release"][0], deployRoot)
		c.JSON(http.StatusOK, response)
	})
	var buffer bytes.Buffer
	buffer.WriteString(":")
	buffer.WriteString(strconv.Itoa(*port))
	router.Run(buffer.String())
}

func init() {
	rootCmd.AddCommand(upCmd)
}
