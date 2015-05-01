package otaru

import (
	"fmt"
	"syscall"

	"github.com/nyaxt/otaru/intn"
)

const (
	ENOENT  = syscall.Errno(syscall.ENOENT)
	ENOTDIR = syscall.Errno(syscall.ENOTDIR)
	EEXIST  = syscall.Errno(syscall.EEXIST)
)

type FileChunk struct {
	Offset   int64
	Length   int64
	BlobPath string
}

func (fc FileChunk) Left() int64 {
	return fc.Offset
}

func (fc FileChunk) Right() int64 {
	return fc.Offset + fc.Length
}

type INodeID uint32

type INodeType int

const (
	FileNodeT = iota
	DirNodeT
	// SymlinkNode
)

type INode interface {
	ID() INodeID
	Type() INodeType
}

type INodeCommon struct {
	INodeID
	INodeType

	// OrigPath contains filepath passed to first create and does not necessary follow "rename" operations.
	// To be used for recovery/debug purposes only
	OrigPath string
}

func (n INodeCommon) ID() INodeID {
	return n.INodeID
}

func (n INodeCommon) Type() INodeType {
	return n.INodeType
}

type FileNode struct {
	INodeCommon
	Size   int64
	Chunks []FileChunk
}

func NewFileNode(db *INodeDB, origpath string) *FileNode {
	id := db.GenerateNewID()
	fn := &FileNode{
		INodeCommon: INodeCommon{
			INodeID:   id,
			INodeType: FileNodeT,
			OrigPath:  origpath,
		},
		Size: 0,
	}
	db.PutMustSucceed(fn)
	return fn
}

type DirNode struct {
	INodeCommon
	Entries map[string]INodeID
}

func NewDirNode(db *INodeDB, origpath string) *DirNode {
	id := db.GenerateNewID()
	dn := &DirNode{
		INodeCommon: INodeCommon{
			INodeID:   id,
			INodeType: DirNodeT,
			OrigPath:  origpath,
		},
		Entries: make(map[string]INodeID),
	}
	db.PutMustSucceed(dn)
	return dn
}

const (
	FileWriteCacheMaxPatches         = 32
	FileWriteCacheMaxPatchContentLen = 256 * 1024
)

type FileWriteCache struct {
	ps intn.Patches
}

func NewFileWriteCache() *FileWriteCache {
	return &FileWriteCache{ps: intn.NewPatches()}
}

func (wc *FileWriteCache) PWrite(offset int64, p []byte) error {
	newp := intn.Patch{Offset: offset, P: p}
	wc.ps = wc.ps.Merge(newp)
	return nil
}

func (wc *FileWriteCache) PReadThrough(offset int64, p []byte, r PReader) error {
	nr := int64(len(p))
	remo := offset
	remp := p

	for _, patch := range wc.ps {
		if nr <= 0 {
			return nil
		}

		if remo > patch.Right() {
			continue
		}

		if remo < patch.Left() {
			fallbackLen := Int64Min(nr, patch.Left()-remo)

			if err := r.PRead(remo, remp[:fallbackLen]); err != nil {
				return err
			}

			remp = remp[fallbackLen:]
			nr -= fallbackLen
			remo += fallbackLen
		}

		if nr <= 0 {
			return nil
		}

		applyOffset := remo - patch.Offset
		applyLen := Int64Min(int64(len(patch.P))-applyOffset, nr)
		copy(remp[:applyLen], patch.P[applyOffset:])

		remp = remp[applyLen:]
		nr -= applyLen
		remo += applyLen
	}

	if err := r.PRead(remo, remp); err != nil {
		return err
	}
	return nil
}

func (wc *FileWriteCache) ContentLen() int64 {
	l := int64(0)
	for _, p := range wc.ps {
		l += int64(len(p.P))
	}
	return l
}

func (wc *FileWriteCache) NeedsFlush() bool {
	if len(wc.ps) > FileWriteCacheMaxPatches {
		return true
	}
	if wc.ContentLen() > FileWriteCacheMaxPatchContentLen {
		return true
	}

	return false
}

func (wc *FileWriteCache) Flush(bh BlobHandle) error {
	for _, p := range wc.ps {
		if err := bh.PWrite(p.Offset, p.P); err != nil {
			return err
		}
	}
	wc.ps = wc.ps[:0]

	return nil
}

func (wc *FileWriteCache) Right() int64 {
	if len(wc.ps) == 0 {
		return 0
	}

	return wc.ps[0].Right()
}

func (wc *FileWriteCache) Truncate(size int64) {
	wc.ps = wc.ps.Truncate(size)
}

type FileSystem struct {
	*INodeDB
	lastID INodeID

	bs RandomAccessBlobStore
	c  Cipher

	newChunkedFileIO func(bs RandomAccessBlobStore, fn *FileNode, c Cipher) BlobHandle

	wcmap map[INodeID]*FileWriteCache
}

