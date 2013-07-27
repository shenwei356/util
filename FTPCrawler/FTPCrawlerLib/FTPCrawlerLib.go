package FTPCrawler

import (
	"errors"
	"fmt"
	"github.com/axgle/mahonia"
	ftp "github.com/shenwei356/goftp"
	"io"
	"log"
	"os"
	"path"
	"sync"
	"time"
)

type Site struct {
	Host    string
	Port    string
	User    string
	Passwd  string
	Path    string
	Charset string

	Conn *ftp.ServerConn
	Root *Node

	Dec mahonia.Decoder
	Enc mahonia.Encoder

	Queue      *DownloadQueue
	FilterTime time.Time
	Logger     *log.Logger
}

func NewSite(host, port, user, passwd, path1, charset string, logger *log.Logger) *Site {
	site := &Site{Host: host,
		Port:    port,
		User:    user,
		Passwd:  passwd,
		Path:    path1,
		Charset: charset,
		Logger:  logger}

	site.Dec = mahonia.NewDecoder(charset)
	site.Enc = mahonia.NewEncoder(charset)
	site.Path = site.Enc.ConvertString(site.Path)

	return site
}

func (site *Site) URL() string {
	return fmt.Sprintf("ftp://%s:%s@%s:%s%s", site.User, site.Passwd, site.Host, site.Port, site.Path)
}

func (site *Site) Login() error {
	conn, err := ftp.Connect(site.Host + ":" + site.Port)
	if err != nil {
		recover()
		return err
	}

	err = conn.Login(site.User, site.Passwd)
	if err != nil {
		recover()
		return err
	}

	site.Conn = conn
	return nil
}

func (site *Site) Logout() error {
	site.Root = nil

	if site.Conn == nil {
		recover()
		return errors.New("Not Login")
	}
	err := site.Conn.Quit()
	if err != nil {
		recover()
		return err
	}
	return nil
}

func (site *Site) DownloadWhileCrawling() error {
	site.Queue = &DownloadQueue{
		Queue: make(chan DownloadQueueItem),
		Site:  site,
		Sum:   0}

	go site.Queue.Start()

	fmt.Println("queue stared")
	err := site.crawlPathAddToDownloadQueue(site.Path)
	if err != nil {
		recover()
		return err
	}
	return nil
}

type DownloadQueueItem struct {
	ScrPath string
	SavPath string
}

type DownloadQueue struct {
	Queue chan DownloadQueueItem
	Site  *Site

	Sum int
	Wg  sync.WaitGroup
}

func (queue *DownloadQueue) AddDownloadQueueItem(item DownloadQueueItem) {
	queue.Queue <- item
}

func (queue *DownloadQueue) Start() {
	for {
		item := <-queue.Queue
		fmt.Printf("RETR %s\n", queue.Site.Dec.ConvertString(item.ScrPath))

		// Mkdir
		err := os.MkdirAll(path.Dir(item.SavPath), 0766)
		if err != nil {
			recover()
			fmt.Println(err)
		}

		//
		_, err = os.Stat(item.SavPath)
		if err != nil {
			recover()
		} else {
			queue.Wg.Done()
			continue
		}

		err = queue.Site.Retr(item.ScrPath, item.SavPath)
		if err != nil {
			recover()
			fmt.Println(err)
			queue.Site.Logger.Println(err)
			queue.Wg.Done()
			continue
		}

		queue.Sum++
		queue.Site.Logger.Printf("RETR %s\n", queue.Site.Dec.ConvertString(item.ScrPath))

		queue.Wg.Done()
	}
}

