package otaru

import (
	"compress/zlib"
	"log"
)

const (
	INodeDBSnapshotBlobpath = "INODEDB_SNAPSHOT"
)

func (idb *INodeDB) SaveToBlobStore(bs RandomAccessBlobStore, c Cipher) error {
	raw, err := bs.Open(INodeDBSnapshotBlobpath, O_RDWR|O_CREATE)
	if err != nil {
		return err
	}
	if err := raw.Truncate(0); err != nil {
		return err
	}

	cio := NewChunkIOWithMetadata(raw, c, ChunkPrologue{
		OrigFilename: "*INODEDB_SNAPSHOT*",
		OrigOffset:   0,
	})
	zw := zlib.NewWriter(&OffsetWriter{cio, 0})
	if err := idb.SerializeSnapshot(zw); err != nil {
		return err
	}
	if err := zw.Close(); err != nil {
		return err
	}
	if err := cio.Close(); err != nil {
		return err
	}

	return nil
}

func LoadINodeDBFromBlobStore(bs RandomAccessBlobStore, c Cipher) (*INodeDB, error) {
	raw, err := bs.Open(INodeDBSnapshotBlobpath, O_RDONLY)
	if err != nil {
		return nil, err
	}

	cio := NewChunkIO(raw, c)
	zr, err := zlib.NewReader(&OffsetReader{cio, 0})
	if err != nil {
		return nil, err
	}
	log.Printf("LoadINodeDBFromBlobStore: zlib init success!")
	idb, err := DeserializeINodeDBSnapshot(zr)
	if err != nil {
		return nil, err
	}
	if err := zr.Close(); err != nil {
		return nil, err
	}
	if err := cio.Close(); err != nil {
		return nil, err
	}

	return idb, nil
}