package otaru

import (
	"fmt"

	"github.com/nyaxt/otaru/blobstore"
	. "github.com/nyaxt/otaru/util" // FIXME
)

const (
	ChunkSplitSize = 256 * 1024 * 1024 // 256MB
)

const (
	NewChunk      = true
	ExistingChunk = false
)

type ChunkedFileIO struct {
	bs blobstore.RandomAccessBlobStore
	fn *FileNode
	c  Cipher

	newChunkIO func(blobstore.BlobHandle, Cipher) blobstore.BlobHandle
}

func NewChunkedFileIO(bs blobstore.RandomAccessBlobStore, fn *FileNode, c Cipher) *ChunkedFileIO {
	return &ChunkedFileIO{
		bs: bs, fn: fn, c: c,
		newChunkIO: func(bh blobstore.BlobHandle, c Cipher) blobstore.BlobHandle { return NewChunkIO(bh, c) },
	}
}

func (cfio *ChunkedFileIO) OverrideNewChunkIOForTesting(newChunkIO func(blobstore.BlobHandle, Cipher) blobstore.BlobHandle) {
	cfio.newChunkIO = newChunkIO
}

func (cfio *ChunkedFileIO) newFileChunk(newo int64) (FileChunk, error) {
	bpath, err := blobstore.GenerateNewBlobPath(cfio.bs)
	if err != nil {
		return FileChunk{}, err
	}
	fc := FileChunk{Offset: newo, Length: 0, BlobPath: bpath}
	fmt.Printf("new chunk %v\n", fc)
	return fc, nil
}

func (cfio *ChunkedFileIO) PWrite(offset int64, p []byte) error {
	remo := offset
	remp := p
	if len(remp) == 0 {
		return nil
	}

	writeToChunk := func(c *FileChunk, isNewChunk bool, maxChunkLen int64) error {
		if !blobstore.IsReadWriteAllowed(cfio.bs.Flags()) {
			return EPERM
		}

		flags := blobstore.O_RDWR
		if isNewChunk {
			flags |= blobstore.O_CREATE | blobstore.O_EXCL
		}
		bh, err := cfio.bs.Open(c.BlobPath, flags)
		if err != nil {
			return fmt.Errorf("Failed to open path \"%s\" for writing (isNewChunk: %t): %v", c.BlobPath, isNewChunk, err)
		}
		cio := cfio.newChunkIO(bh, cfio.c)

		coff := remo - c.Offset
		n := IntMin(len(remp), int(maxChunkLen-coff))
		if n < 0 {
			return nil
		}
		if err := cio.PWrite(coff, remp[:n]); err != nil {
			return err
		}
		if err := cio.Close(); err != nil {
			return err
		}
		c.Length = int64(cio.Size())

		remo += int64(n)
		remp = remp[n:]
		return nil
	}

	fn := cfio.fn

	for i := 0; i < len(fn.Chunks); i++ {
		c := &fn.Chunks[i]
		if c.Left() > remo {
			// Insert a new chunk @ i

			// try best to align offset at ChunkSplitSize
			newo := remo / ChunkSplitSize * ChunkSplitSize
			maxlen := int64(ChunkSplitSize)
			if i > 0 {
				prev := fn.Chunks[i-1]
				pright := prev.Right()
				if newo < pright {
					maxlen -= pright - newo
					newo = pright
				}
			}
			if i < len(fn.Chunks)-1 {
				next := fn.Chunks[i+1]
				if newo+maxlen > next.Left() {
					maxlen = next.Left() - newo
				}
			}

			newc, err := cfio.newFileChunk(newo)
			if err != nil {
				return err
			}
			fn.Chunks = append(fn.Chunks, FileChunk{})
			copy(fn.Chunks[i+1:], fn.Chunks[i:])
			fn.Chunks[i] = newc

			if err := writeToChunk(&newc, NewChunk, maxlen); err != nil {
				return err
			}
			if len(remp) == 0 {
				return nil
			}

			continue
		}

		// Write to the chunk
		maxlen := int64(ChunkSplitSize)
		if i < len(fn.Chunks)-1 {
			next := fn.Chunks[i+1]
			if c.Left()+maxlen > next.Left() {
				maxlen = next.Left() - c.Left()
			}
		}
		if err := writeToChunk(c, ExistingChunk, maxlen); err != nil {
			return err
		}
		if len(remp) == 0 {
			return nil
		}
	}

	for len(remp) > 0 {
		// Append a new chunk at the end
		newo := remo / ChunkSplitSize * ChunkSplitSize
		maxlen := int64(ChunkSplitSize)

		if len(fn.Chunks) > 0 {
			last := fn.Chunks[len(fn.Chunks)-1]
			lastRight := last.Right()
			if newo < lastRight {
				maxlen -= lastRight - newo
				newo = lastRight
			}
		}

		newc, err := cfio.newFileChunk(newo)
		if err != nil {
			return err
		}
		if err := writeToChunk(&newc, NewChunk, maxlen); err != nil {
			return err
		}

		fn.Chunks = append(fn.Chunks, newc)
	}

	return nil
}

