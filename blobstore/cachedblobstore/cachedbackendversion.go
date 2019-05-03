package cachedblobstore

import (
	"encoding/gob"
	"fmt"
	"sync"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"

	"github.com/nyaxt/otaru/blobstore"
	"github.com/nyaxt/otaru/blobstore/version"
	"github.com/nyaxt/otaru/btncrypt"
	"github.com/nyaxt/otaru/flags"
	"github.com/nyaxt/otaru/logger"
	"github.com/nyaxt/otaru/metadata"
	"github.com/nyaxt/otaru/metadata/statesnapshot"
	oprometheus "github.com/nyaxt/otaru/prometheus"
	"github.com/nyaxt/otaru/util"
)

var (
	saveStateCounter = promauto.NewCounter(prometheus.CounterOpts{
		Name: prometheus.BuildFQName(oprometheus.Namespace, promSubsystem, "cbver_state_save"),
		Help: "Number of times CachedBackedVersion.SaveStateToBlobstore() was called.",
	})
	restoreStateCounter = promauto.NewCounter(prometheus.CounterOpts{
		Name: prometheus.BuildFQName(oprometheus.Namespace, promSubsystem, "cbver_state_restore"),
		Help: "Number of times CachedBackedVersion.RestoreStateFromBlobstore() was called.",
	})

	numCacheEntriesGauge = promauto.NewGauge(prometheus.GaugeOpts{
		Namespace: oprometheus.Namespace,
		Subsystem: promSubsystem,
		Name:      "cbver_num_cache_entries",
		Help:      "Number of cached version entries in CachedBackedVersion.",
	})
)

type CachedBackendVersion struct {
	backendbs    blobstore.BlobStore
	queryVersion version.QueryFunc

	mu    sync.Mutex
	cache map[string]version.Version
}

func NewCachedBackendVersion(backendbs blobstore.BlobStore, queryVersion version.QueryFunc) *CachedBackendVersion {
	return &CachedBackendVersion{
		backendbs:    backendbs,
		queryVersion: queryVersion,

		cache: make(map[string]version.Version),
	}
}

func (cbv *CachedBackendVersion) Set(blobpath string, ver version.Version) {
	cbv.mu.Lock()
	defer cbv.mu.Unlock()

	cbv.cache[blobpath] = ver
}

func (cbv *CachedBackendVersion) Query(blobpath string) (version.Version, error) {
	cbv.mu.Lock()
	defer cbv.mu.Unlock() // FIXME: unlock earlier?

	if ver, ok := cbv.cache[blobpath]; ok {
		logger.Debugf(mylog, "return cached ver for \"%s\" -> %d", blobpath, ver)
		return ver, nil
	}

	r, err := cbv.backendbs.OpenReader(blobpath)
	if err != nil {
		if util.IsNotExist(err) {
			cbv.cache[blobpath] = 0
			return 0, nil
		}
		return -1, fmt.Errorf("Failed to open backend blob for ver query: %v", err)
	}
	defer func() {
		if err := r.Close(); err != nil {
			logger.Criticalf(mylog, "Failed to close backend blob handle for querying version: %v", err)
		}
	}()
	ver, err := cbv.queryVersion(r)
	if err != nil {
		return -1, fmt.Errorf("Failed to query backend blob ver: %v", err)
	}

	cbv.cache[blobpath] = ver
	cbv.updateNumCacheEntriesGauge()
	return ver, nil
}

func (cbv *CachedBackendVersion) Delete(blobpath string) {
	cbv.mu.Lock()
	defer cbv.mu.Unlock()
	delete(cbv.cache, blobpath)
	cbv.updateNumCacheEntriesGauge()
}

func (cbv *CachedBackendVersion) decodeCacheFromGob(dec *gob.Decoder) error {
	cbv.mu.Lock()
	defer cbv.mu.Unlock()

	if err := dec.Decode(&cbv.cache); err != nil {
		return fmt.Errorf("Failed to decode cache map: %v", err)
	}
	return nil
}

func (cbv *CachedBackendVersion) RestoreStateFromBlobstore(c *btncrypt.Cipher, bs blobstore.RandomAccessBlobStore) error {
	restoreStateCounter.Inc()

	bp := metadata.VersionCacheBlobpath
	h, err := bs.Open(bp, flags.O_RDONLY)
	if err != nil {
		return err
	}
	defer h.Close()

	return statesnapshot.Restore(
		&blobstore.OffsetReader{h, 0}, c,
		func(dec *gob.Decoder) error { return cbv.decodeCacheFromGob(dec) },
	)
}

func (cbv *CachedBackendVersion) encodeCacheToGob(enc *gob.Encoder) error {
	cbv.mu.Lock()
	defer cbv.mu.Unlock()

	if err := enc.Encode(cbv.cache); err != nil {
		return fmt.Errorf("Failed to encode cache map: %v", err)
	}
	return nil
}

func (cbv *CachedBackendVersion) SaveStateToBlobstore(c *btncrypt.Cipher, bs blobstore.RandomAccessBlobStore) error {
	saveStateCounter.Inc()

	bp := metadata.VersionCacheBlobpath
	h, err := bs.Open(bp, flags.O_RDWRCREATE)
	if err != nil {
		return err
	}
	defer h.Close()

	return statesnapshot.Save(
		&blobstore.OffsetWriter{h, 0}, c,
		func(enc *gob.Encoder) error { return cbv.encodeCacheToGob(enc) },
	)
}

func (cbv *CachedBackendVersion) updateNumCacheEntriesGauge() {
	numCacheEntriesGauge.Set(float64(len(cbv.cache)))
}
