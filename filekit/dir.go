package filekit

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// CopyDir 函数接受两个字符串参数，分别表示源文件夹和目标文件夹的路径
// 返回值是拷贝的文件数量和可能发生的错误
func Copy(src, dst string) (count int, err error) {
	// 使用正则表达式将路径按照 / 或 \ 分割成切片
	regexp1, err := regexp.Compile(`(/|\\)`)
	if err != nil {
		return 0, err
	}
	srcSplits := regexp1.Split(src, 10000)
	dstSplits := regexp1.Split(dst, 10000)

	// 调用 CopyDirInner 函数，传入源文件夹和目标文件夹的前缀和最后一级名称
	return CopyInner(strings.Join(srcSplits[:len(srcSplits)-1], "/"), srcSplits[len(srcSplits)-1], strings.Join(dstSplits[:len(dstSplits)-1], "/"), dstSplits[len(dstSplits)-1])
}

// CopyDirInner 函数接受四个字符串参数，分别表示源文件夹和目标文件夹的前缀和最后一级名称
// 返回值是拷贝的文件数量和可能发生的错误
func CopyInner(srcPrefix, src string, dstPrefix, dst string) (count int, err error) {
	// 如果前缀为空，则设置为当前目录
	if srcPrefix == "" {
		srcPrefix = "."
	}
	if dstPrefix == "" {
		dstPrefix = "."
	}
	// 读取源文件夹下的所有文件和子文件夹
	dirs, err := os.ReadDir(srcPrefix + "/" + src)
	if err != nil {
		return 0, err
	}

	// 在目标文件夹下创建同名的子文件夹
	pathCursor := dstPrefix + "/" + dst + "/" + src
	err = os.MkdirAll(pathCursor, 0600)
	if err != nil {
		return 0, err
	}
	for _, dir := range dirs {
		if dir.IsDir() {
			// 如果是子文件夹，则递归调用 CopyDirInner 函数，传入相应的参数
			countSub, err := CopyInner(srcPrefix+"/"+src, dir.Name(), dstPrefix+"/"+dst, src)
			if err != nil {
				return 0, err
			}
			count += countSub
		} else {
			// 如果是文件，则读取其内容，并写入到目标文件夹下同名的文件中
			bytesFile, err := os.ReadFile(srcPrefix + "/" + src + "/" + dir.Name())
			if err != nil {
				return 0, err
			}

			err = os.WriteFile(pathCursor+"/"+dir.Name(), bytesFile, 0600)
			if err != nil {
				return 0, err
			}
			count++
		}
	}
	return count, nil
}

// CopyDir 递归复制目录及其内容
func CopyDir(src, dst string) error {
	srcInfo, err := os.Stat(src)
	if err != nil {
		return err
	}

	if !srcInfo.IsDir() {
		return fmt.Errorf("%s 不是目录", src)
	}

	// 创建目标目录
	err = os.MkdirAll(dst, srcInfo.Mode())
	if err != nil {
		return err
	}

	entries, err := os.ReadDir(src)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		srcPath := filepath.Join(src, entry.Name())
		dstPath := filepath.Join(dst, entry.Name())

		if entry.IsDir() {
			if err := CopyDir(srcPath, dstPath); err != nil {
				return err
			}
		} else {
			// 处理文件和符号链接
			if entry.Type()&os.ModeSymlink != 0 {
				if err := copySymlink(srcPath, dstPath); err != nil {
					return err
				}
			} else {
				if err := CopyFile(srcPath, dstPath); err != nil {
					return err
				}
			}
		}
	}

	return nil
}

// CopyFile 复制单个文件
func CopyFile(src, dst string) error {
	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	srcInfo, err := srcFile.Stat()
	if err != nil {
		return err
	}

	// 如果目标文件已存在，检查是否需要覆盖
	if _, err := os.Stat(dst); err == nil {
		dstInfo, err := os.Stat(dst)
		if err != nil {
			return err
		}

		// 如果源文件和目标文件相同，跳过复制
		if os.SameFile(srcInfo, dstInfo) {
			return nil
		}
	}

	dstFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer dstFile.Close()

	// 复制文件内容
	_, err = io.Copy(dstFile, srcFile)
	if err != nil {
		return err
	}

	// 复制文件权限
	if err = os.Chmod(dst, srcInfo.Mode()); err != nil {
		return err
	}

	// 尝试复制文件时间

	return nil
}

// copySymlink 复制符号链接
func copySymlink(src, dst string) error {
	target, err := os.Readlink(src)
	if err != nil {
		return err
	}

	return os.Symlink(target, dst)
}
