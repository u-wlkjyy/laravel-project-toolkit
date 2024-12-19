package main

import (
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"syscall"

	"github.com/spf13/cobra"
)

// 全局变量，用于存储当前运行的composer进程
var currentCmd *exec.Cmd

var rootCmd = &cobra.Command{
	Use:   "laravel-project",
	Short: "Laravel项目创建工具",
	Long:  `Laravel项目创建工具 - 用于快速创建指定版本的Laravel项目`,
}

var createCmd = &cobra.Command{
	Use:   "create [项目名称] [版本号]",
	Short: "创建新的Laravel项目",
	Long:  `创建新的Laravel项目，可以指定项目名称和版本号。如果不指定项目名称，将使用当前目录名称；如果不指定版本号，将使用最新版本。`,
	Run: func(cmd *cobra.Command, args []string) {
		projectName := ""
		version := ""

		// 处理项目名称
		if len(args) > 0 {
			projectName = args[0]
		} else {
			currentDir, err := os.Getwd()
			if err != nil {
				fmt.Println("无法获取当前目录名称:", err)
				return
			}
			projectName = filepath.Base(currentDir)
		}

		// 处理版本号
		if len(args) > 1 {
			version = args[1]
		}

		createLaravelProject(projectName, version)
	},
}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "列出所有可用的Laravel版本",
	Run: func(cmd *cobra.Command, args []string) {
		listLaravelVersions()
	},
}

func createLaravelProject(name string, version string) {
	var command string
	if version == "" {
		command = fmt.Sprintf("composer create-project laravel/laravel %s --prefer-dist", name)
	} else {
		command = fmt.Sprintf("composer create-project laravel/laravel %s --prefer-dist %s.*", name, version)
	}

	fmt.Printf("正在创建Laravel项目: %s\n", name)
	if version != "" {
		fmt.Printf("指定版本: %s\n", version)
	} else {
		fmt.Println("使用最新版本")
	}

	currentCmd = exec.Command("powershell", "/c", command)
	currentCmd.Stdout = os.Stdout
	currentCmd.Stderr = os.Stderr

	// 创建信号通道
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	// 在新的goroutine中处理信号
	go func() {
		<-sigChan
		fmt.Println("\n收到终止信号，正在停止项目创建...")
		if currentCmd != nil && currentCmd.Process != nil {
			// 在Windows上终止进程树
			exec.Command("taskkill", "/F", "/T", "/PID", strconv.Itoa(currentCmd.Process.Pid)).Run()
			fmt.Println("已终止所有相关进程")
			os.Exit(1)
		}
	}()

	err := currentCmd.Run()
	if err != nil {
		if err.Error() != "exit status 1" { // 忽略因为手动终止导致的错误
			fmt.Printf("创建项目时发生错误: %v\n", err)
		}
		return
	}

	fmt.Printf("\nLaravel项目 '%s' 创建成功！\n", name)
}

func listLaravelVersions() {
	cmd := exec.Command("powershell", "/c", "composer show laravel/laravel --all")
	output, err := cmd.Output()
	if err != nil {
		fmt.Println("获取版本列表失败:", err)
		return
	}

	versions := string(output)
	if strings.Contains(versions, "versions") {
		versionParts := strings.Split(versions, "versions")
		if len(versionParts) > 1 {
			versionList := strings.TrimSpace(versionParts[1])
			// 提取版本号部分（到下一个关键词之前）
			if idx := strings.Index(versionList, "type"); idx != -1 {
				versionList = versionList[:idx]
			}
			// 分割版本号
			versions := strings.Split(strings.ReplaceAll(versionList, " ", ""), ",")
			
			// 对版本号进行分类
			majorVersions := make(map[string][]string)
			for _, v := range versions {
				if strings.HasPrefix(v, "v") {
					majorVersion := strings.Split(v[1:], ".")[0]
					majorVersions[majorVersion] = append(majorVersions[majorVersion], v)
				}
			}

			// 打印表头
			fmt.Println("\n可用的Laravel版本:")
			fmt.Println(strings.Repeat("-", 100))
			fmt.Printf("| %-10s | %-80s |\n", "主版本", "具体版本")
			fmt.Println(strings.Repeat("-", 100))

			// 按主版本号排序
			var keys []int
			for k := range majorVersions {
				if num, err := strconv.Atoi(k); err == nil {
					keys = append(keys, num)
				}
			}
			sort.Sort(sort.Reverse(sort.IntSlice(keys)))

			// 打印每个主版本的版本列表
			for _, k := range keys {
				versions := majorVersions[strconv.Itoa(k)]
				// 每行最多显示5个版本
				for i := 0; i < len(versions); i += 5 {
					end := i + 5
					if end > len(versions) {
						end = len(versions)
					}
					if i == 0 {
						fmt.Printf("| Laravel %-3d | %-80s |\n", k, strings.Join(versions[i:end], ", "))
					} else {
						fmt.Printf("| %-10s | %-80s |\n", "", strings.Join(versions[i:end], ", "))
					}
				}
				fmt.Println(strings.Repeat("-", 100))
			}
		}
	} else {
		fmt.Println("无法解析版本信息")
	}
}

func main() {
	rootCmd.AddCommand(createCmd)
	rootCmd.AddCommand(listCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}