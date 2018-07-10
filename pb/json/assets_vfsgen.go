// Code generated by vfsgen; DO NOT EDIT.

package json

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	pathpkg "path"
	"time"
)

// Assets statically implements the virtual filesystem provided to vfsgen.
var Assets = func() http.FileSystem {
	fs := vfsgen۰FS{
		"/": &vfsgen۰DirInfo{
			name:    "/",
			modTime: time.Date(2018, 6, 23, 8, 23, 14, 452546057, time.UTC),
		},
		"/otaru.swagger.json": &vfsgen۰CompressedFileInfo{
			name:             "otaru.swagger.json",
			modTime:          time.Date(2018, 7, 10, 14, 16, 30, 126555419, time.UTC),
			uncompressedSize: 18140,

			compressedContent: []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\xec\x5b\xdd\x6e\xdb\x36\x14\xbe\xcf\x53\x10\xda\x2e\x36\x20\xa8\xd2\x62\xd8\x45\x80\x02\x6b\xd3\xa4\x30\x10\x6c\x5d\x32\x74\x17\x43\x21\xd0\xd2\x91\xcc\x96\x22\x55\xf2\x28\xa9\xd7\xf8\xdd\x07\xca\xb2\x2d\x4b\xa2\x65\xfd\x25\xc2\xda\x8b\x02\xa9\xc5\x73\xf4\xf1\xe3\xf9\x37\xfd\xf5\x84\x10\x47\xdf\xd3\x28\x02\xe5\x9c\x13\xe7\xc5\xb3\x33\xe7\xd4\x7c\xc6\x44\x28\x9d\x73\x62\x9e\x13\xe2\x20\x43\x0e\xe6\xf9\x1f\x48\x55\x4a\x5e\xbd\x9b\x65\xab\x08\x71\xee\x40\x69\x26\x85\x79\xf6\x3c\x97\x25\xc4\xf1\xa5\x40\xea\xe3\x56\x01\x21\x8e\xa0\x71\x41\x43\xa2\xe4\x47\xf0\x31\x5f\x4f\x88\x93\x2a\x6e\x9e\x2e\x10\x13\x7d\xee\xba\x11\xc3\x45\x3a\x7f\xe6\xcb\xd8\x15\x4b\xfa\x05\x5d\x69\xc4\x76\xcb\x21\xa6\x2c\x13\x48\x41\xc8\xdf\xb2\x25\x1a\x21\x31\x02\x4e\xb6\x66\x75\x42\xc8\x2a\xdb\x89\xf6\x17\x10\x83\x76\xce\xc9\x3f\x6b\x70\xd9\x3b\xcc\xaa\x0f\xd9\x73\x5f\x0a\x9d\xee\x2d\xa0\x49\xc2\x99\x4f\x91\x49\xe1\x7e\xd4\x52\xec\xd6\x26\x4a\x06\xa9\x7f\xe4\x5a\x8a\x0b\xbd\xa3\xd0\xa5\x09\x73\xef\x9e\xbb\x73\x2e\xe7\x1a\xa5\x02\xd7\x97\x22\x64\x51\x91\xa3\x08\x8a\x94\x11\xe2\xc8\x04\x54\xa6\x7b\x16\x98\xcd\xbe\x05\xbc\x58\x0b\x9d\xee\xd6\x28\xd0\x89\x14\x1a\xf4\x9e\x28\x21\xce\x8b\xb3\xb3\xd2\x47\x84\x38\x01\x68\x5f\xb1\x04\xf3\x33\x2b\x28\xca\x1e\x67\x64\xd1\x8a\x18\x21\xce\x8f\x0a\x42\x23\xf1\x83\x1b\x40\xc8\x04\x33\x1a\xb4\x9b\xcc\xdf\x02\xbe\xde\x6c\x69\x0d\xee\x26\x07\xe4\xec\xa9\x58\x9d\xd4\xfd\xbd\x2a\x6c\x04\x69\xb4\x23\x36\xff\x6c\xab\xfa\x16\xd4\x1d\xf3\x0b\x3a\x3f\x9c\x14\x75\xe5\x7a\x6a\x58\x06\x81\x8a\xed\x91\x73\x0c\xcd\x97\xb9\xd4\xa4\x78\xce\x41\x4d\x8b\x5f\x05\xc6\x23\x3c\x9f\xfa\x0b\x28\x92\x9c\x48\x7d\x98\xe5\x9b\x4c\xf0\x22\x93\x9b\x0e\xcd\x05\x54\x5d\x79\x4e\xa8\xa2\x31\x20\xa8\x32\xdb\x25\xec\x9b\x88\x38\x97\xc1\xb2\x0c\x9c\x09\xdb\x13\x05\x9f\x53\xa6\xc0\x50\x88\x2a\x85\x61\x37\xfc\x39\x05\x8d\xc7\xec\xf7\xc3\x48\x76\x15\x32\x0e\x7a\xa9\x11\x62\x97\x22\x2a\xf7\x2b\x0b\x56\x6d\x5c\xf7\x15\xa2\x9a\x90\x35\x19\x38\x8f\x65\x46\x2c\xa8\x37\x22\x93\x87\xda\x19\x11\x2e\x93\x4c\xa3\x46\xc5\x44\x54\x96\x0d\xa5\x8a\x29\x66\xb9\x97\x09\xfc\xf5\x97\xe2\xbe\x56\xa7\xcd\x38\xeb\xf0\xac\x91\x7e\x4e\x41\x1d\xb2\xf7\x90\x72\xdd\x80\xb5\x93\xb1\x5e\x31\x0e\xb7\x99\xd1\x75\xb7\x56\xf3\x67\x6b\x6b\xbd\x01\x1a\x98\x97\x4f\xc8\x62\x37\x90\xbe\x5b\xed\x3e\x4e\x19\x86\x1a\x70\x34\xbb\x1d\x18\x2d\x07\x11\x8d\xe1\x65\x4c\x20\x98\x4e\xc1\x0a\xb7\x8a\x76\x50\x2f\xdc\xd6\xff\x49\x7a\xd8\xb3\xfe\x56\x0c\x61\x62\xae\xb5\xc5\xf4\xdd\xb7\xf6\x71\x3e\x6d\x01\x54\x38\x96\xc1\xca\x9f\x61\x32\x8a\x08\x7e\x97\x01\xb4\x49\x28\x57\xb9\xcc\x55\xca\xf9\xbb\xfd\x53\x7e\x6a\xeb\x2f\x43\x7b\x2c\x27\xf8\x7f\x96\x1b\xbc\x55\x43\x7b\xcd\x34\xbe\x61\x53\x2a\x8c\x73\x44\x4f\x17\x09\x4b\x3b\x9a\x85\xe4\x81\x83\xf8\x89\x05\x3f\x3f\x90\x97\x2f\xc9\xd9\x29\xc1\x05\x08\x92\xd9\x08\x31\x36\xf4\x6c\x78\x23\xa2\x4a\xd1\x6a\xe0\x43\x88\xcb\x87\x42\x1a\x63\x70\x43\x14\x2e\x71\xfa\xcd\x56\xe9\xa2\x14\x4f\x1b\x87\x14\x17\x0a\x28\x4e\xa9\x88\x58\x03\xfa\x36\x46\x13\x9b\xbd\x3e\x7d\x5a\x66\xc6\x70\x82\xb9\xab\x91\x62\xdb\x51\xe2\xcc\x64\xbd\x37\xaf\x6f\x33\xd1\xe9\x18\x52\x09\xd9\x90\x43\xc5\x8d\xde\x76\x1c\x73\x19\x45\xa0\x5c\x9f\x22\x44\xb2\xc3\xc8\xf6\x62\x27\x38\x29\x96\x77\xb8\x86\xe4\xf8\x3a\x63\xab\x17\xc5\x4b\xf7\xeb\xe6\xaf\x55\xab\xa0\x78\xbb\xdd\xd4\x72\x42\x54\x17\x50\x3d\x56\x78\xf4\xab\x2c\x90\x01\xdb\xac\x96\x29\xfb\x71\x82\xf5\xe1\x51\xc0\xe1\x61\x40\xd7\xb8\xdd\xc7\xd8\x39\x45\xd0\xe8\x71\x19\x79\x20\x50\x2d\x3d\x16\xb4\x0c\x2c\xd7\x99\x86\x6b\x19\x5d\x1a\xf9\x59\x30\x21\x9b\xaf\x41\x37\x95\x20\xc3\x65\xd4\x2a\x82\xff\x69\xca\xca\x6b\x23\x34\x1d\x7a\xb7\x98\x1e\x2b\xa0\xc4\x4c\x78\xb6\xa9\xcd\x23\xcf\xed\x8e\x08\x38\x87\xc3\xdf\x93\x75\x46\xf6\x33\x3a\x66\x72\xca\x62\x36\xc2\x98\x77\xcc\xc1\x69\x27\x2f\xcd\x1b\xa2\xbd\x6b\x1e\xc7\x85\xc3\x75\x01\x3d\x33\x82\xd3\xf1\xd4\x1d\xa8\x21\xe3\xdf\x4e\x6b\x37\x76\x77\x97\x63\xda\x10\xfc\x3e\x97\x9a\x0e\xbb\x39\xa2\xa7\xa1\x76\x7b\x91\xa7\x80\x69\x77\xad\xa6\x34\x48\x32\xff\x35\x31\xa0\xc0\xf8\xc6\x05\xe5\x7c\xff\xc6\x51\xa2\x0c\xf5\xc8\x4a\xa4\x3a\x01\x53\xfb\x35\x02\x69\x98\xbe\x1c\x98\xbc\x14\x59\xc8\xaa\x0f\x9b\xda\x6a\xc0\xb3\x84\x3b\xeb\x19\x65\x4d\xd7\x7b\x06\xf7\xf5\x41\x64\x55\x6b\xaf\xa5\x6f\xa8\x7b\xb0\x56\xbb\xbb\x16\x58\x6d\xf8\xf6\x67\x00\x4f\x7b\xac\x7b\x4f\x6b\x06\x88\xeb\x57\xec\xcd\x0f\xd1\xc0\x27\x0f\x26\xbb\x3c\x10\xaa\x49\x98\x72\xbe\x9e\x26\xd6\xda\x48\x9e\x85\x0e\x41\xac\x95\x4b\xed\x3b\xab\xcb\x3c\xf6\xbc\x53\x54\x1a\x8d\xa1\x34\x01\x15\x7b\xf1\xfe\x18\x6e\x20\xd5\xb1\x0c\x58\xc8\x20\xf0\x90\x35\xb2\xd8\x5e\x7b\xae\xe0\x78\xfb\xfe\xcb\x08\x1c\x6d\xdf\xfd\x3d\x70\x8c\x98\xc5\xb4\x27\xe0\xde\xa6\x77\x2e\x25\x07\x2a\x6c\x8a\x37\x8f\x1b\x39\xb0\x7e\x35\x34\x05\x36\x2c\x98\x0f\x5c\x7b\xec\x81\x7a\x4e\xfd\x4f\x20\x02\x8f\xc5\x09\xf7\xba\x06\x83\x8d\x92\x90\xd3\xa8\x52\x30\x1c\xa3\x20\xbb\xcb\xd7\x0f\xc3\x5a\x45\x3b\x04\x76\xa6\x6b\x46\x68\x3d\x48\xde\xb6\x2d\xa3\xe5\xe2\x75\x45\xbe\x9d\x91\xb5\x4a\xc8\x35\xf7\x3c\x87\x4e\xcb\xc3\x6d\xb4\x0a\x36\x9b\x42\xf4\xdd\xf1\x65\x09\x75\x7b\x3f\xe2\x72\xee\x65\xf3\xb7\x0e\xb6\xab\x91\x62\x37\xc7\x33\x6f\xe5\x20\x06\xcf\x3d\x77\x94\xb3\x60\x14\xcd\x7a\x29\x7c\xcf\x97\xa9\xc0\xc1\x55\x73\xaa\xd1\x4b\x35\x74\x0c\xc4\x4d\x9a\xef\x15\x6b\x3c\xa5\x8e\xaa\x0d\x29\x83\x6b\x16\x69\x3c\x07\xe5\xc9\x70\x8d\x5c\x79\x0b\x2a\x02\x5e\xed\xe8\x06\x7c\xd1\xa0\x6f\xb0\x3b\x6f\xed\xd7\x48\x3d\x9c\x77\xbc\x33\xc8\x34\xe3\x97\x71\xf4\x8e\x51\x7f\x55\x07\x08\x03\x29\x5e\x33\xc1\xfc\x4f\xd0\xd1\xf1\x0f\x29\xdf\x59\xa0\x90\x01\x78\x5c\xfa\x9f\xac\x46\xd8\xb2\xdc\xb7\x5b\xa1\x75\x0e\x3e\x4a\x11\x39\x0c\xec\x5d\x9f\x50\x03\xb2\xc4\xbf\x03\x22\x8d\xf7\xa6\x28\xce\xd5\xec\xfa\xb2\x38\x25\x7a\x33\xbb\xd9\xbc\x73\x3b\x33\x74\x02\x08\x69\xca\x33\x54\xd9\xfa\x3a\x08\x59\x2b\x3e\x85\x62\x7b\x80\x56\xb8\x6f\x9f\x56\x4c\x8e\xec\xdf\xe1\x33\xcc\xf7\x56\x7d\x9c\x56\xdd\xe2\x62\xe5\x8b\x66\x7d\xf2\x52\x65\xb6\x48\x86\xac\xa6\x2d\x93\xcc\x82\x68\x73\x29\x5d\x6a\x3c\x46\x6c\x92\x0e\xb9\x20\x87\x3b\xe0\xe3\x46\xce\xea\x37\x73\xd3\x6d\x92\x2a\x58\x3b\xf4\x48\x16\x1d\x53\xc8\x6d\xc5\x83\x97\x56\xef\x38\x3c\x2c\x98\x86\xad\x15\xb3\xc8\x18\x23\x44\x2e\xd7\x3f\x14\xee\x3b\x11\xa9\xfc\x46\xa7\x4f\xa7\x2c\x83\x26\xde\x6d\xe3\xbd\x25\x1e\x31\xdf\xac\xf9\x79\x61\x9f\x21\xbe\x5a\x7a\x2a\xb5\x12\xd8\x7a\x20\x59\x38\x9c\x00\x34\x53\x10\x78\x47\x24\xfd\x76\x9b\xee\x7f\x46\x3a\xf5\x7d\xd0\xd6\x12\xba\xcf\xae\x41\x29\xa9\xbc\x18\xb4\xa6\x51\xef\x6d\xd7\xdd\xc0\xb2\x6f\xbb\x2c\x5c\xfd\x06\xb7\x07\x65\x91\xf4\x8e\xeb\x9d\x6a\x59\x91\x9d\xc6\xa6\x54\xf9\x9d\xe6\x4d\x22\x8d\xbd\x48\x2a\x99\x22\x13\x23\xd4\x5a\x0b\xa9\xb1\x6b\x31\x9d\x8c\x51\xd9\x1f\xa8\x80\xbb\x2b\x8d\x21\xf6\x28\xe7\xb2\xe3\xd0\xa0\x49\xb5\x5e\x76\x1c\xa3\x34\xf4\xc9\x5e\x64\x05\xdc\xfd\xc4\x8d\xde\x30\x18\xb9\xe5\x2e\xdf\x09\xe8\xe3\xac\x0c\x3d\x5f\xc6\x31\x6b\x9a\x45\xd4\xee\x76\x9e\x32\x1e\x78\x8b\xf2\xa5\xd5\x56\xd2\x47\xa4\xfb\x46\x42\x2a\xbf\x23\x9b\x58\x43\x9d\xff\x82\x75\x70\xbd\x63\x57\x10\xd5\x9f\x4d\x36\xe5\x94\x13\xf3\x6f\x75\xf2\x5f\x00\x00\x00\xff\xff\x88\x1a\x27\xca\xdc\x46\x00\x00"),
		},
	}
	fs["/"].(*vfsgen۰DirInfo).entries = []os.FileInfo{
		fs["/otaru.swagger.json"].(os.FileInfo),
	}

	return fs
}()

