package btvfs

import (
	"io"
	"os"

	"github.com/wrouesnel/go.log"
	"github.com/docker/go-plugins-helpers/graphdriver"
	"path"
	"github.com/go-errors/errors"
	"github.com/satori/go.uuid"
	"io/ioutil"
)

const (
	// The link dir contains a symlink to all the child layer dirs which reference
	// this diff. When emptied, it is safe to delete this layer.
	linkDir = "link"
	// The diff dir contains a flat filesystem and whiteout files for this layer
	// changes over the parent layer.
	diffDir = "diff"
	// The parent dir contains a symlink to the parent layer dir
	parentDir = "parent"
	// The ref dir contains a pid file for each docker process which
	// currently holds a reference to this layer. When it is empty, and linkDir
	// is empty, the layer is free to be GC'd.
	refDir = "refs"

	// TODO: BitTorrent primitives
	// The basic idea
)

// Implements the graphdriver/Driver interface
type graphDriverBTVFS struct {
	rootDirectory string
}

func NewBitTorrentVFSGraphDriver(graphStoragePath string) (*graphdriver.Driver, error) {
	if _, err := os.Stat(graphStoragePath); os.IsNotExist(err) {
		// Create the directory
		cerr := os.MkdirAll(graphStoragePath, os.FileMode(0750))
		if cerr != nil {
			return errors.Wrap(cerr, 1 )
		}
	}

	this := graphDriverBTVFS{
		rootDirectory: graphStoragePath,
	}

	// TODO: check kernel version
	// TODO: check backing FS

	return &this, &ErrUnimplemented{}
}

// Generate an internal ID for a layer
func (this *graphDriverBTVFS) generateLayerId() string {
	return uuid.NewV4().String()
}

func (this *graphDriverBTVFS) dir(id string) string {
	return path.Join(this.rootDirectory, id)
}

func (this *graphDriverBTVFS) Init(home string, options []string) error {
	// Ignore everything Docker sends us.
	// TODO: we could take an option to uniquely identify the docker client
	// here...
	log.Debugln("Docker Client connected! Home:", home, "Options:", options)

	return nil
}

// Create a new layer in the database. The layer is *unpopulated* at the point
// it's created, so no bittorrent metadata can be generated at this point.
// IDs at this level
func (this *graphDriverBTVFS) Create(id, parent string) error {
	log.Debugln("Create ID:", id,"Parent:", parent)

	// The returned error
	var retErr error

	// Generate the full dir name
	dir := this.dir(id)

	// TODO: sensible file mode. UID/GID mapping?
	if err := os.MkdirAll(dir, os.FileMode(0777)); err != nil {
		return errors.Wrap(err, 1)
	}

	defer func() {
		// Clean up on failure
		if retErr != nil {
			os.RemoveAll(dir)
		}
	}()

	layerId := this.generateLayerId()

	// Write layer id to the link file
	if err := ioutil.WriteFile(path.Join(dir, "link"), []byte(layerId), 0644); err != nil {
		return err
	}

	// if no parent directory, done
	if parent == "" {
		return nil
	}



	return nil
}

func (this *graphDriverBTVFS) CreateReadWrite(id, parent string) error {
	log.Debugln("CreateReadWrite ID:", id, "Parent:", parent)
	return &ErrUnimplemented{}
}

func (this *graphDriverBTVFS) Remove(id string) error {
	return &ErrUnimplemented{}
}

func (this *graphDriverBTVFS) Get(id, mountLabel string) (string, error) {
	return "unimplemented", &ErrUnimplemented{}
}

func (this *graphDriverBTVFS) Put(id string) error {
	return &ErrUnimplemented{}
}

func (this *graphDriverBTVFS) Exists(id string) bool {
	return false
}

func (this *graphDriverBTVFS) Status() [][2]string {
	return [][2]string{}
}

func (this *graphDriverBTVFS) GetMetadata(id string) (map[string]string, error) {
	return nil, &ErrUnimplemented{}
}

func (this *graphDriverBTVFS) Cleanup() error {
	return &ErrUnimplemented{}
}

func (this *graphDriverBTVFS) Diff(id, parent string) io.ReadCloser {
	return nil
}

func (this *graphDriverBTVFS) Changes(id, parent string) ([]graphdriver.Change, error) {
	return []graphdriver.Change{}
}

func (this *graphDriverBTVFS) ApplyDiff(id, parent string, archive io.Reader) (int64, error) {
	return 0, &ErrUnimplemented{}
}

func (this *graphDriverBTVFS) DiffSize(id, parent string) (int64, error) {
	return 0, &ErrUnimplemented{}
}
