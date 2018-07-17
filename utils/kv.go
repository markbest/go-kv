package utils

import (
	"bufio"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"
	"sync"
)

type KV struct {
	storage  map[string]string
	keys     map[string]string
	locker   *sync.Mutex
	dataPath string
}

// new kv client
func NewKV(dataPath string) *KV {
	return &KV{
		storage:  make(map[string]string),
		keys:     make(map[string]string),
		locker:   &sync.Mutex{},
		dataPath: dataPath,
	}
}

// init kv data
func (k *KV) Init() int {
	count := 0
	_, err := os.Stat(k.dataPath + "/index")
	if !os.IsNotExist(err) {
		// restore data index
		file, err := os.Open(k.dataPath + "/index")
		if err != nil {
			return count
		}
		defer file.Close()

		buff := bufio.NewReader(file)
		for {
			rData, err := buff.ReadString(';')
			if rData != "" {
				indexData := strings.Split(rData, ":")
				key := indexData[0]
				hashKey := indexData[1]
				k.keys[key] = hashKey[:len(hashKey)-1]
				count++
			}
			if io.EOF == err {
				break
			}
		}

		// restore data
		dbFileList, err := GetFileList(k.dataPath)
		if err != nil {
			return count
		}
		if len(dbFileList) > 0 {
			for _, v := range dbFileList {
				paths := strings.Split(v, string(os.PathSeparator))
				hashKey := paths[1] + paths[2]
				file, _ := os.Open(v)
				value, _ := ioutil.ReadAll(file)
				k.storage[hashKey] = string(value)
				file.Close()
			}
		}
	}
	return count
}

// set key & value
func (k *KV) Set(key, value string) {
	k.locker.Lock()
	defer k.locker.Unlock()

	hashKeyName := k.hash(key)
	k.keys[key] = hashKeyName
	k.storage[hashKeyName] = value
}

// generate key name
func (k *KV) hash(key string) string {
	hash := md5.New()
	hash.Write([]byte(key))
	keyName := hex.EncodeToString(hash.Sum(nil))
	return keyName
}

// get key value
func (k *KV) Get(key string) (string, error) {
	k.locker.Lock()
	defer k.locker.Unlock()

	if str, ok := k.keys[key]; ok {
		if value, ok := k.storage[str]; ok {
			return value, nil
		}
	}
	return "", fmt.Errorf("not exist key %s", key)
}

// delete key
func (k *KV) Del(key string) {
	k.locker.Lock()
	defer k.locker.Unlock()

	if str, ok := k.keys[key]; ok {
		delete(k.keys, key)
		delete(k.storage, str)
		os.Remove(k.dataPath + "/" + str[0:2] + "/" + str[2:])
	}
}

// list all keys
func (k *KV) List() []string {
	keys := make([]string, 0)
	if len(k.keys) > 0 {
		for k := range k.keys {
			keys = append(keys, k)
		}
	}
	return keys
}

// clear all kv data
func (k *KV) Clear() {
	k.storage = make(map[string]string)
	k.keys = make(map[string]string)
	os.RemoveAll(k.dataPath)
}

// persistent data
func (k *KV) Persistent() {
	if len(k.storage) > 0 {
		for key, val := range k.storage {
			err := k.storageData(key, val)
			if err != nil {
				panic(err)
			}
		}
	}

	if len(k.keys) > 0 {
		k.storageIndex()
	}
}

// save data to file
func (k *KV) storageData(hashKeyName, value string) error {
	err := os.MkdirAll(k.dataPath+"/"+hashKeyName[0:2], 0666)
	if err != nil {
		return err
	}
	file, err := os.OpenFile(k.dataPath+"/"+hashKeyName[0:2]+"/"+hashKeyName[2:], os.O_CREATE|os.O_RDWR, 0666)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.WriteString(value)
	if err != nil {
		return err
	}
	return nil
}

// save data index
func (k *KV) storageIndex() error {
	file, err := os.OpenFile(k.dataPath+"/index", os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0666)
	defer file.Close()
	if err != nil {
		return err
	}

	for key, val := range k.keys {
		_, err := file.WriteString(key + ":" + val + ";")
		if err != nil {
			continue
		}
	}
	return nil
}
