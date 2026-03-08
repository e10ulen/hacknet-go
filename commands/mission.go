// commands/mission.go
package commands

import (
	"fmt"

	"github.com/e10ulen/hacknet-go/vfs"
)

// HandleMission は mission コマンドのサブコマンド処理
// modelにアクセスできないので、必要に応じて log に「modelメソッドを呼ぶ」旨を書き込むだけ
// 実際の処理は model側で実行される形にする
func HandleMission(args []string, vfs *vfs.VFS, log *[]string) (bool, string, error) {
	if len(args) == 0 {
		*log = append(*log, "使い方: mission [list | accept <id> | status]")
		return false, "", nil
	}

	subCmd := args[0]
	switch subCmd {
	case "list":
		*log = append(*log, "利用可能なミッション一覧:")
		// ここでは model.GetMissionsList() を呼べないので、
		// 「modelにアクセスできない」旨をログに書いておく（本当はmodel側で処理）
		*log = append(*log, "[システム] ミッション一覧は model.GetMissionsList() で取得可能ですが、")
		*log = append(*log, "          commandsパッケージからは直接アクセスできません。")
		*log = append(*log, "          mission list を model.Update 内で処理するように設計変更が必要です。")

	case "accept":
		if len(args) < 2 {
			*log = append(*log, "使い方: mission accept <ミッションID>")
			return false, "", nil
		}
		id := args[1]
		*log = append(*log, fmt.Sprintf("[受注試行] ミッション %s を受注中...", id))
		*log = append(*log, "[システム] model.AcceptMission(id) を model側で呼ぶ必要があります。")

	case "status":
		*log = append(*log, "進行中のミッション:")
		*log = append(*log, "[システム] model.activeMissions を model側で表示してください。")

	default:
		*log = append(*log, "不明なサブコマンド: "+subCmd)
	}

	return false, "", nil
}
