package main

import (
	"fmt"
	"os"
	"path/filepath"
)

func Helper(model int) bool {
	fmt.Println()
	fmt.Println("可用命令:")
	template_example := "     [例]  %s\n"
	template_short := "   %-15s %-40s\n"
	template_long := "   %-30s %-40s\n"
	template_superlong := "   %-30s:\n      %-30s\n"
	switch model {
	case 0:
		fmt.Printf(template_short, "gameasm server", "管理游戏安装包")
		fmt.Printf(template_short, "gameasm client", "解包和还原游戏安装包")
	case 1:
		fmt.Printf(template_long, "gameasm server push <version name> [previous]", "打包一个新的游戏版本，指定版本号")
		fmt.Printf(template_example, "\"gameasm server push 1.0\"")
		fmt.Printf(template_example, "这将会在工作目录的 gameasm/artifacts/ 文件夹下创建一个 1_0.zip 安装包。")
		fmt.Printf(template_example, "生成的是增量安装包，支持对文件的增删改变化进行检测。")
		fmt.Printf(template_example, "指定 previous 参数，程序将从指定的版本对比文件更改。")
	case 2:
		fmt.Printf(template_superlong, "gameasm client install <filename>", "将工作目录的指定文件解压并还原到 gameasm/game/ 文件夹中")
		fmt.Printf(template_example, "例如，将当前目录下的 1_0.zip 文件还原，请输入：")
		fmt.Printf(template_example, "\"gameasm client install 1_0.zip\"")
		fmt.Printf(template_example, "包含的内容将会被恢复至 gameasm/game/ 文件夹内")
	}

	fmt.Println()
	return true
}

func main() {

	fmt.Println()
	fmt.Println("Game Assets Manager [v1]")
	fmt.Println("当前工作目录: " + filepath.Dir(os.Args[0]))

	args := os.Args
	if len(args) > 1 {
		switch args[1] {
		case "server":
			Server(args)
			return
		case "client":
			Client(args)
			return
		}
	}
	Helper(0)
}