func NewFileSystem(bs RandomAccessBlobStore, c Cipher) *FileSystem {
	fs := &FileSystem{
		INodeDB: NewINodeDB(),
		lastID:  0,

		bs: bs,
		c:  c,

		newChunkedFileIO: func(bs RandomAccessBlobStore, fn *FileNode, c Cipher) BlobHandle {
			return NewChunkedFileIO(bs, fn, c)
		},

		wcmap: make(map[INodeID]*FileWriteCache),
	}

	rootdir := NewDirNode(fs.INodeDB, "/")
	if rootdir.ID() != 1 {
		panic("rootdir must have INodeID 1")
	}

	return fs
}

func (fs *FileSystem) getOrCreateFileWriteCache(id INodeID) *FileWriteCache {
	wc := fs.wcmap[id]
	if wc == nil {
		wc = NewFileWriteCache()
		fs.wcmap[id] = wc
	}
	return wc
}

func (fs *FileSystem) OverrideNewChunkedFileIOForTesting(newChunkedFileIO func(RandomAccessBlobStore, *FileNode, Cipher) BlobHandle) {
	fs.newChunkedFileIO = newChunkedFileIO
}

type DirHandle struct {
	fs *FileSystem
	n  *DirNode
}

func (fs *FileSystem) OpenDir(id INodeID) (*DirHandle, error) {
	node := fs.INodeDB.Get(id)
	if node == nil {
		return nil, ENOENT
	}
	if node.Type() != DirNodeT {
		return nil, ENOTDIR
	}
	dirnode := node.(*DirNode)

	h := &DirHandle{fs: fs, n: dirnode}
	return h, nil
}

func (dh *DirHandle) FileSystem() *FileSystem {
	return dh.fs
}

func (dh *DirHandle) ID() INodeID {
	return dh.n.ID()
}

func (dh *DirHandle) Entries() map[string]INodeID {
	return dh.n.Entries
}

// FIXME
func (dh *DirHandle) Path() string {
	return ""
}

func (dh *DirHandle) CreateFile(name string) (INodeID, error) {
	_, ok := dh.n.Entries[name]
	if ok {
		return 0, EEXIST
	}

	fullorigpath := fmt.Sprintf("%s/%s", dh.Path(), name)
	n := NewFileNode(dh.fs.INodeDB, fullorigpath)

	id := n.ID()
	dh.n.Entries[name] = id
	return id, nil
}

// FIXME: Multiple FileHandle may exist for same file at once. Support it!
type FileHandle struct {
	fs   *FileSystem
	n    *FileNode
	wc   *FileWriteCache
	cfio BlobHandle
}

func (fs *FileSystem) openFileNode(n *FileNode) (*FileHandle, error) {
	wc := fs.getOrCreateFileWriteCache(n.ID())
	cfio := fs.newChunkedFileIO(fs.bs, n, fs.c)
	h := &FileHandle{fs: fs, n: n, wc: wc, cfio: cfio}
	return h, nil
}

func (fs *FileSystem) OpenFile(id INodeID) (*FileHandle, error) {
	node := fs.INodeDB.Get(id)
	if node == nil {
		return nil, ENOENT
	}
	if node.Type() != FileNodeT {
		return nil, ENOTDIR
	}
	filenode := node.(*FileNode)

	h, err := fs.openFileNode(filenode)
	if err != nil {
		return nil, err
	}
	return h, nil
}

type Attr struct {
	INodeID
	INodeType
	Size int64
}

func (fs *FileSystem) Attr(id INodeID) (Attr, error) {
	n := fs.INodeDB.Get(id)
	if n == nil {
		return Attr{}, ENOENT
	}

	size := int64(0)
	if fn, ok := n.(*FileNode); ok {
		size = fn.Size
	}

	a := Attr{
		INodeID:   n.ID(),
		INodeType: n.Type(),
		Size:      size,
	}
	return a, nil
}

func (fs *FileSystem) IsDir(id INodeID) (bool, error) {
	n := fs.INodeDB.Get(id)
	if n == nil {
		return false, ENOENT
	}

	return n.Type() == DirNodeT, nil
}

func (h *FileHandle) ID() INodeID {
	return h.n.ID()
}

func (h *FileHandle) PWrite(offset int64, p []byte) error {
	if err := h.wc.PWrite(offset, p); err != nil {
		return err
	}

	if h.wc.NeedsFlush() {
		if err := h.wc.Flush(h.cfio); err != nil {
			return err
		}
	}

	right := offset + int64(len(p))
	if right > h.n.Size {
		h.n.Size = right
	}

	return nil
}

func (h *FileHandle) PRead(offset int64, p []byte) error {
	return h.wc.PReadThrough(offset, p, h.cfio)
}

func (h *FileHandle) Flush() error {
	return h.wc.Flush(h.cfio)
}

func (h *FileHandle) Size() int64 {
	return h.n.Size
}

func (h *FileHandle) Truncate(newsize int64) error {
	if newsize > h.n.Size {
		h.n.Size = newsize
		return nil
	}

	if newsize < h.n.Size {
		h.wc.Truncate(newsize)
		h.cfio.Truncate(newsize)
		h.n.Size = newsize
	}
	return nil
}