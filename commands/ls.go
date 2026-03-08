// commands/ls.go
package commands

import (
	"fmt"
	"strings"

	"github.com/e10ulen/hacknet-go/vfs"
)

// HandleLs は ls コマンドの処理
func HandleLs(args []string, vfs *vfs.VFS, log *[]string) (bool, string, error) {
	files := vfs.ListFiles()
	if len(files) == 0 {
		*log = append(*log, "ディレクトリは空です")
		return false, "", nil
	}

	*log = append(*log, fmt.Sprintf("現在のディレクトリ: %s", vfs.GetPath()))
	*log = append(*log, "  "+strings.Join(files, "  "))
	// 状態変更なし → false, "", nil を返す
	return false, "", nil
}
