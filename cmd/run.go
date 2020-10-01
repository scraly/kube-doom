// Copyright © 2020 NAME HERE <EMAIL ADDRESS>
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
	"log"
	"net"
	"strconv"
	"strings"
	"time"
	"context"

	"github.com/scraly/kube-doom/kubernetes"
	"github.com/scraly/kube-doom/utils"
	"github.com/spf13/cobra"
)

var mode, namespace string

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Run comand launch kubedoom",
	Long:  `Run command launch kubedoom program in order to kill pods or namespaces through Doom game`,
	Run: func(cmd *cobra.Command, args []string) {
		options := fmt.Sprintf("run command called with options = mode:%s - namespace:%s", mode, namespace)
		fmt.Println(options)

		//TODO: depending on the mode, retrieves Pods or namespaces
		//TODO: For mode:pods, retrieve pods in all namespaces if namespace="" or in a given namespace
		resources := getResources(mode, namespace)

		startDoom(resources)
	},
}

//getResources alllows you to retrieve pods or namespaces
func getResources(mode, namespace string) []string {
	ctx := context.Background()

	//TODO: continue ...
	resources := []string{"pod1", "pod2"}

	// var getKubeCmd string
	if mode == "namespaces" {
		// getKubeCmd = fmt.Sprintf("kubectl get %s", mode)
	} else if namespace == "" && mode == "pods" {
		clt := kubernetes.GetClient()
		resources := kubernetes.GetPods(ctx, clt, namespace)
		fmt.Println(resources)
	} else {
		// getKubeCmd = fmt.Sprintf("kubectl get %s -ns %s", mode, namespace)
	}

	// getKubeCmd = fmt.Sprintf(getKubeCmd, `-o go-template --template="{{range .items}}{{.metadata.name}} {{end}}"`)
	//"-o", "go-template", "--template={{range .items}}{{.metadata.name}} {{end}}"

	// fmt.Println(getKubeCmd)
	//kubectl get <mode>

	//TOOO: use Kube Go lib in order to retrieve pods or namespaces to Kill
	

	resources := []string{"pod1", "pod2"}

	return resources
}

func startDoom(resources []string) {
	listener, err := net.Listen("unix", "/dockerdoom.socket")
	if err != nil {
		log.Fatalf("Could not create socket file")
	}

	fmt.Println("Create virtual display")
	utils.LaunchCmd("/usr/bin/Xvfb :99 -ac -screen 0 640x480x24")
	time.Sleep(time.Duration(2) * time.Second)
	utils.LaunchCmd("x11vnc -geometry 640x480 -forever -usepw -display :99")
	fmt.Println("You can now connect to it with a VNC viewer at port 5900")

	fmt.Println("Trying to start DOOM ...")
	utils.LaunchCmd("/usr/bin/env DISPLAY=:99 /usr/local/games/psdoom -warp -E1M1")

	socketLoop(listener, resources)
}

//TODO: to refactor !
func socketLoop(listener net.Listener, resources []string) {
	for true {
		conn, err := listener.Accept()
		if err != nil {
			panic(err)
		}
		stop := false
		for !stop {
			bytes := make([]byte, 40960)
			n, err := conn.Read(bytes)
			if err != nil {
				stop = true
			}
			bytes = bytes[0:n]
			strbytes := strings.TrimSpace(string(bytes))
			// entities := mode.getResources()
			if strbytes == "list" {
				for _, entity := range resources {
					padding := strings.Repeat("\n", 255-len(entity))
					_, err = conn.Write([]byte(entity + padding))
					if err != nil {
						log.Fatal("Could not write to socker file")
					}
				}
				conn.Close()
				stop = true
			} else if strings.HasPrefix(strbytes, "kill ") {
				parts := strings.Split(strbytes, " ")
				killhash, err := strconv.ParseInt(parts[1], 10, 32)
				if err != nil {
					log.Fatal("Could not parse kill hash")
				}
				for _, resource := range resources {
					if hash(resource) == int32(killhash) {
						//TODO: Kill resource
						// mode.deleteEntity(entity)
						break
					}
				}
				conn.Close()
				stop = true
			}
		}
	}
}

func hash(input string) int32 {
	var hash int32
	hash = 5381
	for _, char := range input {
		hash = ((hash << 5) + hash + int32(char))
	}
	if hash < 0 {
		hash = 0 - hash
	}
	return hash
}

func init() {

	runCmd.PersistentFlags().StringVarP(&mode, "mode", "m", "pods", "Mode: pods or namespaces are allowed")
	runCmd.PersistentFlags().StringVarP(&namespace, "namespace", "n", "", "In which namespace do you want to kill Pods")

	rootCmd.AddCommand(runCmd)
}