func (site *Site) crawlPathAddToDownloadQueue(path1 string) error {
	listDatas, err := site.Conn.List(path1)
	if err != nil {
		recover()
		return err
	}

	fmt.Printf("crawl %s\n", site.Dec.ConvertString(path1))
	for _, data := range listDatas {
		if data == nil {
			continue
		}
		if data.Name == "." || data.Name == ".." {
			continue
		}
		if data.TryCwd { // directory
			var p string
			if path1 == "/" {
				p = "/" + data.Name
			} else {
				p = path1 + "/" + data.Name
			}
			e := site.crawlPathAddToDownloadQueue(p)
			if e != nil {
				recover()
				return err
			}
		} else if data.TryRetr { // file
			if data.Mtime.After(site.FilterTime) {
				savePath := site.Dec.ConvertString(path.Join(site.Host, path1, data.Name))
				srcPath := path.Join(path1, data.Name)

				site.Queue.Wg.Add(1)
				site.Queue.AddDownloadQueueItem(DownloadQueueItem{srcPath, savePath})
			}

		} else if data.LinkDest != "" {
			//fmt.Printf("LINK: %s\n", data.Name)
			continue
		}
	}
	return nil
}

func (site *Site) Retr(srcPath, dstPath string) error {
	// CWD
	err := site.Conn.ChangeDir(path.Dir(srcPath))
	if err != nil {
		return err
	}

	// RETR
	fh, err := os.OpenFile(dstPath, os.O_CREATE|os.O_WRONLY, 0766)
	if err != nil {
		return err
	}
	reader, err := site.Conn.Retr(path.Base(srcPath))
	if err != nil {
		return err
	}
	_, err = io.Copy(fh, reader)
	if err != nil {
		return err
	}

	err = fh.Close()
	if err != nil {
		recover()
		return err
	}
	err = reader.Close()
	if err != nil {
		recover()
		return err
	}
	return nil
}

//===================================================================================
//===================================================================================
//===================================================================================
//===================================================================================
//===================================================================================
func (site *Site) Crawl() error {
	site.Root = &Node{&ftp.FTPListData{"", "", true, true, 0, ftp.UNKNOWN_MTIME_TYPE, time.Date(0, time.January, 1, 0, 0, 0, 0, time.Local), ftp.UNKNOWN_ID_TYPE, "", ""},
		nil,
		map[string]*Node{},
		""}
	return site.crawlPath(site.Path, site.Root)
}

func (site *Site) crawlPath(path1 string, node *Node) error {
	listDatas, err := site.Conn.List(path1)
	if err != nil {
		recover()
		return err
	}

	fmt.Printf("crawl %s\n", site.Dec.ConvertString(path1))
	for _, data := range listDatas {
		if data == nil {
			continue
		}
		if data.Name == "." || data.Name == ".." {
			continue
		}
		if data.TryCwd { // directory
			node.Children[data.Name] = &Node{data, node, map[string]*Node{}, ""}
			var p string
			if path1 == "/" {
				p = "/" + data.Name
			} else {
				p = path1 + "/" + data.Name
			}
			e := site.crawlPath(p, node.Children[data.Name])
			if e != nil {
				recover()
				return err
			}
		} else if data.TryRetr { // file
			node.Children[data.Name] = &Node{data, node, map[string]*Node{}, ""}
		} else if data.LinkDest != "" {
			//fmt.Printf("LINK: %s\n", data.Name)
			// TODO
			continue
		}
	}
	return nil
}

func (site *Site) FilterFiles(time time.Time) []*Node {
	var list []*Node
	list = site.filterFiles(site.Root, list, time)
	return list
}

func (site *Site) filterFiles(n *Node, list []*Node, time time.Time) []*Node {
	for _, v := range n.Children {
		if v.Entry.TryCwd {
			list = site.filterFiles(v, list, time)
		} else if v.Entry.TryRetr {
			if v.Entry.Mtime.After(time) {
				v.FullPath = v.Path()
				list = append(list, v)
			}
		} else if v.Entry.LinkDest != "" {
		}
	}
	return list
}

type Node struct {
	Entry    *ftp.FTPListData
	Parent   *Node
	Children map[string]*Node
	FullPath string
}

func (node *Node) Path() string {
	n := node
	path := n.Entry.Name
	for n.Parent != nil {
		n = n.Parent
		if n.Entry.Name == "/" {
			path = "/" + path
		} else {
			path = n.Entry.Name + "/" + path
		}
	}
	node.FullPath = path
	return path
}
