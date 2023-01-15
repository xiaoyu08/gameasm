package main

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"

	cp "github.com/otiai10/copy"
)

type GamePackageIndex struct {
	hash string
	name string
}

func Server(args []string) {
	CheckConfig()
	if len(args) > 2 {
		switch args[2] {
		case "push":
			Push(args)
			return
		}
	}
	Helper(1)
}

func CheckConfig() {
	data, err := os.ReadFile("./config.dat")
	if err != nil || data == nil {
		fmt.Println()
		fmt.Println("打扰一下, 无法找到配置文件, 请回答下列问题:")

		input := ConversationQuestion("游戏资源原始目录")

		d1 := []byte(input)
		os.WriteFile("config.dat", d1, 0644)
		fmt.Println("谢谢，已创建配置文件，正在继续下一步...")
	}
}

func Push(args []string) {
	fmt.Println()
	var version_string string = ""
	reg := regexp.MustCompile(`[^a-zA-Z0-9.]+`)
	if len(args) != 4 {
		fmt.Println("一些必要的参数尚未被满足, 请回答下列问题:")
		version_string = ConversationQuestionRegexp("游戏版本字符串", reg, "字母数字和点")
	} else {
		version_string = reg.ReplaceAllString(args[3], "")
	}
	PackageFiles(version_string)
	fmt.Println()
}

func getmd5(path string) string {
	file, err := os.Open(path)

	if err != nil {
		panic(err)
	}

	defer file.Close()

	hash := md5.New()
	_, err = io.Copy(hash, file)
	return hex.EncodeToString(hash.Sum(nil))
}

func PackageFiles(version_string string) error {
	data, err := os.ReadFile("./config.dat")
	HandleError(err, "读取配置文件")

	var str_config string = strings.TrimSpace(string(data))
	version_formatted_string := strings.ReplaceAll(version_string, ".", "_")
	cache_dictionary := filepath.Join(filepath.Dir(os.Args[0]), "gameasm/cache/"+version_formatted_string+"/")
	artifact_dictionary := filepath.Join(filepath.Dir(os.Args[0]), "gameasm/artifacts/")
	MkdirAllHandleError(artifact_dictionary)
	MkdirAllHandleError(cache_dictionary)

	fmt.Println("打包根目录:", str_config)
	fmt.Println("打包缓存文件目录:", cache_dictionary)
	fmt.Println()

	data, err = os.ReadFile(filepath.Join(artifact_dictionary, "_latest_index.json"))
	latest := map[string]string{}
	if err == nil {
		LoggerInfo("检测到有旧版本，将进行增量更新...")
		var index_config string = strings.TrimSpace(string(data))

		err = json.Unmarshal([]byte(index_config), &latest)
		HandleError(err, "Unmarshal _latest_index.json")
	} else {
		LoggerInfo("这是第一次发布版本，安装包可能较大...")
	}

	index := make(map[string]string)
	err = filepath.Walk(str_config, func(path string, info os.FileInfo, err error) error {

		if info.IsDir() {
			return nil
		}

		md5 := getmd5(path)
		clean_path := strings.ReplaceAll(path, str_config, "")
		clean_path = strings.ReplaceAll(clean_path, "\\", "/")
		copy_file := false

		val1, ok1 := latest["+"+md5]
		_, ok2 := latest["-"+md5]
		val3, ok3 := latest["*"+md5]

		if ok1 || ok2 || ok3 {
			if val1 == clean_path {
				// latest 中的文件有 hash 和文件路径都相同的文件
				LoggerInfo("发现未改变的文件: " + clean_path)
			} else if ok1 == true {
				// 上个版本中 hash 相同但文件路径不同
				index["*"+md5] = val1 + "=>" + clean_path
				LoggerInfo("发现内容相同但路径不同的文件: " + val1 + ", " + clean_path)
				copy_file = false
			} else if ok3 == true {
				to := strings.Split(val3, "=>")[1]
				if to != clean_path {
					// 上次移动的地方不是这次的地方，所以还得继续移动
					index["*"+md5] = to + "=>" + clean_path
					LoggerInfo("重新移动文件: " + to + ", " + clean_path)
				} else {
					LoggerInfo("发现未改变的文件: " + to + ", " + clean_path)
				}
			} else if ok2 == true {
				// 上次要求删除，这次又给加回来了
				LoggerInfo("新增文件: " + clean_path)
				index["+"+md5] = clean_path
				copy_file = true
			}
		} else {
			LoggerInfo("新增文件: " + clean_path)
			index["+"+md5] = clean_path
			copy_file = true
		}

		delete(latest, "+"+md5)
		delete(latest, "-"+md5)
		delete(latest, "*"+md5)

		if copy_file {
			dest_path := filepath.Join(cache_dictionary, md5)
			err = cp.Copy(path, dest_path)
			HandleError(err, "复制文件")
		}

		return nil
	})

	HandleError(err, "遍历文件")

	for key, value := range latest {

		if strings.HasPrefix(key, "+") {
			// 增量更新
			key = strings.TrimPrefix(key, "+")
			index["-"+key] = value
			LoggerInfo("发现已删除的文件: " + value)
		}

	}

	LoggerInfo("写入资源索引文件...")
	index_json, err := json.Marshal(index)
	err = os.WriteFile(filepath.Join(cache_dictionary, "index.json"), index_json, 0644)
	HandleError(err, "写入 index.json")
	err = os.WriteFile(filepath.Join(artifact_dictionary, "_latest_index.json"), index_json, 0644)
	HandleError(err, "写入 _latest_index.json")

	LoggerInfo("打包为 zip 文件...")
	zipDirectory(filepath.Join(artifact_dictionary, version_formatted_string+".zip"), cache_dictionary)
	err = os.RemoveAll(cache_dictionary)

	LoggerInfo("请不要删除 _latest_index.json 文件，这是为了下一次增量更新准备的。")

	exec.Command(`explorer.exe`, `/select,`, filepath.Join(artifact_dictionary, version_formatted_string+".zip")).Run()

	return nil
}