type vfsgen۰FS map[string]interface{}

func (fs vfsgen۰FS) Open(path string) (http.File, error) {
	path = pathpkg.Clean("/" + path)
	f, ok := fs[path]
	if !ok {
		return nil, &os.PathError{Op: "open", Path: path, Err: os.ErrNotExist}
	}

	switch f := f.(type) {
	case *vfsgen۰CompressedFileInfo:
		gr, err := gzip.NewReader(bytes.NewReader(f.compressedContent))
		if err != nil {
			// This should never happen because we generate the gzip bytes such that they are always valid.
			panic("unexpected error reading own gzip compressed bytes: " + err.Error())
		}
		return &vfsgen۰CompressedFile{
			vfsgen۰CompressedFileInfo: f,
			gr: gr,
		}, nil
	case *vfsgen۰DirInfo:
		return &vfsgen۰Dir{
			vfsgen۰DirInfo: f,
		}, nil
	default:
		// This should never happen because we generate only the above types.
		panic(fmt.Sprintf("unexpected type %T", f))
	}
}

// vfsgen۰CompressedFileInfo is a static definition of a gzip compressed file.
type vfsgen۰CompressedFileInfo struct {
	name              string
	modTime           time.Time
	compressedContent []byte
	uncompressedSize  int64
}

