package etcd

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/client/v3/concurrency"
)

var (
	dialTimeout    = 5 * time.Second
	requestTimeout = 1 * time.Second
)

type Etcd struct {
	Client      *clientv3.Client
	KV          clientv3.KV
	Concurrency *concurrency.Session
}

type Lock struct {
	mutex   *concurrency.Mutex
	ctx     context.Context
	session *concurrency.Session
	cancel  context.CancelFunc
}

func (etcd *Etcd) Connect(endpointlist string, user string, password string, use_tls bool) {
	var err error
	endpoints := strings.Split(endpointlist, ",")
	cfg := clientv3.Config{
		Endpoints:   endpoints,
		DialTimeout: dialTimeout,
		Username:    user,
		Password:    password,
	}
	if use_tls {
		cfg.TLS = &tls.Config{InsecureSkipVerify: true}
	}
	etcd.Client, err = clientv3.New(cfg)
	if err != nil {
		log.Fatal(err)
	}
	etcd.KV = clientv3.NewKV(etcd.Client)
	//etcd.Concurrency, _ = concurrency.NewSession(etcd.Client)
}

func (etcd *Etcd) Watch(ctx context.Context, key_path string) clientv3.WatchChan {
	rch := etcd.Client.Watch(ctx, key_path, clientv3.WithPrefix())
	return rch
}

func (etcd *Etcd) WatchKey(ctx context.Context, key_path string) clientv3.WatchChan {
	rch := etcd.Client.Watch(ctx, key_path)
	return rch
}

func (etcd *Etcd) GetLock(lock_path string) (Lock, error) {
	var lock Lock
	var err error
	lock.session, err = concurrency.NewSession(etcd.Client, concurrency.WithTTL(1))
	if err != nil {
		return lock, err
	}

	m := concurrency.NewMutex(lock.session, lock_path)
	d := time.Now().Add(2 * time.Second)
	//ctx := context.Background()
	lock.ctx, lock.cancel = context.WithDeadline(context.Background(), d)
	// acquire lock (or wait to have it)
	if err := m.Lock(lock.ctx); err != nil {
		return lock, err
	}
	lock.mutex = m
	return lock, nil
}

func (etcd *Etcd) ReleaseLock(lock Lock) error {
	err := lock.mutex.Unlock(lock.ctx)
	return err
}

func (etcd *Etcd) DeferLock(lock Lock) {
	lock.session.Close()
	lock.cancel()
}

func (etcd *Etcd) PutWithLease(key_path string, value string, ttl int64) (clientv3.LeaseID, error) {
	ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
	resp, err := etcd.Client.Grant(ctx, ttl)
	if err != nil {
		cancel()
		return resp.ID, err
	}
	_, err = etcd.KV.Put(ctx, key_path, value, clientv3.WithLease(resp.ID))
	cancel()
	return resp.ID, err
}

func (etcd *Etcd) KeepLeaseAliveOnce(lease_id clientv3.LeaseID) error {
	ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
	_, err := etcd.Client.KeepAliveOnce(ctx, lease_id)
	cancel()
	return err
}

func (etcd *Etcd) KeepLeaseAlive(lease_id clientv3.LeaseID) error {
	ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
	_, err := etcd.Client.KeepAlive(ctx, lease_id)
	cancel()
	return err
}

func (etcd *Etcd) RevokeLease(lease_id clientv3.LeaseID) error {
	ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
	_, err := etcd.Client.Revoke(ctx, lease_id)
	cancel()
	return err
}

func (etcd *Etcd) GetValue(key_path string) ([]byte, error) {
	ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
	gr, err := etcd.KV.Get(ctx, key_path)
	cancel()
	if err != nil {
		return nil, err
	}
	if len(gr.Kvs) == 0 {
		return nil, errors.New("key does not exists")
	}
	return gr.Kvs[0].Value, err
}

func (etcd *Etcd) GetObject(key_path string) (map[string][]byte, error) {
	var err error
	ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
	resp, err := etcd.KV.Get(ctx, key_path, clientv3.WithPrefix())
	cancel()
	var result = make(map[string][]byte)
	if err != nil {
		return result, err
	}
	var exists bool = false
	for _, ev := range resp.Kvs {
		exists = true
		if len(key_path) > 0 {
			var path string
			if string(key_path[len(key_path)-1]) == "/" {
				path = key_path
			} else {
				path = key_path + "/"
			}
			shorted := strings.TrimPrefix(string(ev.Key), path)
			// shorted := strings.TrimPrefix(string(ev.Key), key_path)
			result[shorted] = ev.Value
		} else {
			result[string(ev.Key)] = ev.Value
		}
	}
	if !exists {
		err = fmt.Errorf("path does not exists")
	}
	return result, err
}

func (etcd *Etcd) GetObjectList(key_path string) ([]map[string][]byte, error) {
	var err error
	ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
	resp, err := etcd.KV.Get(ctx, key_path, clientv3.WithPrefix())
	cancel()
	var result = make([]map[string][]byte, 0)
	var tmp = make(map[string]map[string][]byte)
	if err != nil {
		return result, err
	}
	var exists bool = false
	for _, ev := range resp.Kvs {
		exists = true
		keys := strings.Split(string(ev.Key), "/")
		k1 := keys[len(keys)-1]
		k2 := keys[len(keys)-2]
		if tmp[k2] == nil {
			tmp[k2] = make(map[string][]byte)
		}
		tmp[k2][k1] = ev.Value
	}
	if !exists {
		err = fmt.Errorf("path does not exists")
	}
	for _, v := range tmp {
		result = append(result, v)
	}
	return result, err
}

func (etcd *Etcd) Put(key_path string, value string) error {
	ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
	_, err := etcd.KV.Put(ctx, key_path, value)
	cancel()
	return err
}

func (etcd *Etcd) Delete(key_path string) error {
	ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
	_, err := etcd.KV.Delete(ctx, key_path, clientv3.WithPrefix())
	cancel()
	return err
}
