package tsoutil

import (
	"path"
	"time"

	etcdkv "github.com/zilliztech/milvus-distributed/internal/kv/etcd"
	"go.etcd.io/etcd/clientv3"
)

const (
	logicalBits     = 18
	logicalBitsMask = (1 << logicalBits) - 1
)

func ComposeTS(physical, logical int64) uint64 {
	return uint64((physical << logicalBits) + logical)
}

// ParseTS parses the ts to (physical,logical).
func ParseTS(ts uint64) (time.Time, uint64) {
	logical := ts & logicalBitsMask
	physical := ts >> logicalBits
	physicalTime := time.Unix(int64(physical/1000), int64(physical)%1000*time.Millisecond.Nanoseconds())
	return physicalTime, logical
}

func NewTSOKVBase(etcdAddr []string, tsoRoot, subPath string) *etcdkv.EtcdKV {
	client, _ := clientv3.New(clientv3.Config{
		Endpoints:   etcdAddr,
		DialTimeout: 5 * time.Second,
	})
	return etcdkv.NewEtcdKV(client, path.Join(tsoRoot, subPath))
}
