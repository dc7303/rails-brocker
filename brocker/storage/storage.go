package storage

import (
	"context"
	"log"

	"github.com/yorkie-team/yorkie/client"
	"github.com/yorkie-team/yorkie/pkg/document"
	"github.com/yorkie-team/yorkie/pkg/document/proxy"
)

type Storage struct {
	addr string

	cli *client.Client
	doc *document.Document

	logLen int
}

func New(addr string) *Storage {
	return &Storage{
		addr: addr,
	}
}

func (s *Storage) Run() error {
	log.Println("Run storage")

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

	err = s.doc.Update(func(root *proxy.ObjectProxy) error {
		root.SetNewText("log")
		return nil
	})
	if err != nil {
		return err
	}

	log.Println("Create document: test-collection$doc")

	return nil
}

func (s *Storage) Write(logText string) error {
	if err := s.doc.Update(func(root *proxy.ObjectProxy) error {
		text := root.GetText("log")
		text.Edit(s.logLen, s.logLen, logText)
		s.logLen += len(logText)
		return nil
	}); err != nil {
		return err
	}

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