func (cfio *ChunkedFileIO) PRead(offset int64, p []byte) error {
	remo := offset
	remp := p

	if offset < 0 {
		return fmt.Errorf("negative offset %d given", offset)
	}

	cs := cfio.fn.Chunks
	// fmt.Printf("Chunks: %v\n", cs)
	for i := 0; i < len(cs) && len(remp) > 0; i++ {
		c := cs[i]
		if c.Left() > remo+int64(len(remp)) {
			break
		}
		if c.Right() <= remo {
			continue
		}

		coff := remo - c.Left()
		if coff < 0 {
			// Fill gap with zero
			n := Int64Min(int64(len(remp)), -coff)
			for j := int64(0); j < n; j++ {
				remp[j] = 0
			}
			remo += n
			coff = 0
			if len(remp) == 0 {
				return nil
			}
		}

		if !blobstore.IsReadAllowed(cfio.bs.Flags()) {
			return EPERM
		}

		bh, err := cfio.bs.Open(c.BlobPath, blobstore.O_RDONLY)
		if err != nil {
			return fmt.Errorf("Failed to open path \"%s\" for reading: %v", c.BlobPath, err)
		}
		cio := cfio.newChunkIO(bh, cfio.c)

		n := Int64Min(int64(len(p)), c.Length-coff)
		if err := cio.PRead(coff, remp[:n]); err != nil {
			return err
		}
		if err := cio.Close(); err != nil {
			return err
		}

		remo += n
		remp = remp[n:]

		if len(remp) == 0 {
			return nil
		}
	}

	return fmt.Errorf("Attempt to read over file size by %d", len(remp))
}

func (cfio *ChunkedFileIO) Size() int64 {
	cs := cfio.fn.Chunks
	if len(cs) == 0 {
		return 0
	}
	return cs[len(cs)-1].Right()
}

func (cfio *ChunkedFileIO) Close() error {
	return nil
}

func (cfio *ChunkedFileIO) Truncate(size int64) error {
	if !blobstore.IsReadWriteAllowed(cfio.bs.Flags()) {
		return EPERM
	}

	chunks := cfio.fn.Chunks

	for i := len(chunks) - 1; i >= 0; i-- {
		c := &chunks[i]

		if c.Left() >= size {
			// drop the chunk
			continue
		}

		if c.Right() > size {
			// trim the chunk
			chunksize := size - c.Left()

			bh, err := cfio.bs.Open(c.BlobPath, blobstore.O_RDWR)
			if err != nil {
				return err
			}
			cio := cfio.newChunkIO(bh, cfio.c)
			if err := cio.Truncate(chunksize); err != nil {
				return err
			}
			if err := cio.Close(); err != nil {
				return err
			}
			c.Length = int64(cio.Size())
		}

		cfio.fn.Chunks = chunks[:i+1]
		return nil
	}
	cfio.fn.Chunks = []FileChunk{}
	return nil
}
