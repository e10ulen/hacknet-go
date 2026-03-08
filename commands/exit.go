// commands/exit.go
package commands

import "github.com/e10ulen/hacknet-go/vfs"

// HandleExit は exit コマンド：接続切断
func HandleExit(args []string, vfs *vfs.VFS, log *[]string) (bool, string, error) {
	*log = append(*log, "切断しました。")
	// connectとは逆に false を返す（切断なので）
	return false, "", nil
}
