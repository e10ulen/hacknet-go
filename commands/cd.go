// commands/cd.go
package commands

import (
	"fmt"

	"github.com/e10ulen/hacknet-go/vfs"
)

// HandleCd は cd コマンドの処理
func HandleCd(args []string, vfs *vfs.VFS, log *[]string) (bool, string, error) {
	if len(args) < 1 {
		*log = append(*log, "使い方: cd <ディレクトリ>")
		return false, "", nil
	}

	target := args[0]
	err := vfs.ChangeDir(target)
	if err != nil {
		*log = append(*log, fmt.Sprintf("エラー: %v", err))
		return false, "", nil
	}

	*log = append(*log, fmt.Sprintf("ディレクトリを変更しました: %s", vfs.GetPath()))
	// 状態変更なし → false, "", nil を返す
	return false, "", nil
}
