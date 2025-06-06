package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/88250/lute"
	"github.com/Wsine/feishu2md/core"
	"github.com/Wsine/feishu2md/utils"
	"github.com/chyroc/lark"
	"github.com/pkg/errors"
)

type DownloadOpts struct {
	outputDir string
	dump      bool
	batch     bool
	wiki      bool
	force     bool
	appId     string
	appSecret string
}

var dlOpts = DownloadOpts{}
var dlConfig core.Config

func downloadDocument(ctx context.Context, client *core.Client, url string, opts *DownloadOpts) error {
	// Validate the url to download
	docType, docToken, err := utils.ValidateDocumentURL(url)
	if err != nil {
		return err
	}
	fmt.Println("Captured document token:", docToken)

	// for a wiki page, we need to renew docType and docToken first
	if docType == "wiki" {
		node, err := client.GetWikiNodeInfo(ctx, docToken)
		if err != nil {
			err = fmt.Errorf("GetWikiNodeInfo err: %v for %v", err, url)
		}
		utils.CheckErr(err)
		docType = node.ObjType
		docToken = node.ObjToken
	}
	if docType == "docs" {
		return errors.Errorf(
			`Feishu Docs is no longer supported. ` +
				`Please refer to the Readme/Release for v1_support.`)
	}

	// Process the download
	docx, blocks, err := client.GetDocxContent(ctx, docToken)
	utils.CheckErr(err)

	parser := core.NewParser(dlConfig.Output)

	title := docx.Title
	markdown := parser.ParseDocxContent(docx, blocks)

	if !dlConfig.Output.SkipImgDownload {
		for _, imgToken := range parser.ImgTokens {
			localLink, err := client.DownloadImage(
				ctx, imgToken, filepath.Join(opts.outputDir, dlConfig.Output.ImageDir),
			)
			if err != nil {
				return err
			}
			markdown = strings.Replace(markdown, imgToken, localLink, 1)
		}
	}

	// Format the markdown document
	engine := lute.New(func(l *lute.Lute) {
		l.RenderOptions.AutoSpace = true
	})
	result := engine.FormatStr("md", markdown)

	// Handle the output directory and name
	if _, err := os.Stat(opts.outputDir); os.IsNotExist(err) {
		if err := os.MkdirAll(opts.outputDir, 0o755); err != nil {
			return err
		}
	}

	if dlOpts.dump {
		jsonName := fmt.Sprintf("%s.json", docToken)
		outputPath := filepath.Join(opts.outputDir, jsonName)
		data := struct {
			Document *lark.DocxDocument `json:"document"`
			Blocks   []*lark.DocxBlock  `json:"blocks"`
		}{
			Document: docx,
			Blocks:   blocks,
		}
		pdata := utils.PrettyPrint(data)

		if err = os.WriteFile(outputPath, []byte(pdata), 0o644); err != nil {
			return err
		}
		fmt.Printf("Dumped json response to %s\n", outputPath)
	}

	// Write to markdown file
	mdName := fmt.Sprintf("%s.md", docToken)
	if dlConfig.Output.TitleAsFilename {
		mdName = fmt.Sprintf("%s.md", utils.SanitizeFileName(title))
	}
	outputPath := filepath.Join(opts.outputDir, mdName)

	if dlConfig.Output.Delta && !opts.force {
		if _, err := os.Stat(outputPath); !os.IsNotExist(err) {
			fmt.Printf("File %s already exists, skipping download.\n", outputPath)
			return nil
		}
	}

	if err = os.WriteFile(outputPath, []byte(result), 0o644); err != nil {
		return err
	}
	fmt.Printf("Downloaded markdown file to %s\n", outputPath)

	return nil
}

func downloadDocuments(ctx context.Context, client *core.Client, url string) error {
	// Validate the url to download
	folderToken, err := utils.ValidateFolderURL(url)
	if err != nil {
		return err
	}
	fmt.Println("Captured folder token:", folderToken)

	// Error channel and wait group
	errChan := make(chan error)
	wg := sync.WaitGroup{}

	// Recursively go through the folder and download the documents
	var processFolder func(ctx context.Context, folderPath, folderToken string) error
	processFolder = func(ctx context.Context, folderPath, folderToken string) error {
		files, err := client.GetDriveFolderFileList(ctx, nil, &folderToken)
		if err != nil {
			return err
		}
		opts := DownloadOpts{outputDir: folderPath, dump: dlOpts.dump, batch: false, force: dlOpts.force}
		for _, file := range files {
			if file.Type == "folder" {
				_folderPath := filepath.Join(folderPath, file.Name)
				if err := processFolder(ctx, _folderPath, file.Token); err != nil {
					return err
				}
			} else if file.Type == "docx" {
				// concurrently download the document
				wg.Add(1)
				go func(_url string) {
					if err := downloadDocument(ctx, client, _url, &opts); err != nil {
						errChan <- err
					}
					wg.Done()
				}(file.URL)
			}
		}
		return nil
	}
	if err := processFolder(ctx, dlOpts.outputDir, folderToken); err != nil {
		return err
	}

	// Wait for all the downloads to finish
	go func() {
		wg.Wait()
		close(errChan)
	}()
	for err := range errChan {
		return err
	}
	return nil
}

