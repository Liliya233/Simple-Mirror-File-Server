package main

import (
	"encoding/json"
	"fileserver/utils"
	"fmt"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	_ "embed"

	"github.com/pterm/pterm"
)

var (
	//go:embed assets/default_config.json
	defaultConfig []byte
	CurrentConfig *Config = &Config{}
)

// 配置文件结构
type Config struct {
	LocalDir   string    `json:"本地文件夹"`
	Address    string    `json:"IP地址"`
	Port       uint32    `json:"开放端口"`
	MirrorUrl  string    `json:"Github镜像站"`
	CheckCycle int64     `json:"检查更新频率(秒)"`
	SourceList []*Source `json:"资源列表"`
}

// 资源项
type Source struct {
	SubDirName   string `json:"子文件夹名称"`
	Url          string `json:"资源地址"`
	HashFileName string `json:"Hash文件名称"`
}

func download(DownloadUrl, dstFile string) (err error) {
	var downloadBytes []byte
	if downloadBytes, err = utils.DownloadSmallContent(DownloadUrl); err != nil {
		return err
	}
	if err := utils.WriteFileData(dstFile, downloadBytes); err != nil {
		return err
	}
	return nil
}

func updateRes(source *Source) {
	utils.PrintWarn("检查更新: " + source.SubDirName)
	// 如有设置镜像站, 检查是否为Github链接
	if CurrentConfig.MirrorUrl != "" && !strings.HasPrefix(source.Url, "https://github.com/") {
		utils.PrintError("非 Github 链接无法使用镜像站: " + source.SubDirName)
		return
	}
	// 设置下载url
	downloadUrl := CurrentConfig.MirrorUrl + source.Url
	// 获取 hash 及文件列表
	hashesBytes, err := utils.DownloadSmallContent(downloadUrl + source.HashFileName)
	if err != nil {
		utils.PrintError("下载hash文件时出现错误: " + err.Error())
		return
	}
	// 解析 hash 及文件列表
	hashMap := make(map[string]string, 0)
	if err := json.Unmarshal([]byte(hashesBytes), &hashMap); err != nil {
		utils.PrintError("解析hash文件时出现错误, 可能指定的文件不正确: " + err.Error())
		return
	}
	// 设置文件夹
	downloadDir := filepath.Join(utils.GetCurrentDir(), CurrentConfig.LocalDir, source.SubDirName)
	// 初始化文件夹
	utils.MkDir(downloadDir)
	// 将hash文件保存到文件夹
	utils.WriteFileData(filepath.Join(downloadDir, source.HashFileName), hashesBytes)
	// 开始更新
	p, _ := pterm.DefaultProgressbar.WithTotal(len(hashMap)).WithTitle(pterm.Sprintf("%s %s %s", pterm.White(time.Now().Format("[15:04:05]")), pterm.Yellow("正在更新 ->"), pterm.White("FileName"))).Start()
	p.RemoveWhenDone = true
	success := true
	// 遍历文件列表
	for fileName, fileHash := range hashMap {
		// hash不一致时进行更新
		if fileHash != utils.GetFileHash(filepath.Join(downloadDir, fileName)) {
			// 打印开始更新信息
			p.UpdateTitle(pterm.Sprintf("%s %s %s", pterm.White(time.Now().Format("[15:04:05]")), pterm.Yellow("正在更新 ->"), fileName))
			// 尝试更新
			if err := download(downloadUrl+fileName, filepath.Join(downloadDir, fileName)); err != nil {
				utils.PrintError("更新失败 -> " + fileName)
				success = false
			} else {
				utils.PrintSuccess("更新完成 -> " + fileName)
			}
		} else {
			utils.PrintInfo("无需更新 -> " + fileName)
		}
		p.Increment()
	}
	if success {
		utils.PrintSuccess("全部资源更新完成: " + source.SubDirName)
	} else {
		// 失败时继续更新
		utils.PrintWarn("存在未成功更新的资源, 将再次尝试更新: " + source.SubDirName)
		updateRes(source)
	}
}

func main() {
	// 打印项目地址
	pterm.DefaultBox.Println("https://github.com/Liliya233/simple_mirror_file_site")
	// 读取配置
	configPath := filepath.Join(utils.GetCurrentDir(), "config.json")
	if utils.GetJsonData(configPath, CurrentConfig) != nil {
		// 读取失败, 使用默认配置
		utils.PrintWarn("未能读取到配置文件, 将生成并使用默认配置")
		utils.WriteFileData(configPath, defaultConfig)
		utils.GetJsonData(configPath, CurrentConfig)
	}
	// 打印运行信息
	utils.PrintInfo("将基于此目录搭建文件服务器: " + CurrentConfig.LocalDir)
	utils.PrintInfo("将使用此IP搭建文件服务器: " + fmt.Sprintf("%s:%d", CurrentConfig.Address, CurrentConfig.Port))
	// 初始化文件夹
	utils.MkDir(CurrentConfig.LocalDir)
	// 启动文件更新协程
	ticker := time.NewTicker(time.Duration(CurrentConfig.CheckCycle) * time.Second)
	go func() {
		for {
			for _, source := range CurrentConfig.SourceList {
				updateRes(source)
			}
			<-ticker.C
		}
	}()
	http.HandleFunc(fmt.Sprintf("/%s/", CurrentConfig.LocalDir), func(w http.ResponseWriter, r *http.Request) {
		ip, _ := utils.GetIP(r)
		utils.PrintInfo("接受访问: " + ip + "->" + r.URL.Path)
		http.StripPrefix(fmt.Sprintf("/%s/", CurrentConfig.LocalDir), http.FileServer(http.Dir(filepath.Join(utils.GetCurrentDir(), CurrentConfig.LocalDir)))).ServeHTTP(w, r)
	})
	utils.PrintSuccess("文件服务器已启动")
	http.ListenAndServe(fmt.Sprintf("%s:%d", CurrentConfig.Address, CurrentConfig.Port), nil)
}
