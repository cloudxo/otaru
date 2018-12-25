package otaruapiserver

import (
	"fmt"
	"io"
	"mime"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/nyaxt/otaru/apiserver"
	"github.com/nyaxt/otaru/apiserver/jwt"
	"github.com/nyaxt/otaru/blobstore"
	"github.com/nyaxt/otaru/filesystem"
	"github.com/nyaxt/otaru/flags"
	"github.com/nyaxt/otaru/inodedb"
	"github.com/nyaxt/otaru/logger"
)

type fileHandler struct {
	fs      *filesystem.FileSystem
	jwtauth *jwt.JWTAuthProvider
}

type content struct {
	h      *filesystem.FileHandle
	offset int64
	size   int64
}

var _ = io.ReadSeeker(&content{})

func (c *content) Read(p []byte) (int, error) {
	// FIXME: not sure if this handles eof correctly
	n, err := c.h.ReadAt(p, c.offset)
	c.offset += int64(n)
	return n, err
}

func (c *content) Seek(offset int64, whence int) (int64, error) {
	switch whence {
	case io.SeekStart:
		if offset < 0 || c.size < offset {
			return c.offset, os.ErrInvalid
		}
		c.offset = offset
		break
	case io.SeekCurrent:
		logger.Panicf(mylog, "No implemention for Seek(whence=io.SeekCurrent), as we don't expect net/http to use this")
		break
	case io.SeekEnd:
		c.offset = c.size
		break
	default:
		return c.offset, os.ErrInvalid
	}

	return c.offset, nil
}

func (fh *fileHandler) serveGet(w http.ResponseWriter, r *http.Request, id inodedb.ID, filename string) {
	auth := r.Header.Get("Authorization")
	ui, err := fh.jwtauth.UserInfoFromAuthHeader(auth)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	if ui.Role < jwt.RoleReadOnly {
		http.Error(w, "", http.StatusForbidden)
		return
	}

	h, err := fh.fs.OpenFile(id, flags.O_RDONLY)
	if err != nil {
		logger.Debugf(mylog, "serveGet(id: %v). OpenFile failed: %v", id, err)
		http.Error(w, "Failed to open file", http.StatusInternalServerError)
		return
	}
	defer h.Close()

	a, err := h.Attr()
	if err != nil {
		logger.Debugf(mylog, "serveGet(id: %v). Attr failed: %v", id, err)
		http.Error(w, "Failed to attr", http.StatusInternalServerError)
		return
	}

	if filename == "" {
		filename = filepath.Base(a.OrigPath)
		if filename == "" {
			filename = fmt.Sprintf("%d.bin", id)
		}
	}

	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.Header().Set("Content-Disposition", "attachment; filename*=UTF-8''"+url.QueryEscape(filename))
	ctype := mime.TypeByExtension(filepath.Ext(filename))
	if ctype == "" {
		ctype = "application/octet-stream"
	}
	w.Header().Set("Content-Type", ctype)

	c := &content{h: h, offset: 0, size: a.Size}
	http.ServeContent(w, r, filename, a.ModifiedT, c)
}

func (fh *fileHandler) servePut(w http.ResponseWriter, r *http.Request, id inodedb.ID, filename string) {
	auth := r.Header.Get("Authorization")
	ui, err := fh.jwtauth.UserInfoFromAuthHeader(auth)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	if ui.Role < jwt.RoleAdmin {
		http.Error(w, "", http.StatusForbidden)
		return
	}

	h, err := fh.fs.OpenFile(id, flags.O_WRONLY)
	if err != nil {
		logger.Debugf(mylog, "servePut(id: %v). OpenFile failed: %v", id, err)
		http.Error(w, "Failed to open file", http.StatusInternalServerError)
		return
	}
	defer h.Close()

	// FIXME: parse offset
	offset := int64(0)
	nw, err := io.Copy(&blobstore.OffsetWriter{h, offset}, r.Body)
	if err != nil {
		logger.Debugf(mylog, "servePut(id: %v). io.Copy failed: %v", id, err)
		http.Error(w, "Failed to write file", http.StatusInternalServerError)
		return
	}
	logger.Debugf(mylog, "servePut(id: %v). written %d", id, nw)

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Write([]byte("\"ok\""))
}

func (fh *fileHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// path: /inodeid/filename
	args := strings.Split(strings.TrimPrefix(r.URL.Path, "/"), "/")

	if len(args) == 0 || len(args) > 2 {
		http.Error(w, "Error parsing url", http.StatusNotFound)
		return
	}

	nid, err := strconv.ParseUint(args[0], 10, 64)
	if err != nil {
		http.Error(w, "Error parsing inodeid", http.StatusBadRequest)
		return
	}
	id := inodedb.ID(nid)

	filename := ""
	if len(args) == 2 {
		filename = args[1]
	}

	if r.Method == "GET" || r.Method == "HEAD" {
		fh.serveGet(w, r, id, filename)
	} else if r.Method == "PUT" {
		fh.servePut(w, r, id, filename)
	} else {
		http.Error(w, "Unsupported method", http.StatusBadRequest)
		return
	}
}

func InstallFileHandler(fs *filesystem.FileSystem, jwtauth *jwt.JWTAuthProvider) apiserver.Option {
	return apiserver.AddMuxHook(func(mux *http.ServeMux) {
		filePrefix := "/file/"
		mux.Handle(filePrefix, http.StripPrefix(filePrefix, &fileHandler{fs, jwtauth}))
	})
}
