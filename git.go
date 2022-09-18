package main

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/wailsapp/wails"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"
)

type Git struct {
	RepoPath string
	runtime  *wails.Runtime
}

func NewGit() (*Git, error) {
	result := &Git{}
	return result, nil
}

func (g *Git) WailsInit(runtime *wails.Runtime) error {
	g.runtime = runtime
	return nil
}

func (g *Git) SetRepoPath(RepoPath string) {
	g.RepoPath = RepoPath
	if done != nil && !isClosed(done) {
		close(done)
	}

	// 设置完仓库目录后开始监听文件变化
	go g.Watch()
}

func isClosed(ch <-chan bool) bool {
	select {
	case <-ch:
		return true
	default:
	}

	return false
}

func (g *Git) Files() []*File {
	var files []*File
	root := g.RepoPath

	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		// 路径不存在
		if err != nil {
			return nil
		}
		// 跳过 .git 目录
		if info.IsDir() && info.Name() == ".git" {
			return filepath.SkipDir
		}
		// 只扫描文件，如果是目录则跳过
		if info.IsDir() {
			return nil
		}
		// 隐藏 .swp 文件
		if strings.HasSuffix(info.Name(), ".swp") {
			return nil
		}

		f := &File{}
		f.Path = path
		f.Name = info.Name()
		f.ModTime = info.ModTime()
		content, _ := os.ReadFile(path)
		f.Content = string(content)

		files = append(files, f)

		return nil
	})
	if err != nil {
		panic(err)
	}

	sort.SliceStable(files, sortFileByModTime(files))

	return files
}

func (g *Git) Index() *File {
	path := g.RepoPath + "/.git/index"
	info, err := os.Stat(path)
	if err != nil {
		return nil
	}

	f := &File{}
	f.Path = path
	f.Name = info.Name()
	f.ModTime = info.ModTime()
	f.Content = g.lsFiles()

	return f
}

func (g *Git) lsFiles() string {
	err := os.Chdir(g.RepoPath)
	if err != nil {
		log.Fatal(err)
	}
	out, err := exec.Command("git", "ls-files", "--stage").Output()
	if err != nil {
		log.Fatal(err)
	}
	str := string(out)
	str = strings.TrimSuffix(str, "\n")

	return str
}

func (g *Git) lsFileEntries() []*Entry {
	err := os.Chdir(g.RepoPath)
	if err != nil {
		log.Fatal(err)
	}
	out, err := exec.Command("git", "ls-files", "--stage").Output()
	if err != nil {
		log.Fatal(err)
	}
	str := string(out)
	str = strings.TrimSuffix(str, "\n")
	lines := strings.Split(str, "\n")

	entries := make([]*Entry, 0)
	for _, line := range lines {
		p := strings.Split(line, "\t")
		q := strings.Split(p[0], " ")
		e := &Entry{}
		e.Filename = p[1]
		e.Mode = q[0]
		e.Hash = q[1]
		e.Slot, _ = strconv.Atoi(q[2])
		entries = append(entries, e)
	}

	return entries
}

func (g *Git) HEAD() *File {
	path := g.RepoPath + "/.git/HEAD"
	info, err := os.Stat(path)
	if err != nil {
		return nil
	}

	f := &File{}
	f.Path = path
	f.Name = info.Name()
	f.ModTime = info.ModTime()
	content, _ := os.ReadFile(path)
	f.Content = string(content)

	return f
}

func (g *Git) Heads() []*File {
	var files []*File
	root := g.RepoPath + "/.git/refs/heads/"

	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		if !info.IsDir() {
			f := &File{}
			f.Path = path
			f.Name = info.Name()
			f.ModTime = info.ModTime()
			content, _ := os.ReadFile(path)
			f.Content = string(content)

			files = append(files, f)
		}
		return nil
	})
	if err != nil {
		panic(err)
	}

	sort.SliceStable(files, sortFileByModTime(files))

	return files
}

func sortFileByModTime(files []*File) func(i int, j int) bool {
	return func(i, j int) bool {
		return files[i].ModTime.Before(files[j].ModTime)
	}
}

func (g *Git) Objects() []*Object {
	var objects []*Object
	root := g.RepoPath + "/.git/objects"

	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		if !info.IsDir() {
			pars := strings.Split(path, "/")
			size := len(pars)
			hash := pars[size-2] + pars[size-1]

			o := &Object{}
			o.Hash = hash
			o.Type = g.catFile(g.RepoPath, hash, "-t")
			o.Content = g.catFile(g.RepoPath, hash, "-p")
			o.ModeTime = info.ModTime()

			objects = append(objects, o)
		}
		return nil
	})
	if err != nil {
		panic(err)
	}

	sort.SliceStable(objects, func(i, j int) bool {
		return objects[i].ModeTime.Before(objects[j].ModeTime)
	})

	return objects
}

func (g *Git) catFile(repo string, hash string, option string) string {
	err := os.Chdir(repo)
	if err != nil {
		log.Fatal(err)
	}
	out, err := exec.Command("git", "cat-file", option, hash).Output()
	if err != nil {
		log.Fatal(err)
	}
	str := string(out)
	str = strings.TrimSuffix(str, "\n")
	return str
}

type File struct {
	Path    string
	Name    string
	Content string
	ModTime time.Time
}

type Object struct {
	Hash     string
	Type     string
	Content  string
	ModeTime time.Time
}

type Entry struct {
	Mode     string
	Hash     string
	Slot     int
	Filename string
}

var watcher *fsnotify.Watcher
var done chan bool

func (g *Git) Watch() {

	// creates a new file watcher
	watcher, _ = fsnotify.NewWatcher()
	defer watcher.Close()

	// starting at the root of the project, walk each file/directory searching for
	// directories
	root := g.RepoPath
	if err := filepath.Walk(root, watchDir); err != nil {
		fmt.Println("ERROR", err)
	}

	done = make(chan bool)

	go func() {
		for {
			select {
			// watch for events
			case event := <-watcher.Events:
				// TODO: 为什么会监听到文件名为空，Op为0的操作
				if event.Name == "" {
					break
				}
				// 读取文件会修改文件的最后访问时间
				if event.Op&fsnotify.Chmod == fsnotify.Chmod {
					break
				}

				fmt.Printf("EVENT! %#v\n", event)

				if event.Op&fsnotify.Create == fsnotify.Create {
					fi, err := os.Stat(event.Name)
					if err != nil {
						// 如果文件已经不存在会报错误（比如 .swp）
						break
					}
					// 如果新建了文件夹，需要对其监听
					if fi.Mode().IsDir() {
						_ = watcher.Add(event.Name)
					} else {
						g.runtime.Events.Emit("file_changed")
					}
				} else {
					g.runtime.Events.Emit("file_changed")
				}
			// watch for errors
			case err := <-watcher.Errors:
				fmt.Println("ERROR", err)
			}
		}
	}()

	<-done
}

// watchDir gets run as a walk func, searching for directories to add watchers to
func watchDir(path string, fi os.FileInfo, err error) error {
	if err != nil {
		return nil
	}

	// 忽略 .git 目录可以解决 too many open files 问题，但是无法监听 git add、commit 命令导致的内部文件变化
	//if strings.Index(path, "/.git") > -1 {
	//	return nil
	//}

	// since fsnotify can watch all the files in a directory, watchers only need
	// to be added to each nested directory
	if fi.Mode().IsDir() {
		return watcher.Add(path)
	}

	return nil
}