func (f *vfsgen۰CompressedFileInfo) Readdir(count int) ([]os.FileInfo, error) {
	return nil, fmt.Errorf("cannot Readdir from file %s", f.name)
}
func (f *vfsgen۰CompressedFileInfo) Stat() (os.FileInfo, error) { return f, nil }

func (f *vfsgen۰CompressedFileInfo) GzipBytes() []byte {
	return f.compressedContent
}

func (f *vfsgen۰CompressedFileInfo) Name() string       { return f.name }
func (f *vfsgen۰CompressedFileInfo) Size() int64        { return f.uncompressedSize }
func (f *vfsgen۰CompressedFileInfo) Mode() os.FileMode  { return 0444 }
func (f *vfsgen۰CompressedFileInfo) ModTime() time.Time { return f.modTime }
func (f *vfsgen۰CompressedFileInfo) IsDir() bool        { return false }
func (f *vfsgen۰CompressedFileInfo) Sys() interface{}   { return nil }

// vfsgen۰CompressedFile is an opened compressedFile instance.
type vfsgen۰CompressedFile struct {
	*vfsgen۰CompressedFileInfo
	gr      *gzip.Reader
	grPos   int64 // Actual gr uncompressed position.
	seekPos int64 // Seek uncompressed position.
}

func (f *vfsgen۰CompressedFile) Read(p []byte) (n int, err error) {
	if f.grPos > f.seekPos {
		// Rewind to beginning.
		err = f.gr.Reset(bytes.NewReader(f.compressedContent))
		if err != nil {
			return 0, err
		}
		f.grPos = 0
	}
	if f.grPos < f.seekPos {
		// Fast-forward.
		_, err = io.CopyN(ioutil.Discard, f.gr, f.seekPos-f.grPos)
		if err != nil {
			return 0, err
		}
		f.grPos = f.seekPos
	}
	n, err = f.gr.Read(p)
	f.grPos += int64(n)
	f.seekPos = f.grPos
	return n, err
}
func (f *vfsgen۰CompressedFile) Seek(offset int64, whence int) (int64, error) {
	switch whence {
	case io.SeekStart:
		f.seekPos = 0 + offset
	case io.SeekCurrent:
		f.seekPos += offset
	case io.SeekEnd:
		f.seekPos = f.uncompressedSize + offset
	default:
		panic(fmt.Errorf("invalid whence value: %v", whence))
	}
	return f.seekPos, nil
}
func (f *vfsgen۰CompressedFile) Close() error {
	return f.gr.Close()
}

