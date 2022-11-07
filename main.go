package main

import (
	"bytes"
	"encoding/json"
	"fileserver/utils"
	"io"
	"net/http"
	"path"
	"time"

	"github.com/pterm/pterm"
)

var REMOTE_REPO = "https://storage.fastbuilder.pro/"
var LOCAL_REPO = "./files"
var PORT = ":12333"
var UPDATETIME = 1 * time.Hour

var isUpdating = false

func download(fileName string) ([]byte, error) {
	var compressedData []byte
	var execBytes []byte
	var err error
	path := path.Join(utils.GetCurrentDir(), LOCAL_REPO, fileName)
	url := REMOTE_REPO + fileName
	if compressedData, err = utils.DownloadSmallContent(url); err != nil {
		return nil, err
	}
	if execBytes, err = io.ReadAll(bytes.NewReader(compressedData)); err != nil {
		return nil, err
	}
	if err := utils.WriteFileData(path, execBytes); err != nil {
		return nil, err
	}
	return compressedData, nil
}

func updateRes() {
	utils.PrintfWithTime(pterm.Yellow("尝试进行资源更新.."))
	if isUpdating {
		utils.PrintfWithTime(pterm.Yellow("仍有资源处于下载状态，终止本次更新"))
		return
	}
	isUpdating = true
	jsonData, err := download("hashes.json")
	if err != nil {
		utils.PrintfWithTime(pterm.Red(pterm.Sprintf("无法从远程仓库获取hashes，将在 %s 后再次尝试更新", UPDATETIME)))
		isUpdating = false
		return
	}
	hashMap := make(map[string]string, 0)
	if err := json.Unmarshal([]byte(jsonData), &hashMap); err != nil {
		panic(err)
	}
	p, _ := pterm.DefaultProgressbar.WithTotal(len(hashMap)).WithTitle(pterm.Sprintf("%s %s %s", pterm.White(time.Now().Format("[15:04:05]")), pterm.Yellow("正在更新 ->"), pterm.White("FileName"))).Start()
	p.RemoveWhenDone = true
	for k, v := range hashMap {
		if v != utils.GetFileHash(path.Join(utils.GetCurrentDir(), LOCAL_REPO, k)) {
			p.UpdateTitle(pterm.Sprintf("%s %s %s", pterm.White(time.Now().Format("[15:04:05]")), pterm.Yellow("正在更新 ->"), k))
			if _, err := download(k); err != nil {
				utils.PrintfWithTime(pterm.Sprintf("%s %s", pterm.Red("更新失败 ->"), k))
				p.Increment()
				continue
			}
			utils.PrintfWithTime(pterm.Sprintf("%s %s", pterm.Green("更新完成 ->"), k))
		} else {
			utils.PrintfWithTime(pterm.Sprintf("%s %s", pterm.LightCyan("无需更新 ->"), k))
		}
		p.Increment()
	}
	utils.PrintfWithTime(pterm.Yellow(pterm.Sprintf("资源更新执行完毕，将在 %s 后再次检查更新", UPDATETIME)))
	isUpdating = false
}

func main() {
	pterm.DefaultBox.Println("https://github.com/Liliya233/simple_mirror_file_site")
	filePath := path.Join(utils.GetCurrentDir(), LOCAL_REPO)
	utils.PrintfWithTime(pterm.Sprintf("%s %s", pterm.LightCyan("将使用此地址作为远程仓库:"), REMOTE_REPO))
	utils.PrintfWithTime(pterm.Sprintf("%s %s", pterm.LightCyan("将使用此目录搭建文件服务器:"), filePath))
	utils.PrintfWithTime(pterm.Sprintf("%s %s", pterm.LightCyan("将使用此IP搭建文件服务器:"), PORT))
	if !utils.IsDir(filePath) {
		utils.MkDir(filePath)
	}
	utils.PrintfWithTime(pterm.Yellow("文件服务器将在首次资源更新完成后启动"))
	updateRes()
	ticker := time.NewTicker(UPDATETIME)
	go func() {
		for {
			<-ticker.C
			updateRes()
		}
	}()
	http.HandleFunc("/res/", func(w http.ResponseWriter, r *http.Request) {
		ip, _ := utils.GetIP(r)
		utils.PrintfWithTime(pterm.Sprintf("%s %s %s %s", pterm.Green("接受访问:"), pterm.Yellow(ip), pterm.Cyan("->"), r.URL.Path))
		http.StripPrefix("/res/", http.FileServer(http.Dir(filePath))).ServeHTTP(w, r)
	})
	utils.PrintfWithTime(pterm.Green("更新完毕，将启动文件服务器"))
	http.ListenAndServe(PORT, nil)
}
