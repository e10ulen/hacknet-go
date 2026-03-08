// model.go
package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/e10ulen/hacknet-go/commands" // コマンドハンドラーのDispatch
	"github.com/e10ulen/hacknet-go/vfs"      // 仮想ファイルシステム
)

// スタイル定義（Hacknet風の緑基調）
var (
	titleStyle = lipgloss.NewStyle().
			MarginLeft(2).
			Foreground(lipgloss.Color("#00FF00"))

	logStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#00AA00"))

	inputStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFFFF")).
			Background(lipgloss.Color("#003300"))

	statusStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFF00"))
)

// model はゲーム全体の状態を保持
type model struct {
	log            []string            // ターミナルログの履歴（表示用）
	viewport       viewport.Model      // ログのスクロール表示領域
	input          textinput.Model     // コマンド入力欄
	trace          float64             // トレース値（0.0〜100.0）
	lastTick       time.Time           // 最後のtick時刻（Trace増加計算用）
	connected      bool                // 現在接続中か
	currentIP      string              // 接続中のIP（表示用）
	vfs            *vfs.VFS            // 仮想ファイルシステム
	missions       map[string]*Mission // ミッション一覧（ID -> Mission）
	activeMissions []string            // 進行中のミッションIDリスト
}

type Mission struct {
	ID          string // ミッションID（例: "entropy-001"）
	Title       string // 表示タイトル
	Description string // 詳細説明
	TargetIP    string // ターゲットIP
	TargetFile  string // 盗むファイル名など
	Completed   bool   // クリア済みか
	Reward      string // クリア報酬（ログ表示用）
}

// model.go に追加（model構造体の下あたり）

// GetMissionsList はミッション一覧を文字列スライスで返す（表示用）
func (m *model) GetMissionsList() []string {
	var list []string
	for id, miss := range m.missions {
		status := "未受注"
		if contains(m.activeMissions, id) {
			status = "進行中"
		}
		if miss.Completed {
			status = "完了"
		}
		list = append(list, fmt.Sprintf("[%s] %s - %s", status, miss.Title, miss.Description))
	}
	return list
}

// AcceptMission はミッションを受注（activeMissionsに追加）
func (m *model) AcceptMission(id string) bool {
	if _, ok := m.missions[id]; !ok {
		return false // 存在しないミッション
	}
	if !contains(m.activeMissions, id) {
		m.activeMissions = append(m.activeMissions, id)
		return true
	}
	return false // すでに受注済み
}

// CheckMissionComplete は cat コマンドなどでクリア判定（必要に応じて呼ぶ）
func (m *model) CheckMissionComplete(fileName string) bool {
	for _, mid := range m.activeMissions {
		miss := m.missions[mid]
		if !miss.Completed && miss.TargetFile == fileName && miss.TargetIP == m.currentIP {
			miss.Completed = true
			return true
		}
	}
	return false
}

// contains はスライスに要素が含まれるかチェック（ヘルパー）
func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

// initialModel は初期状態を作成
func initialModel() model {
	ti := textinput.New()
	ti.Placeholder = "コマンドを入力... (help で一覧)"
	ti.Focus()
	ti.CharLimit = 156
	ti.Width = 80

	vp := viewport.New(80, 20)
	vp.SetContent("")

	vfsInstance := vfs.NewVFS()

	m := model{
		log:      []string{"[Bit]: ようこそ、ハッカー。まずは 'help' でコマンド一覧を見ろ。"},
		viewport: vp,
		input:    ti,
		trace:    0.0,
		lastTick: time.Now(),
		vfs:      vfsInstance,
	}
	m.missions = map[string]*Mission{
		"entropy-001": {
			ID:          "entropy-001",
			Title:       "First Entropy Contract",
			Description: "Entropy Test Server (192.168.1.100) から secret.txt を盗め。プロキシ経由でTraceを抑えろ。",
			TargetIP:    "192.168.1.100",
			TargetFile:  "secret.txt",
			Completed:   false,
			Reward:      "TraceKill.exe 入手可能！",
		}, // 追加ミッション例（後で増やせる）
		"entropy-002": {
			ID:          "entropy-002",
			Title:       "Proxy Bypass Challenge",
			Description: "Bypass Proxy Server経由で target.corp に侵入せよ。",
			TargetIP:    "10.0.0.5",
			TargetFile:  "classified.doc",
			Completed:   false,
			Reward:      "BruteSSH.exe アップグレード",
		},
	}
	m.activeMissions = []string{"entropy-001"} // 最初は1つだけアクティブ
	m.updateViewport()
	return m
}

// Init は初回に呼ばれる（定期tickと入力カーソル点滅開始）
func (m model) Init() tea.Cmd {
	return tea.Batch(textinput.Blink, tickCmd())
}

