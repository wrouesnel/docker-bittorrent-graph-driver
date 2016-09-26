package bittorrent_overlay2

import (
	"io"

	"github.com/docker/go-plugins-helpers/graphdriver"
)

// Implements the graphdriver/Driver interface
type bittorrentOverlay2GraphDriver struct {
}

func NewBitTorrentOverlay2GraphDriver() (*graphdriver.Driver, error) {
	return nil, &ErrUnimplemented{}
}

func (this *bittorrentOverlay2GraphDriver) Init(home string, options []string) error {
	return &ErrUnimplemented{}
}

func (this *bittorrentOverlay2GraphDriver) Create(id, parent string) error {
	return &ErrUnimplemented{}
}

func (this *bittorrentOverlay2GraphDriver) CreateReadWrite(id, parent string) error {
	return &ErrUnimplemented{}
}

func (this *bittorrentOverlay2GraphDriver) Remove(id string) error {
	return &ErrUnimplemented{}
}

func (this *bittorrentOverlay2GraphDriver) Get(id, mountLabel string) (string, error) {
	return "unimplemented", &ErrUnimplemented{}
}

func (this *bittorrentOverlay2GraphDriver) Put(id string) error {
	return &ErrUnimplemented{}
}

func (this *bittorrentOverlay2GraphDriver) Exists(id string) bool {
	return false
}

func (this *bittorrentOverlay2GraphDriver) Status() [][2]string {
	return [][2]string{}
}

func (this *bittorrentOverlay2GraphDriver) GetMetadata(id string) (map[string]string, error) {
	return nil, &ErrUnimplemented{}
}

func (this *bittorrentOverlay2GraphDriver) Cleanup() error {
	return &ErrUnimplemented{}
}

func (this *bittorrentOverlay2GraphDriver) Diff(id, parent string) io.ReadCloser {
	return nil
}

func (this *bittorrentOverlay2GraphDriver) Changes(id, parent string) ([]graphdriver.Change, error) {
	return []graphdriver.Change{}
}

func (this *bittorrentOverlay2GraphDriver) ApplyDiff(id, parent string, archive io.Reader) (int64, error) {
	return 0, &ErrUnimplemented{}
}

func (this *bittorrentOverlay2GraphDriver) DiffSize(id, parent string) (int64, error) {
	return 0, &ErrUnimplemented{}
}
