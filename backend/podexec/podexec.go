package podexec

import (
	"io"
	"log"
	"os"

	"github.com/creack/pty"
	tfoclicmd "github.com/isaaguilar/terraform-operator-cli/cmd"
	"github.com/sorenisanerd/gotty/server"
)

type PodExecFactory struct{}
type PodExec struct {
	file   *os.File
	reader io.Reader
	writer io.Writer
}

func (f *PodExecFactory) Name() string {
	return "Pod Exec"
}

func (f *PodExecFactory) New(params map[string][]string, headers map[string][]string) (server.Slave, error) {
	// var s server.Slave

	return New()
}

// command string, argv []string, headers map[string][]string, options ...Option
func New() (*PodExec, error) {

	// command := "tfo"
	// cmd := exec.Command(command, "debug", "hello-tfo")

	// cmd.Env = append(os.Environ(), "TERM=xterm-256color")
	// pty, err := pty.Start(cmd)
	// if err != nil {
	// 	// todo close cmd?
	// 	return nil, errors.Wrapf(err, "failed to start command `%s`", command)
	// }
	// ptyClosed := make(chan struct{})
	// go func() {
	// 	defer func() {
	// 		pty.Close()
	// 		close(ptyClosed)
	// 	}()

	// 	cmd.Wait()
	// }()

	// return &PodExec{}, nil

	pty, tty, err := pty.Open()
	if err != nil {
		log.Fatal(err)
	}
	os.Stdin = tty
	os.Stderr = tty
	os.Stdout = tty

	// stop := make(chan bool)
	go func() {
		tfoclicmd.RemoteDebug("hello-tfo", tty)
		pty.Close()
		// stop <- true
	}()

	return &PodExec{
		file:   nil,
		reader: pty,
		writer: pty,
	}, nil

}

func (podExec *PodExec) Read(p []byte) (n int, err error) {
	if podExec.reader != nil {
		return podExec.reader.Read(p)
	}
	return podExec.file.Read(p)
	// return
}

func (podExec *PodExec) Write(p []byte) (n int, err error) {
	if podExec.writer != nil {
		return podExec.writer.Write(p)
	}
	return podExec.file.Write(p)
	// return
}

// WindowTitleVariables returns any values that can be used to fill out
// the title of a terminal.
func (p *PodExec) WindowTitleVariables() map[string]interface{} {
	return map[string]interface{}{
		"command":  "tfo",
		"hostname": "localhost",
	}

}

// ResizeTerminal sets a new size of the terminal.
func (p *PodExec) ResizeTerminal(columns int, rows int) error {
	return nil
}

func (p *PodExec) Close() error {
	log.Println("Going away...")
	return nil
}