// Update は各種メッセージに応じて状態を更新
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {

	case tea.KeyMsg:
		// 終了ショートカット
		switch msg.String() {
		case "ctrl+c", "esc":
			return m, tea.Quit
		}

		// 入力処理
		if m.input.Focused() {
			var cmd tea.Cmd
			m.input, cmd = m.input.Update(msg)
			cmds = append(cmds, cmd)

			if msg.Type == tea.KeyEnter {
				inputCmd := m.input.Value()
				if inputCmd != "" {
					m.log = append(m.log, fmt.Sprintf("> %s", inputCmd))
					m.executeCommand(inputCmd)
					m.input.Reset()
				}
			}
			// model.go の Update() の tea.KeyMsg case 内に追加（Enter処理の後）
			if msg.Type == tea.KeyEnter {
				inputCmd := m.input.Value()
				if inputCmd != "" {
					m.log = append(m.log, fmt.Sprintf("> %s", inputCmd))

					parts := strings.Fields(inputCmd)
					if len(parts) > 0 && parts[0] == "mission" {
						// missionコマンドはここで直接処理（modelからアクセス可能）
						m.handleMissionCommand(parts[1:])
					} else {
						m.executeCommand(inputCmd)
					}
					m.input.Reset()
				}
			}
		}

	case tea.WindowSizeMsg:
		// ウィンドウサイズ変更時のレイアウト調整
		m.viewport.Width = msg.Width - 4
		m.viewport.Height = msg.Height - 8
		m.input.Width = msg.Width - 10
		m.updateViewport()

	case tickMsg:
		// 1秒ごとのTrace増加処理
		now := time.Now()
		delta := now.Sub(m.lastTick).Seconds()
		m.lastTick = now

		if m.connected {
			increase := delta * (0.5 + float64(len(m.log)%10)/10)
			m.trace += increase

			if m.trace > 100 {
				m.log = append(m.log, "[TRACE 100%] 検知されました。追跡されました。終了。")
				return m, tea.Quit
			}
		}

		cmds = append(cmds, tickCmd())
	}

	m.updateViewport()
	return m, tea.Batch(cmds...)
}

// executeCommand は入力コマンドを解析し、対応ハンドラーを実行
func (m *model) executeCommand(cmd string) {
	parts := strings.Fields(cmd)
	if len(parts) == 0 {
		return
	}

	commandName := parts[0]
	args := parts[1:]

	// Dispatchでハンドラーを取得
	handler, ok := commands.Dispatch(commandName)
	if !ok {
		m.log = append(m.log, fmt.Sprintf("不明なコマンド: %s", cmd))
		m.log = append(m.log, "help と入力してコマンド一覧を表示してください")
		return
	}

	// ハンドラー実行（VFSとログを渡す）
	handler(args, m.vfs, &m.log)

	// 特殊処理例: logs.txt削除でTrace減少（Hacknet風バランス）
	if commandName == "rm" && len(args) > 0 && args[0] == "logs.txt" {
		m.trace = max(0.0, m.trace-8.0)
		m.log = append(m.log, "[TRACE] ログ削除成功 → 追跡リスクが軽減されました")
	}
}

// updateViewport はログをviewportに反映し、最下部に自動スクロール
func (m *model) updateViewport() {
	content := strings.Join(m.log, "\n")
	m.viewport.SetContent(content)
	m.viewport.GotoBottom()
}

// View は画面全体の描画
func (m model) View() string {
	s := ""

	s += titleStyle.Render("HACKNET TUI - Bitからの指令") + "\n\n"
	s += m.viewport.View() + "\n"

	status := fmt.Sprintf("Trace: %.1f%% | RAM: %d/%d MB | %s",
		m.trace,
		len(m.log)%512, 1024,
		ifConnected(m.connected, fmt.Sprintf("Connected: %s", m.currentIP), "Disconnected"))
	s += statusStyle.Render(status) + "\n"

	s += inputStyle.Render("> " + m.input.View())

	return s
}

// 補助関数
func ifConnected(cond bool, a, b string) string {
	if cond {
		return a
	}
	return b
}

func max(a, b float64) float64 {
	if a > b {
		return a
	}
	return b
}

// tickMsg と tickCmd（1秒ごとタイマー）
type tickMsg time.Time

func tickCmd() tea.Cmd {
	return tea.Tick(time.Second, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

// model.go に追加
func (m *model) handleMissionCommand(args []string) {
	if len(args) == 0 {
		m.log = append(m.log, "使い方: mission [list | accept <id> | status]")
		return
	}

	subCmd := args[0]
	switch subCmd {
	case "list":
		m.log = append(m.log, "利用可能なミッション一覧:")
		for _, line := range m.GetMissionsList() {
			m.log = append(m.log, line)
		}

	case "accept":
		if len(args) < 2 {
			m.log = append(m.log, "使い方: mission accept <ミッションID>")
			return
		}
		id := args[1]
		if m.AcceptMission(id) {
			m.log = append(m.log, fmt.Sprintf("ミッション受注成功: %s", id))
		} else {
			m.log = append(m.log, fmt.Sprintf("受注失敗: %s（存在しないか既に受注済み）", id))
		}

	case "status":
		m.log = append(m.log, "進行中のミッション:")
		for _, id := range m.activeMissions {
			miss := m.missions[id]
			status := "進行中"
			if miss.Completed {
				status = "完了"
			}
			m.log = append(m.log, fmt.Sprintf("  %s - %s", miss.Title, status))
		}

	default:
		m.log = append(m.log, "不明なサブコマンド: "+subCmd)
	}
}
