// main.go
package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	// Bubble Teaプログラムの起動
	// initialModel() は model.go で定義
	p := tea.NewProgram(initialModel(), tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		fmt.Printf("エラー: %v\n", err)
		os.Exit(1)
	}
}
