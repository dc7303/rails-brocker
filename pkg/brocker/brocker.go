package brocker

import (
	"io"
	"log"
	"os/exec"
)

type Brocker struct {
	dir    string
	stdin  io.WriteCloser
	stdout io.ReadCloser
}

func New(dir string) *Brocker {
	return &Brocker{
		dir: dir,
	}
}

func (b *Brocker) Run() {
	cmd := exec.Command("dc")
	cmd.Dir = b.dir

	var err error
	b.stdin, err = cmd.StdinPipe()
	if err != nil {
		log.Fatal(err)
	}
	// defer b.stdin.Close()

	b.stdout, err = cmd.StdoutPipe()
	if err != nil {
		log.Fatal(err)
	}

	if err := cmd.Start(); err != nil {
		log.Fatal(err)
	}

	go func() {
		for {
			out := make([]byte, 1024)
			n, err := b.stdout.Read(out)
			if err != nil {
				log.Fatal(err)
			}

			log.Printf("output: %s\n", string(out[:n]))
		}
	}()
}

func (b *Brocker) Write(text string) {
	b.stdin.Write([]byte(text))
}

func (b *Brocker) Close() {
	if b.stdin != nil {
		b.stdin.Close()
	}
	if b.stdout != nil {
		b.stdout.Close()
	}
}
