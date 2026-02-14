package collector

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

type Options struct {
	Sources []string
	Dest    string
	DryRun  bool
}

type ImageFile struct {
	Path    string
	ModTime time.Time
}

type Result struct {
	Copied  int
	Skipped int
	Failed  int
	Errors  []error
}

var imageExts = map[string]bool{
	".jpg": true, ".jpeg": true, ".png": true, ".gif": true, ".bmp": true, ".webp": true,
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
	defer func() {
		_ = out.Close()
	}()

	if _, err := io.Copy(out, in); err != nil {
		return err
	}
	return nil
}

func Run(opt Options) (Result, error) {
	var res Result

	if err := os.MkdirAll(opt.Dest, 0755); err != nil {
		return res, err
	}

	for _, srcDir := range opt.Sources {
		// folderName := filepath.Base(srcDir)

		// 1) 対象ファイル列挙（順序固定のためスライスに集めてsort）
		var files []ImageFile
		err := filepath.WalkDir(srcDir, func(path string, d os.DirEntry, err error) error {
			if err != nil {
				return err
			}
			if d.IsDir() {
				return nil
			}
			if !isImage(path) {
				return nil
			}

			info, err := d.Info()
			if err != nil {
				return err
			}

			files = append(files, ImageFile{
				Path:    path,
				ModTime: info.ModTime(),
			})
			return nil
		})
		if err != nil {
			return res, err
		}

		// 更新日時でソート
		sort.Slice(files, func(i, j int) bool {
			return files[i].ModTime.Before(files[j].ModTime)
		})

		// 2) フォルダごとに連番
		seqByFolder := map[string]int{}
		for _, file := range files {

			folderName := filepath.Base(filepath.Dir(file.Path))

			seqByFolder[folderName]++
			seq := seqByFolder[folderName]

			ext := strings.ToLower(filepath.Ext(file.Path))
			newName := fmt.Sprintf("%s_%04d%s", folderName, seq, ext)
			dstPath := filepath.Join(opt.Dest, newName)

			// ※MVP：同名があったらスキップ（後で重複回避に進化させる）
			if _, err := os.Stat(dstPath); err == nil {
				res.Skipped++
				seq++
				continue
			}

			if opt.DryRun {
				fmt.Printf("[dry-run] %s -> %s\n", file.Path, dstPath)
				res.Skipped++ // dry-runは実コピーしないのでskipped扱いでもOK（好みで別カウントにしても）
				seq++
				continue
			}

			if err := copyFile(file.Path, dstPath); err != nil {
				res.Failed++
				res.Errors = append(res.Errors, fmt.Errorf("copy failed: %s: %w", file.Path, err))
			} else {
				res.Copied++
			}
			seq++
		}
	}

	return res, nil
}
