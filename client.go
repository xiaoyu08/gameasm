package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	cp "github.com/otiai10/copy"
)

func Client(args []string) {
	if len(args) > 2 {
		switch args[2] {
		case "install":
			Recover(args)
			return
		}
	}
	Helper(2)
}

func Recover(args []string) {
	fmt.Println()
	var filename string = ""
	if len(args) != 4 {
		fmt.Println("Need parameters")
		os.Exit(1)
	} else {
		filename = args[3]
	}
	UnpackageFiles(filename)
}

func UnpackageFiles(filename string) error {

	recover_dictionary := filepath.Join(filepath.Dir(os.Args[0]), "gameasm/game/")
	unzipped_dictionary := filepath.Join(filepath.Dir(os.Args[0]), "gameasm/game/unzipped/")
	MkdirAllHandleError(unzipped_dictionary)
	err := Unzip(filepath.Join(filepath.Dir(os.Args[0]), filename), unzipped_dictionary)
	if err != nil {
		HandleError(err, "Unzip file")
	}

	data, err := os.ReadFile(filepath.Join(unzipped_dictionary, "index.json"))
	HandleError(err, "Reading index.json")

	var index_config string = strings.TrimSpace(string(data))

	x := map[string]string{}
	err = json.Unmarshal([]byte(index_config), &x)
	HandleError(err, "Unmarshal index.json")
	count := 0

	for key, value := range x {

		if strings.HasPrefix(key, "+") {
			// 增量更新
			key = strings.TrimPrefix(key, "+")
			original_file := filepath.Join(recover_dictionary, value)
			hash_file := filepath.Join(unzipped_dictionary, key)
			MkdirAllHandleError(filepath.Dir(original_file))
			LoggerInfo("拷贝新文件: " + key + " => " + original_file)
			cp.Copy(filepath.Join(unzipped_dictionary, key), original_file)
			os.Remove(hash_file)
			count += 1
		} else if strings.HasPrefix(key, "-") {
			// 要求删除文件
			original_file := filepath.Join(recover_dictionary, value)
			os.Remove(original_file)
			LoggerInfo("删除文件: " + value)
			count += 1
		} else if strings.HasPrefix(key, "*") {
			// 要求移动文件
			from := strings.Split(value, "=>")[0]
			to := strings.Split(value, "=>")[1]
			original_file := filepath.Join(recover_dictionary, from)
			dest_file := filepath.Join(recover_dictionary, to)
			cp.Copy(original_file, dest_file)
			os.Remove(original_file)
			LoggerInfo("移动文件: " + from + " --> " + to)
			count += 1
		}

	}

	os.RemoveAll(unzipped_dictionary)

	LoggerInfo("Done. " + fmt.Sprint(count) + " files recovered.")

	return nil
}
