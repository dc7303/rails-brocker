package storage

import (
	"context"
	"log"

	"github.com/yorkie-team/yorkie/client"
	"github.com/yorkie-team/yorkie/pkg/document"
)

type Storage struct {
	addr string

	cli *client.Client
	doc *document.Document
}

func New(addr string) *Storage {
	return &Storage{
		addr: addr,
	}
}

func (s *Storage) Run() error {
	var err error
	s.cli, err = client.Dial(s.addr)
	if err != nil {
		return err
	}

	ctx := context.Background()
	err = s.cli.Activate(ctx)
	if err != nil {
		return err
	}

	s.doc = document.New("test-collection", "doc")
	if err = s.cli.Attach(ctx, s.doc); err != nil {
		return err
	}

	log.Println("Run storage")
	log.Println("Create document: test-collection$doc")

	return nil
}

func (s *Storage) Close(ctx context.Context) error {
	if err := s.cli.Detach(ctx, s.doc); err != nil {
		return err
	}

	if err := s.cli.Deactivate(ctx); err != nil {
		return err
	}

	log.Println("Close storage")

	return nil
}