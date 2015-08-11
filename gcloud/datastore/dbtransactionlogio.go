package datastore

import (
	"fmt"
	"sort"
	"sync"
	"time"

	"golang.org/x/net/context"
	"google.golang.org/cloud/datastore"

	"github.com/nyaxt/otaru/btncrypt"
	gcutil "github.com/nyaxt/otaru/gcloud/util"
	"github.com/nyaxt/otaru/inodedb"
	"github.com/nyaxt/otaru/logger"
)

var txlog = logger.Registry().Category("dbtxlogio")

type DBTransactionLogIO struct {
	cfg     *Config
	rootKey *datastore.Key

	mu        sync.Mutex
	nextbatch []inodedb.DBTransaction

	muSync     sync.Mutex
	committing []inodedb.DBTransaction
}

const kindTransaction = "OtaruINodeDBTx"

var _ = inodedb.DBTransactionLogIO(&DBTransactionLogIO{})

func NewDBTransactionLogIO(cfg *Config) *DBTransactionLogIO {
	return &DBTransactionLogIO{
		cfg:       cfg,
		rootKey:   datastore.NewKey(ctxNoNamespace, kindTransaction, cfg.rootKeyStr, 0, nil),
		nextbatch: make([]inodedb.DBTransaction, 0),
	}
}

type storedbtx struct {
	OpsJSON []byte `datastore:,noindex`
}

func (txio *DBTransactionLogIO) encodeKey(id inodedb.TxID) *datastore.Key {
	return datastore.NewKey(ctxNoNamespace, kindTransaction, "", int64(id), txio.rootKey)
}

func (txio *DBTransactionLogIO) encode(tx inodedb.DBTransaction) (*datastore.Key, *storedbtx, error) {
	key := txio.encodeKey(tx.TxID)

	jsonops, err := inodedb.EncodeDBOperationsToJson(tx.Ops)
	if err != nil {
		return nil, nil, fmt.Errorf("Failed to encode dbtx: %v", err)
	}

	c := txio.cfg.c
	env, err := btncrypt.Encrypt(c, jsonops)
	if err != nil {
		return nil, nil, fmt.Errorf("Failed to encrypt OpsJSON: %v", err)
	}

	return key, &storedbtx{OpsJSON: env}, nil
}

func decode(c btncrypt.Cipher, key *datastore.Key, stx *storedbtx) (inodedb.DBTransaction, error) {
	jsonop, err := btncrypt.Decrypt(c, stx.OpsJSON, len(stx.OpsJSON)-c.FrameOverhead())
	if err != nil {
		return inodedb.DBTransaction{}, fmt.Errorf("Failed to decrypt OpsJSON: %v", err)
	}

	ops, err := inodedb.DecodeDBOperationsFromJson(jsonop)
	if err != nil {
		return inodedb.DBTransaction{}, err
	}

	return inodedb.DBTransaction{TxID: inodedb.TxID(key.ID()), Ops: ops}, nil
}

func (txio *DBTransactionLogIO) AppendTransaction(tx inodedb.DBTransaction) error {
	txio.mu.Lock()
	defer txio.mu.Unlock()

	txio.nextbatch = append(txio.nextbatch, tx)
	return nil
}

func (txio *DBTransactionLogIO) Sync() error {
	start := time.Now()

	txio.muSync.Lock()
	defer txio.muSync.Unlock()

	txio.mu.Lock()
	if len(txio.committing) != 0 {
		panic("I should be the only one committing.")
	}
	txio.committing = txio.nextbatch
	batch := txio.committing
	txio.nextbatch = make([]inodedb.DBTransaction, 0)
	txio.mu.Unlock()
	rollback := func() {
		txio.mu.Lock()
		txio.nextbatch = append(txio.committing, txio.nextbatch...)
		txio.committing = []inodedb.DBTransaction{}
		txio.mu.Unlock()
	}

	if len(batch) == 0 {
		return nil
	}

	cli, err := txio.cfg.getClient(context.Background())
	if err != nil {
		rollback()
		return err
	}

	keys := make([]*datastore.Key, 0, len(batch))
	stxs := make([]*storedbtx, 0, len(batch))
	for _, tx := range batch {
		key, stx, err := txio.encode(tx)
		if err != nil {
			rollback()
			return err
		}

		keys = append(keys, key)
		stxs = append(stxs, stx)
	}

	dstx, err := cli.NewTransaction(context.Background(), datastore.Serializable)
	if err != nil {
		rollback()
		return err
	}

	if _, err := dstx.PutMulti(keys, stxs); err != nil {
		rollback()
		dstx.Rollback()
		return err
	}

	if _, err := dstx.Commit(); err != nil {
		rollback()
		return err
	}

	txio.mu.Lock()
	txio.committing = []inodedb.DBTransaction{}
	txio.mu.Unlock()

	logger.Infof(txlog, "Sync() took %s. Committed %d txs", time.Since(start), len(stxs))
	return nil
}

