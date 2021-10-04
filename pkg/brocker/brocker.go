package brocker

import (
	"io"
	"log"
	"os/exec"
	"path"
)

type Brocker struct {
	dir    string
	stdin  io.WriteCloser
	stdout io.ReadCloser
	stderr io.ReadCloser
}

func New(dir string) *Brocker {
	return &Brocker{
		dir: dir,
	}
}

func (b *Brocker) Run() {
	cmd := exec.Command(
		path.Join(b.dir, "bin/bundle"),
		"exec",
		"rails",
		"c",
	)
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

	b.stderr, err = cmd.StderrPipe()
	if err != nil {
		log.Fatal(err)
	}

	reader := io.MultiReader(b.stdout, b.stderr)
	if err := cmd.Start(); err != nil {
		log.Fatal(err)
	}

	go func() {
		for {
			out := make([]byte, 2_147_483_647)
			n, err := reader.Read(out)
			if err != nil {
				log.Fatal(err)
			}

			log.Printf(string(out[:n]))
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
	if b.stderr != nil {
		b.stderr.Close()
	}
}
