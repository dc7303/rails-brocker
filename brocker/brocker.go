package brocker

import (
	"context"
	"io"
	"log"
	"os/exec"
	"path"

	"github.com/dc7303/rails-brocker/brocker/storage"
)

type Brocker struct {
	dir     string
	stdin   io.WriteCloser
	stdout  io.ReadCloser
	stderr  io.ReadCloser
	storage *storage.Storage
}

func New(dir string) *Brocker {
	return &Brocker{
		dir: dir,
	}
}

func (b *Brocker) Run() error {
	log.Println("Run brocker")

	strg := storage.New("localhost:11101")
	if err := strg.Run(); err != nil {
		return err
	}
	ctx := context.Background()
	if err := strg.Close(ctx); err != nil {
		return err
	}

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
		return err
	}

	b.stdout, err = cmd.StdoutPipe()
	if err != nil {
		return err
	}

	b.stderr, err = cmd.StderrPipe()
	if err != nil {
		return err
	}

	reader := io.MultiReader(b.stdout, b.stderr)
	if err := cmd.Start(); err != nil {
		return err
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

	return nil
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
