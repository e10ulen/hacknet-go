// vfs/vfs.go
package vfs

import "fmt"

// File は仮想ファイルの構造
type File struct {
	Name     string
	Content  string
	IsDir    bool
	Children map[string]*File // ディレクトリの場合の子ノード
}

// VFS は仮想ファイルシステム全体
type VFS struct {
	Root    *File
	Current *File
}

// NewVFS は初期状態のVFSを作成
func NewVFS() *VFS {
	root := &File{
		Name:     "/",
		IsDir:    true,
		Children: make(map[string]*File),
	}

	// サンプルディレクトリとファイル
	home := &File{Name: "home", IsDir: true, Children: make(map[string]*File)}
	root.Children["home"] = home

	secret := &File{Name: "secret.txt", Content: "これは機密情報だ…\nパスワード: hunter2"}
	home.Children["secret.txt"] = secret

	logs := &File{Name: "logs.txt", Content: "アクセスログ: 192.168.1.100 -> 接続試行"}
	home.Children["logs.txt"] = logs

	virus := &File{Name: "virus.exe", Content: "[危険] 実行厳禁"}
	home.Children["virus.exe"] = virus

	return &VFS{
		Root:    root,
		Current: home, // 初期位置は /home
	}
}

// GetPath は現在のフルパスを返す
func (vfs *VFS) GetPath() string {
	if vfs.Current == vfs.Root {
		return "/"
	}
	// 簡易実装（本格的にはスタックで辿るべき）
	return "/home" // 今は固定で簡略化
}

// ChangeDir は cd コマンド相当
func (vfs *VFS) ChangeDir(path string) error {
	if path == ".." {
		if vfs.Current != vfs.Root {
			// 親に戻る（簡易実装）
			vfs.Current = vfs.Root
		}
		return nil
	}

	if dir, ok := vfs.Current.Children[path]; ok && dir.IsDir {
		vfs.Current = dir
		return nil
	}
	return fmt.Errorf("ディレクトリが見つかりません: %s", path)
}

// ListFiles は ls 相当
func (vfs *VFS) ListFiles() []string {
	var files []string
	for name, f := range vfs.Current.Children {
		if f.IsDir {
			files = append(files, name+"/")
		} else {
			files = append(files, name)
		}
	}
	return files
}

// ReadFile は cat 相当
func (vfs *VFS) ReadFile(name string) (string, error) {
	if f, ok := vfs.Current.Children[name]; ok && !f.IsDir {
		return f.Content, nil
	}
	return "", fmt.Errorf("ファイルが見つかりません: %s", name)
}

// RemoveFile は rm 相当（ログ削除でTrace減少用に使う想定）
func (vfs *VFS) RemoveFile(name string) error {
	if _, ok := vfs.Current.Children[name]; ok {
		delete(vfs.Current.Children, name)
		return nil
	}
	return fmt.Errorf("ファイルが見つかりません: %s", name)
}