type txsorter []inodedb.DBTransaction

func (s txsorter) Len() int           { return len(s) }
func (s txsorter) Less(i, j int) bool { return s[i].TxID < s[j].TxID }
func (s txsorter) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (txio *DBTransactionLogIO) QueryTransactions(minID inodedb.TxID) (txs []inodedb.DBTransaction, err error) {
	err = gcutil.RetryIfNeeded(func() error {
		txs, err = txio.queryTransactionsOnce(minID)
		return err
	})
	return
}

func (txio *DBTransactionLogIO) queryTransactionsOnce(minID inodedb.TxID) ([]inodedb.DBTransaction, error) {
	start := time.Now()
	txs := []inodedb.DBTransaction{}

	txio.mu.Lock()
	for _, tx := range txio.committing {
		if tx.TxID >= minID {
			txs = append(txs, tx)
		}
	}
	for _, tx := range txio.nextbatch {
		if tx.TxID >= minID {
			txs = append(txs, tx)
		}
	}
	txio.mu.Unlock()

	cli, err := txio.cfg.getClient(context.Background())
	if err != nil {
		return nil, err
	}

	dstx, err := cli.NewTransaction(context.Background(), datastore.Serializable)
	if err != nil {
		return nil, err
	}

	q := datastore.NewQuery(kindTransaction).Ancestor(txio.rootKey).Transaction(dstx)
	if minID != inodedb.AnyVersion {
		minKey := txio.encodeKey(minID)
		q = q.Filter("__key__ >=", minKey)
	}

	it := cli.Run(context.Background(), q)
	for {
		var stx storedbtx
		key, err := it.Next(&stx)
		if err != nil {
			if err == datastore.Done {
				break
			}
			dstx.Commit()
			return []inodedb.DBTransaction{}, err
		}

		tx, err := decode(txio.cfg.c, key, &stx)
		if err != nil {
			dstx.Commit()
			return []inodedb.DBTransaction{}, err
		}

		txs = append(txs, tx)
	}

	// FIXME: not sure if Rollback() is better
	if _, err := dstx.Commit(); err != nil {
		return nil, err
	}

	sort.Sort(txsorter(txs))
	uniqed := make([]inodedb.DBTransaction, 0, len(txs))
	var prevId inodedb.TxID
	for _, tx := range txs {
		if tx.TxID == prevId {
			continue
		}

		uniqed = append(uniqed, tx)
		prevId = tx.TxID
	}

	logger.Infof(txlog, "QueryTransactions(%v) took %s", minID, time.Since(start))
	return uniqed, nil
}

const maxWriteEntriesPerTx = 500 // Google Cloud Datastore limit on number of write entries per tx

func (txio *DBTransactionLogIO) DeleteTransactions(smallerThanID inodedb.TxID) error {
	start := time.Now()

	txio.mu.Lock()
	batch := make([]inodedb.DBTransaction, 0, len(txio.nextbatch))
	for _, tx := range txio.nextbatch {
		if tx.TxID < smallerThanID {
			continue
		}
		batch = append(batch, tx)
	}
	txio.nextbatch = batch
	txio.mu.Unlock()

	cli, err := txio.cfg.getClient(context.Background())
	if err != nil {
		return err
	}

	ndel := 0
	for {
		needAnotherTx := false
		txStart := time.Now()
		dstx, err := cli.NewTransaction(context.Background(), datastore.Serializable)
		if err != nil {
			return err
		}

		keys := []*datastore.Key{}
		ltkey := txio.encodeKey(smallerThanID)
		q := datastore.NewQuery(kindTransaction).Ancestor(txio.rootKey).Filter("__key__ <", ltkey).KeysOnly().Transaction(dstx)
		it := cli.Run(context.Background(), q)
		for {
			k, err := it.Next(nil)
			if err != nil {
				if err == datastore.Done {
					break
				}
				dstx.Rollback()
				return err
			}

			keys = append(keys, k)
			if len(keys) == maxWriteEntriesPerTx {
				needAnotherTx = true
				break
			}
		}

		if err := dstx.DeleteMulti(keys); err != nil {
			dstx.Rollback()
			return err
		}

		if _, err := dstx.Commit(); err != nil {
			return err
		}
		ndel += len(keys)

		if needAnotherTx {
			logger.Infof(txlog, "DeleteTransactions(%v): A tx deleting %d entries took %s. Starting next tx to delete more.", smallerThanID, len(keys), time.Since(txStart))
		} else {
			break
		}
	}
	logger.Infof(txlog, "DeleteTransactions(%v) deleted %d entries. tx took %s", smallerThanID, ndel, time.Since(start))
	return nil
}

func (txio *DBTransactionLogIO) DeleteAllTransactions() error {
	return txio.DeleteTransactions(inodedb.LatestVersion)
}
