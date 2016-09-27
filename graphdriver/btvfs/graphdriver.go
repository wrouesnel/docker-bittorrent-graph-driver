package btvfs

import (
	"io"
	"os"

	"github.com/wrouesnel/go.log"
	"github.com/docker/go-plugins-helpers/graphdriver"
	"path"
	"github.com/go-errors/errors"
	"fmt"
	"github.com/docker/docker/pkg/chrootarchive"
	"io/ioutil"
	"strings"
	"path/filepath"
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
	// Work dir's are created when a writeable layer is created. They are a full
	// standalone copy of the underlying layers.
	workDir = "work"

	// These are the standard docker whiteouts.
	whiteoutFormat = ".wh."
)

var (
	// CopyWithTar defines the copy method to use.
	CopyWithTar = chrootarchive.CopyWithTar
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

	return &this, nil
}

// Return a description of the driver
func (this *graphDriverBTVFS) String() string {
	return "btvfs"
}

// Generate an internal ID for a layer
//func (this *graphDriverBTVFS) generateLayerId() string {
//	return uuid.NewV4().String()
//}

func (this *graphDriverBTVFS) dir(id string) string {
	return path.Join(this.rootDirectory, id)
}

func (this *graphDriverBTVFS) Init(home string, options []string) error {
	// Ignore everything Docker sends us.
	// TODO: we could take an option to uniquely identify the docker client
	// here...
	log.Infoln("Docker Client connected! Home:", home, "Options:", options)

	return nil
}

// Create a new layer in the database. The layer is *unpopulated* at the point
// it's created, so no bittorrent metadata can be generated at this point.
// IDs at this level
func (this *graphDriverBTVFS) Create(id, parent string) error {
	log.Infoln("Create ID:", id,"Parent:", parent)

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

	//layerId := this.generateLayerId()

	if err := os.MkdirAll(path.Join(dir,linkDir), os.FileMode(0750)) ; err != nil {
		return errors.Wrap(err, 1)
	}

	if err := os.MkdirAll(path.Join(dir,diffDir), os.FileMode(0750)) ; err != nil {
		return errors.Wrap(err, 1)
	}

	if err := os.MkdirAll(path.Join(dir,parentDir), os.FileMode(0750)) ; err != nil {
		return errors.Wrap(err, 1)
	}

	if err := os.MkdirAll(path.Join(dir,refDir), os.FileMode(0750)) ; err != nil {
		return errors.Wrap(err, 1)
	}

	// If no parent directory then we are finished.
	if parent == "" {
		return nil
	}

	// There is a parent. Write a symlink named after the parent ID, pointing to
	// the parent directory.
	parentDir := this.dir(parent)

	// Check the parent exists and is a directory
	if pst, serr := os.Stat(parentDir) ; os.IsNotExist(serr) || !pst.IsDir() {
		return errors.New(fmt.Sprintf("parent layer does not exist or is not a correctly formatted directory: %v", parent))
	}

	// Make a symlink named after the ID to the parent.
	if err := os.Symlink(parentDir, path.Join(dir,parentDir,parent)); err != nil {
		return errors.Wrap(err, 1)
	}

	// Success
	retErr = nil

	return retErr
}

func (this *graphDriverBTVFS) CreateReadWrite(id, parent string) error {
	log.Infoln("CreateReadWrite ID:", id, "Parent:", parent)

	dir := this.dir(id)

	// Create a regular read-only layer
	if err := this.Create(id, parent) ; err != nil {
		return errors.Wrap(err,1)
	}

	// Make the working directory
	if err := os.MkdirAll(path.Join(dir,workDir), os.FileMode(0750)) ; err != nil {
		return errors.Wrap(err, 1)
	}

	// Get the parent layer dirs
	layers := []string{id}
	next := parent
	for next != "" {
		layers = append(layers, next)

		nextParentDir := path.Join(dir,next,parentDir)
		if st, err := os.Stat(nextParentDir) ; os.IsNotExist(err) || !st.IsDir() {
			return errors.New(fmt.Sprintf("layer has no parent directory, it may be damaged: %v", next))
		}

		files, rerr := ioutil.ReadDir(nextParentDir)
		if rerr != nil {
			return errors.Wrap(rerr, 1)
		}

		if len(files) > 1 {
			return errors.New(fmt.Sprintf("layer has more then 1 parent, it may be damaged: %v", next))
		}

		// The parent is just the symlink name
		if len(files) == 0 {
			next = files[0]
		} else {
			next = ""
		}
	}

	// Got list of directories from top->bottom.
	// Copy up the parent, delete whiteouts ala docker-squash to the working
	// directory.
	for i := len(layers)-1 ; i != 0; i-- {
		log.Debugln("Copying parent layer:", i-len(layers)-1 , path.Join(dir,next,diffDir))
		if err := CopyWithTar(path.Join(dir,next,diffDir), path.Join(dir,workDir)); err != nil {
			return errors.WrapPrefix(err, "Error copying up layers", 1)
		}

		// Delete white-outs
		if err := deleteWhiteouts(path.Join(dir,workDir)); err != nil {
			return errors.WrapPrefix(err, "Error deleting whiteouts while copying up", 1)
		}
	}

	return nil
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

func deleteWhiteouts(location string) error {
	return filepath.Walk(location, func(p string, info os.FileInfo, err error) error {
		if err != nil && !os.IsNotExist(err) {
			return err
		}

		if info == nil {
			return nil
		}

		name := info.Name()
		parent := filepath.Dir(p)
		// if start with whiteout
		if strings.Index(name, whiteoutFormat) == 0 {
			deletedFile := path.Join(parent, name[len(whiteoutFormat):len(name)])
			// remove deleted files
			if err := os.RemoveAll(deletedFile); err != nil {
				return err
			}
			// remove the whiteout itself
			if err := os.RemoveAll(p); err != nil {
				return err
			}
		}
		return nil
	})
}