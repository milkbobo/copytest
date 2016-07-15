package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

var Config struct {
	originDir string
	targetDir string
}

var err error

func main() {

	//获取参数
	ReadConfig()

	//读取需要的文件
	data, err := ReadDir(Config.originDir)
	if err != nil {
		fmt.Println("read dir error " + err.Error())
		os.Exit(1)
		return
	}

	//便利复制文件
	for _, singleFile := range data {
		_, err := CopyFile(singleFile, Config.targetDir)
		if err != nil {
			panic(err.Error())
		}
	}
}

func ReadConfig() {

	argv := os.Args
	argv = argv[1:]
	if len(argv) < 1 {
		panic("请填写第一参数，源文件夹，第二参数，目标文件夹")
	}
	//获取绝对地址
	Config.originDir, err = filepath.Abs(argv[0])
	if err != nil {
		panic(err)
	}
	Config.targetDir, err = filepath.Abs(argv[1])
	if err != nil {
		panic(err)
	}

}

func combineDirInfo(a1 []string, a2 []string) []string {
	for _, value := range a2 {
		if len(value) != 0 {
			a1 = append(a1, value)
		}
	}
	return a1
}

func ReadDir(path string) ([]string, error) {
	result := []string{}
	//获取文件夹所有的文件列表
	tempResult, err := ioutil.ReadDir(path)
	if err != nil {
		return nil, err
	}
	for _, singleFileInfo := range tempResult {
		name := path + "/" + singleFileInfo.Name()

		if singleFileInfo.IsDir() {
			//发现dir
			result2, err := ReadDir(name)
			if err != nil {
				return nil, err
			}
			result = combineDirInfo(result, result2)
		} else {
			//发现源代码文件
			if strings.HasSuffix(name, "_testing.go") == false && strings.HasSuffix(name, "inittestdatabase.go") == false {
				continue
			}
			result = append(result, name)
		}
	}
	return result, nil

}

func CopyFile(src, dst string) (w int64, err error) {

	srcFile, err := os.Open(src)
	if err != nil {
		panic(err.Error())
	}
	defer srcFile.Close()

	//获取输出的文件
	dst = Config.targetDir + strings.Replace(src, Config.originDir, "", -1)

	//新建整条目录文件夹
	err = os.MkdirAll((filepath.Dir(dst)), os.ModePerm)
	if err != nil {
		panic(err)
	}

	dstFile, err := os.Create(dst)

	if err != nil {
		panic(err.Error())
	}

	defer dstFile.Close()
	// log(src)
	// log(dst)
	return io.Copy(dstFile, srcFile)
}

//方便调试输出使用
func log(data interface{}) {
	fmt.Printf("%#v \n", data)
}
