package manager

import (
	"io/ioutil"
	"os"
	"os/exec"
	"path"

	"github.com/micro-community/stream/engine/util"
)

//InstanceDesc for streaming
type InstanceDesc struct {
	Name    string
	Path    string
	Plugins []string
	Config  string
}

//Command to run
func (p *InstanceDesc) Command(name string, args ...string) (cmd *exec.Cmd) {
	cmd = exec.Command(name, args...)
	cmd.Dir = p.Path
	return
}

//CreateDir create new dir
func (p *InstanceDesc) CreateDir(sse *util.SSE, clearDir bool) (err error) {
	if clearDir {
		err = os.RemoveAll(p.Path)
	}
	if err = os.MkdirAll(p.Path, 0666); err != nil {
		return
	}
	sse.WriteEvent("step", []byte("2:目录创建成功！"))
	if err = ioutil.WriteFile(path.Join(p.Path, "config.toml"), []byte(p.Config), 0666); err != nil {
		return
	}

	return nil
}
