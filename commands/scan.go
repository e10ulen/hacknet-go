// commands/scan.go
package commands

import "github.com/e10ulen/hacknet-go/vfs"

// HandleScan は scan コマンドの処理
func HandleScan(args []string, vfs *vfs.VFS, log *[]string) (bool, string, error) {
	*log = append(*log, "スキャン中... 3つのノードを発見:")
	*log = append(*log, "  - 192.168.1.1 (proxy.local)")
	*log = append(*log, "  - 10.0.0.5 (target.corp)")
	*log = append(*log, "  - 172.16.0.1 (unknown)")
	*log = append(*log, "詳細は 'probe' または 'connect <IP>' で確認可能")
	// 状態変更なし → false, "", nil を返す
	return false, "", nil
}
