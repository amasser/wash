package plugin

import (
	"context"
	"io"
	"os"
	"time"
)

// ==== Wash Protocols and Resources ====

// Entry is a basic named resource type.
type Entry interface {
	Name() string
}

// Cached is a resource that expects its method invocations to be cached.
// If CacheConfig() returns nil, a default config will be used.
type Cached interface {
	CacheConfig() *CacheConfig
}

// EntryBase implements Entry, making it easy to create new entries.
// EntryBase includes a default caching configuration.
type EntryBase struct {
	name        string
	cacheConfig *CacheConfig
}

// Name returns the entry's name.
func (e *EntryBase) Name() string { return e.name }

// CacheConfig returns the entry's cache config.
func (e *EntryBase) CacheConfig() *CacheConfig { return e.cacheConfig }

// MetadataMap maps keys to arbitrary structured data.
type MetadataMap = map[string]interface{}

// Resource is an entry that has metadata.
type Resource interface {
	Entry
	Metadata(context.Context) (MetadataMap, error)
}

// Group is an entry that can list its contents, an array of entries.
// It will be represented as a directory in the wash filesystem.
type Group interface {
	Entry
	List(context.Context) ([]Entry, error)
}

// Root is the root object of a plugin.
type Root interface {
	Group
	Init() error
}

// ExecOptions is a struct we can add new features to that must be serializable to JSON.
// Examples of potential features: user, privileged, tty, map of environment variables, string of stdin, timeout.
type ExecOptions struct {
	Stdin io.Reader
}

// ExecOutputChunk is a struct containing a chunk of the Exec'ed cmd's output.
type ExecOutputChunk struct {
	StreamID  int8
	Timestamp time.Time
	Data      string
	Err       error
}

// ExecResult is a struct that contains the result of invoking Execable#exec.
// Any of these fields can be nil.
type ExecResult struct {
	OutputCh   <-chan ExecOutputChunk
	ExitCodeCB func() (int, error)
}

// Execable is an entry that can have a command run on it.
type Execable interface {
	Entry
	Exec(ctx context.Context, cmd string, args []string, opts ExecOptions) (ExecResult, error)
}

// File is an entry that specifies filesystem attributes.
type File interface {
	Entry
	Attr() Attributes
}

// Pipe is an entry that returns a stream of updates.
type Pipe interface {
	Entry
	Stream(context.Context) (io.Reader, error)
}

// SizedReader returns a ReaderAt that can report its Size.
type SizedReader interface {
	io.ReaderAt
	Size() int64
}

// Readable is an entry that has a fixed amount of content we can read.
type Readable interface {
	Entry
	Open(context.Context) (SizedReader, error)
}

// Writable is an entry that we can write new data to.
type Writable interface {
	Entry
	Save(context.Context, io.Reader) error
}

// Attributes of resources.
type Attributes struct {
	Atime time.Time     `json:"atime"`
	Mtime time.Time     `json:"mtime"`
	Ctime time.Time     `json:"ctime"`
	Mode  os.FileMode   `json:"mode"`
	Size  uint64        `json:"size"`
	Valid time.Duration `json:"valid"`
}

// SizeUnknown can be used to denote that the size is unknown and should be queried from content.
const SizeUnknown = ^uint64(0)

// The Registry contains the core filesystem data.
// Plugins: maps plugin mount points to their implementations.
type Registry struct {
	Plugins map[string]Root
}