func downloadWiki(ctx context.Context, client *core.Client, url string) error {
	prefixURL, spaceID, err := utils.ValidateWikiURL(url)
	if err != nil {
		return err
	}

	// 获取知识库名称，但不用作输出目录
	wikiName, err := client.GetWikiName(ctx, spaceID)
	if err != nil {
		return err
	}
	if wikiName == "" {
		return fmt.Errorf("failed to GetWikiName")
	}

	// 使用用户指定的输出目录
	folderPath := dlOpts.outputDir
	// 打印知识库名称供用户参考
	fmt.Printf("Downloading wiki: %s to directory: %s\n", wikiName, folderPath)

	errChan := make(chan error)

	var maxConcurrency = 10 // Set the maximum concurrency level
	wg := sync.WaitGroup{}
	semaphore := make(chan struct{}, maxConcurrency) // Create a semaphore with the maximum concurrency level

	var downloadWikiNode func(ctx context.Context,
		client *core.Client,
		spaceID string,
		parentPath string,
		parentNodeToken *string) error

	downloadWikiNode = func(ctx context.Context,
		client *core.Client,
		spaceID string,
		folderPath string,
		parentNodeToken *string) error {
		nodes, err := client.GetWikiNodeList(ctx, spaceID, parentNodeToken)
		if err != nil {
			return err
		}
		for _, n := range nodes {
			if n.HasChild {
				_folderPath := filepath.Join(folderPath, n.Title)
				if err := downloadWikiNode(ctx, client,
					spaceID, _folderPath, &n.NodeToken); err != nil {
					return err
				}
			}
			if n.ObjType == "docx" {
				opts := DownloadOpts{outputDir: folderPath, dump: dlOpts.dump, batch: false, force: dlOpts.force}
				wg.Add(1)
				semaphore <- struct{}{}
				go func(_url string) {
					if err := downloadDocument(ctx, client, _url, &opts); err != nil {
						errChan <- err
					}
					wg.Done()
					<-semaphore
				}(prefixURL + "/wiki/" + n.NodeToken)
				// downloadDocument(ctx, client, prefixURL+"/wiki/"+n.NodeToken, &opts)
			}
		}
		return nil
	}

	if err = downloadWikiNode(ctx, client, spaceID, folderPath, nil); err != nil {
		return err
	}

	// Wait for all the downloads to finish
	go func() {
		wg.Wait()
		close(errChan)
	}()
	for err := range errChan {
		return err
	}
	return nil
}

func handleDownloadCommand(url string) error {
	// Load config
	configPath, err := core.GetConfigFilePath()
	if err != nil {
		return err
	}

	// 检查配置文件是否存在，如果不存在则创建默认配置
	var config *core.Config
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		// 创建默认配置，使用命令行传入的appId和appSecret
		config = core.NewConfig(dlOpts.appId, dlOpts.appSecret)
		if err = config.WriteConfig2File(configPath); err != nil {
			return err
		}
		fmt.Println("Created default configuration file at: " + configPath)
		if dlOpts.appId == "" || dlOpts.appSecret == "" {
			fmt.Println("Please set your AppID and AppSecret using 'feishu2md config' command")
		}
	} else {
		// 读取现有配置
		config, err = core.ReadConfigFromFile(configPath)
		if err != nil {
			return err
		}
		// 如果命令行传入了appId或appSecret，则更新配置
		if dlOpts.appId != "" || dlOpts.appSecret != "" {
			if dlOpts.appId != "" {
				config.Feishu.AppId = dlOpts.appId
			}
			if dlOpts.appSecret != "" {
				config.Feishu.AppSecret = dlOpts.appSecret
			}
			// 将更新后的配置写回文件
			if err = config.WriteConfig2File(configPath); err != nil {
				return err
			}
			fmt.Println("Updated configuration with provided AppID and/or AppSecret")
		}
	}
	dlConfig = *config

	// Instantiate the client
	client := core.NewClient(
		dlConfig.Feishu.AppId, dlConfig.Feishu.AppSecret,
	)
	ctx := context.Background()

	if dlOpts.batch {
		return downloadDocuments(ctx, client, url)
	}

	if dlOpts.wiki {
		return downloadWiki(ctx, client, url)
	}

	return downloadDocument(ctx, client, url, &dlOpts)
}
