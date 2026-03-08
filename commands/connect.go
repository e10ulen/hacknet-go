// commands/connect.go
package commands

import (
	"fmt"

	"github.com/e10ulen/hacknet-go/vfs"
)

// HandleConnect は connect コマンド：サーバー接続
func HandleConnect(args []string, vfs *vfs.VFS, log *[]string) (bool, string, error) {
	if len(args) < 1 {
		*log = append(*log, "使い方: connect <IP>")
		return false, "", nil
	}

	ip := args[0]
	*log = append(*log, fmt.Sprintf("接続試行中: %s ...", ip))

	// ここでは常に成功とする（実際はランダム失敗や認証を追加可能）
	*log = append(*log, fmt.Sprintf("接続成功: %s", ip))

	// model に反映させるために true と新しいIPを返す
	return true, ip, nil
}
