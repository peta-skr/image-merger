package collector

import (
	"io"
	"os"
	"path/filepath"
	"strings"
)

var imageExts = map[string]bool{
	".jpg":  true,
	".jpeg": true,
	".png":  true,
	".gif":  true,
	".bmp":  true,
	".webp": true,
}

func isImage(path string) bool {
	ext := strings.ToLower(filepath.Ext(path))
	return imageExts[ext]
}

func copyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}

	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, in)
	return err
}

func CollectImages(srcDir, dstDir string) error {
	// 出力フォルダ作成
	err := os.MkdirAll(dstDir, 0755)

	if err != nil {
		return err
	}

	return filepath.Walk(srcDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		if !isImage(path) {
			return nil
		}

		dstPath := filepath.Join(dstDir, info.Name())

		return copyFile(path, dstPath)
	})
}
