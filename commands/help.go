// 例: commands/help.go の修正版
package commands

import "github.com/e10ulen/hacknet-go/vfs"

// HandleHelp は help コマンド
func HandleHelp(args []string, vfs *vfs.VFS, log *[]string) (bool, string, error) {
	helpText := []string{
		"利用可能なコマンド一覧:",
		"  help          - このヘルプを表示",
		"  scan          - 周囲のノードをスキャン",
		"  connect <IP>  - 指定IPに接続",
		"  exit          - 現在の接続を切断",
		"  ls            - 現在のディレクトリの内容を表示",
		"  cd <dir>      - ディレクトリ移動",
		"  cat <file>    - ファイル内容を表示",
		"  rm <file>     - ファイルを削除",
		"  mission list     - ミッション一覧表示",
		"  mission accept <id> - ミッション受注",
		"  mission status   - 進行状況確認",
	}

	for _, line := range helpText {
		*log = append(*log, line)
	}

	// 状態変更なし → false, "", nil を返す
	return false, "", nil
}
