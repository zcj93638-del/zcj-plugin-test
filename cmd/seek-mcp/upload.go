package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
	"unicode"
)

type uploadArgs struct {
	FilePath    string `json:"filePath"`
	ProjectGUID string `json:"projectGuid"`
	FolderGUID  string `json:"parentFolderGuid"`
	Conflict    string `json:"onNameConflict"`
}

func uploadTool() obj {
	return obj{"name": "seek_upload", "description": "把用户指定的一个本地普通文件上传到指定 Seek 项目或文件夹。使用 MCP 运行机器可读的绝对路径。默认同名保留两份，仅在用户明确要求时使用 overwrite。返回素材 GUID。", "inputSchema": schema(obj{
		"filePath":         obj{"type": "string", "description": "所选文件在插件所在电脑上的绝对路径"},
		"projectGuid":      obj{"type": "string", "description": "目标项目 GUID"},
		"parentFolderGuid": obj{"type": "string", "description": "目标项目内的文件夹 GUID。省略或空串表示项目根目录"},
		"onNameConflict":   obj{"type": "string", "enum": []string{"keepBoth", "skip", "overwrite"}, "default": "keepBoth"},
	}, "filePath", "projectGuid")}
}

func (c connection) upload(a uploadArgs) (obj, error) {
	if !c.ready() {
		return nil, errors.New(incomplete)
	}
	if !filepath.IsAbs(a.FilePath) || strings.TrimSpace(a.ProjectGUID) == "" {
		return nil, errors.New("请提供本地文件绝对路径和目标项目 GUID。")
	}
	if a.Conflict == "" {
		a.Conflict = "keepBoth"
	}
	switch a.Conflict {
	case "keepBoth", "skip", "overwrite":
	default:
		return nil, errors.New("无效的同名处理策略。")
	}
	name := filepath.Base(a.FilePath)
	if strings.IndexFunc(name, unicode.IsControl) >= 0 || strings.Contains(name, `\`) {
		return nil, errors.New("文件名含有不支持的字符。")
	}
	info, err := os.Stat(a.FilePath)
	if err != nil {
		return nil, err
	}
	if !info.Mode().IsRegular() {
		return nil, errors.New("请选择一个普通文件。")
	}
	f, err := os.Open(a.FilePath)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	info, err = f.Stat()
	if err != nil {
		return nil, err
	}
	if !info.Mode().IsRegular() {
		return nil, errors.New("请选择一个普通文件。")
	}
	// The NAS stops parsing fields at the file part. Buffer just the framing so
	// the HTTP client can stream the file with an exact Content-Length.
	var framing bytes.Buffer
	writer := multipart.NewWriter(&framing)
	for _, field := range [][2]string{{"projectGuid", a.ProjectGUID}, {"parentFolderGuid", a.FolderGUID}, {"onNameConflict", a.Conflict}, {"size", strconv.FormatInt(info.Size(), 10)}} {
		if err := writer.WriteField(field[0], field[1]); err != nil {
			return nil, err
		}
	}
	if _, err = writer.CreateFormFile("file", name); err != nil {
		return nil, err
	}
	split := framing.Len()
	if err = writer.Close(); err != nil {
		return nil, err
	}
	parts := framing.Bytes()
	body := io.MultiReader(bytes.NewReader(parts[:split]), io.NewSectionReader(f, 0, info.Size()), bytes.NewReader(parts[split:]))
	req, err := http.NewRequest(http.MethodPost, c.Address+"/upload/api/v1/asset", body)
	if err != nil {
		return nil, err
	}
	req.ContentLength = int64(len(parts)) + info.Size()
	req.Header.Set("Content-Type", writer.FormDataContentType())
	uploadClient := *client
	uploadClient.Timeout = 30 * time.Minute
	data, _, err := c.send(&uploadClient, req)
	if err != nil {
		return nil, err
	}
	var asset struct {
		GUID string `json:"guid"`
	}
	if json.Unmarshal(data, &asset) != nil || asset.GUID == "" {
		return nil, errors.New("Seek 返回了成功状态但缺少素材 GUID，请检查目标文件夹后再重试。")
	}
	return obj{"guid": asset.GUID, "fileName": name, "projectGuid": a.ProjectGUID, "parentFolderGuid": a.FolderGUID, "onNameConflict": a.Conflict}, nil
}
