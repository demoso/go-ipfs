package iface

import (
	"context"
	"errors"
	"io"

	options "github.com/ipfs/go-ipfs/core/coreapi/interface/options"

	cid "gx/ipfs/QmNp85zy9RLrQ5oQD4hPyS39ezrrXpcaa7R4Y9kxdWQLLQ/go-cid"
	ipld "gx/ipfs/QmPN7cwmpcc4DWXb4KTB9dNAJgjuPY69h3npsMfhRrQL9c/go-ipld-format"
)

type Path interface {
	String() string
	Cid() *cid.Cid
	Root() *cid.Cid
	Resolved() bool
}

// TODO: should we really copy these?
//       if we didn't, godoc would generate nice links straight to go-ipld-format
type Node ipld.Node
type Link ipld.Link

type Reader interface {
	io.ReadSeeker
	io.Closer
}

type CoreAPI interface {
	Unixfs() UnixfsAPI
	Dag() DagAPI

	ResolvePath(context.Context, Path) (Path, error)
	ResolveNode(context.Context, Path) (Node, error)
}

type UnixfsAPI interface {
	Add(context.Context, io.Reader) (Path, error)
	Cat(context.Context, Path) (Reader, error)
	Ls(context.Context, Path) ([]*Link, error)
}

// DagAPI specifies the interface to IPLD
type DagAPI interface {
	// Put inserts data using specified format and input encoding.
	// If format is not specified (nil), default dag-cbor/sha256 is used
	Put(ctx context.Context, src io.Reader, opts ...options.DagPutOption) ([]Node, error)

	// WithInputEnc is an option for Put which specifies the input encoding of the
	// data. Default is "json", most formats/codecs support "raw"
	WithInputEnc(enc string) options.DagPutOption

	// WithCodec is an option for Put which specifies the multicodec to use to
	// serialize the object. Default is cid.DagCBOR (0x71)
	WithCodec(codec uint64) options.DagPutOption

	// WithHash is an option for Put which specifies the multihash settings to use
	// when hashing the object. Default is based on the codec used
	// (mh.SHA2_256 (0x12) for DagCBOR). If mhLen is set to -1, default length for
	// the hash will be used
	WithHash(mhType uint64, mhLen int) options.DagPutOption

	// Get attempts to resolve and get the node specified by the path
	Get(ctx context.Context, path Path) (Node, error)

	// Tree returns list of paths within a node specified by the path.
	Tree(ctx context.Context, path Path, opts ...options.DagTreeOption) ([]Path, error)

	// WithDepth is an option for Tree which specifies maximum depth of the
	// returned tree. Default is -1 (no depth limit)
	WithDepth(depth int) options.DagTreeOption
}

// type ObjectAPI interface {
// 	New() (cid.Cid, Object)
// 	Get(string) (Object, error)
// 	Links(string) ([]*Link, error)
// 	Data(string) (Reader, error)
// 	Stat(string) (ObjectStat, error)
// 	Put(Object) (cid.Cid, error)
// 	SetData(string, Reader) (cid.Cid, error)
// 	AppendData(string, Data) (cid.Cid, error)
// 	AddLink(string, string, string) (cid.Cid, error)
// 	RmLink(string, string) (cid.Cid, error)
// }

// type ObjectStat struct {
// 	Cid            cid.Cid
// 	NumLinks       int
// 	BlockSize      int
// 	LinksSize      int
// 	DataSize       int
// 	CumulativeSize int
// }

var ErrIsDir = errors.New("object is a directory")
var ErrOffline = errors.New("can't resolve, ipfs node is offline")