// vfsgen۰DirInfo is a static definition of a directory.
type vfsgen۰DirInfo struct {
	name    string
	modTime time.Time
	entries []os.FileInfo
}

func (d *vfsgen۰DirInfo) Read([]byte) (int, error) {
	return 0, fmt.Errorf("cannot Read from directory %s", d.name)
}
func (d *vfsgen۰DirInfo) Close() error               { return nil }
func (d *vfsgen۰DirInfo) Stat() (os.FileInfo, error) { return d, nil }

func (d *vfsgen۰DirInfo) Name() string       { return d.name }
func (d *vfsgen۰DirInfo) Size() int64        { return 0 }
func (d *vfsgen۰DirInfo) Mode() os.FileMode  { return 0755 | os.ModeDir }
func (d *vfsgen۰DirInfo) ModTime() time.Time { return d.modTime }
func (d *vfsgen۰DirInfo) IsDir() bool        { return true }
func (d *vfsgen۰DirInfo) Sys() interface{}   { return nil }

// vfsgen۰Dir is an opened dir instance.
type vfsgen۰Dir struct {
	*vfsgen۰DirInfo
	pos int // Position within entries for Seek and Readdir.
}

func (d *vfsgen۰Dir) Seek(offset int64, whence int) (int64, error) {
	if offset == 0 && whence == io.SeekStart {
		d.pos = 0
		return 0, nil
	}
	return 0, fmt.Errorf("unsupported Seek in directory %s", d.name)
}

func (d *vfsgen۰Dir) Readdir(count int) ([]os.FileInfo, error) {
	if d.pos >= len(d.entries) && count > 0 {
		return nil, io.EOF
	}
	if count <= 0 || count > len(d.entries)-d.pos {
		count = len(d.entries) - d.pos
	}
	e := d.entries[d.pos : d.pos+count]
	d.pos += count
	return e, nil
}
