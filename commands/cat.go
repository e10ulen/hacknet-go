// commands/cat.go
package commands

import (
	"fmt"

	"github.com/e10ulen/hacknet-go/vfs"
)

// HandleCat は cat コマンドの処理
func HandleCat(args []string, vfs *vfs.VFS, log *[]string) (bool, string, error) {
	if len(args) < 1 {
		*log = append(*log, "使い方: cat <ファイル>")
		return false, "", nil
	}

	content, err := vfs.ReadFile(args[0])
	if err != nil {
		*log = append(*log, fmt.Sprintf("エラー: %v", err))
		return false, "", nil
	}

	*log = append(*log, fmt.Sprintf("--- %s ---", args[0]))
	*log = append(*log, content)
	*log = append(*log, "--- EOF ---")

	// 状態変更なし → false, "", nil を返す
	return false, "", nil
}
