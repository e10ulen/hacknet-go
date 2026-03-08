// commands/rm.go
package commands

import (
	"fmt"

	"github.com/e10ulen/hacknet-go/vfs"
)

// HandleRm は rm コマンドの処理（ログ削除でTrace減少の想定）
func HandleRm(args []string, vfs *vfs.VFS, log *[]string) (bool, string, error) {
	if len(args) < 1 {
		*log = append(*log, "使い方: rm <ファイル>")
		return false, "", nil
	}

	err := vfs.RemoveFile(args[0])
	if err != nil {
		*log = append(*log, fmt.Sprintf("エラー: %v", err))
		return false, "", nil
	}

	*log = append(*log, fmt.Sprintf("削除しました: %s", args[0]))
	// ログ削除ならTraceを少し減らす（ゲームバランス調整用）
	// m.trace -= 5.0  ← model側で処理するのでここではログだけ
	// 状態変更なし → false, "", nil を返す
	return false, "", nil
}
