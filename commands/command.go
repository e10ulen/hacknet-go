// commands/command.go
package commands

import (
	"github.com/e10ulen/hacknet-go/vfs"
)

// CommandHandler はコマンド処理の関数型
// 戻り値で model に反映すべき状態変更を伝える
type CommandHandler func(args []string, vfs *vfs.VFS, log *[]string) (changedConnected bool, newIP string, err error)

// Dispatch はコマンド名からハンドラーを返す
func Dispatch(cmd string) (CommandHandler, bool) {
	handlers := map[string]CommandHandler{
		"help":    HandleHelp,
		"scan":    HandleScan,
		"connect": HandleConnect,
		"exit":    HandleExit,
		"ls":      HandleLs,
		"cd":      HandleCd,
		"cat":     HandleCat,
		"rm":      HandleRm,
		"mission": HandleMission,
	}

	handler, ok := handlers[cmd]
	return handler, ok
}
