package mybolt

import (
	"github.com/kwuyoucloud/spider/pkg/log"
	"os"
	"time"

	"github.com/boltdb/bolt"
)

const (
	// ConnectionTimeout is the timeout second for database connection.
	ConnectionTimeout = 1 * time.Second
)

// AddKeyValueonBolt is the function that add key string, value string into bucket_name of file [db_file_name]
// The db_file_name shoule be one file that at the root path of this application.
// The bucket_name is the bucket in the db file.
// If add sucess then return nil.
// If fail then return the error.
func AddKeyValueonBolt(dbFileName string, bucketName string, key string, value string) error {
	// Open db file with write permission.
	db, err := bolt.Open(dbFileName, 0666, nil)
	defer db.Close()

	if err != nil {
		log.FatalLog(err)
		return err
	}

	// Insert key, value into bucket.
	err = db.Update(func(tx *bolt.Tx) error {
		// If the bucket is not exist, create it.
		// If create bucket with error, return the error.
		b, createErr := tx.CreateBucketIfNotExists([]byte(bucketName))
		if createErr != nil {
			log.FatalLog(createErr)
			return createErr
		}
		// Insert key, value into bucket.
		errB := b.Put([]byte(key), []byte(value))
		if errB != nil {
			return errB
		}
		return nil
	})
	return err
}

// AddKeyValuesonBolt comments
func AddKeyValuesonBolt(dbFileName string, bucketName string, keyValues map[string]string) error {
	// Open db file with write permission.
	db, err := bolt.Open(dbFileName, 0666, nil)
	defer db.Close()

	if err != nil {
		log.FatalLog(err)
		return err
	}

	// Insert key, value into bucket.
	err = db.Update(func(tx *bolt.Tx) error {
		// If the bucket is not exist, create it.
		// If create bucket with error, return the error.
		b, createErr := tx.CreateBucketIfNotExists([]byte(bucketName))
		if createErr != nil {
			log.FatalLog(createErr)
			return createErr
		}
		// Insert key, value into bucket.
		for k, v := range keyValues {
			// If the key has exists, continue;
			// if the key is not exist, put into the bucket.
			// bucket put achieve this naturally.
			errB := b.Put([]byte(k), []byte(v))
			if errB != nil {
				return errB
			}
		}
		return nil
	})
	return err
}

// GetValueFromBolt can get value of the key in the bucket of db file.
func GetValueFromBolt(dbFileName string, bucketName string, key string) string {
	// Open the database file.
	db, err := bolt.Open(dbFileName, 0444, nil)
	defer db.Close()

	if err != nil {
		log.ErrLog(err)
		return ""
	}

	// initlization value as type []byte
	var value []byte

	if err = db.View(func(tx *bolt.Tx) error {
		// Open bucket use bucketName
		b := tx.Bucket([]byte(bucketName))
		// Get value from bucket.
		value = b.Get([]byte(key))
		return nil
	}); err != nil {
		log.ErrLog(err)
	}
	return string(value)
}

// GetValuesFromBolt can get value of the key in the bucket of db file.
func GetValuesFromBolt(dbFileName string, bucketName string, keys []string) map[string]string {
	var resultKeyValues map[string]string
	db, err := bolt.Open(dbFileName, 0444, nil)
	defer db.Close()

	if err != nil {
		log.ErrLog(err)
		// return nil as result
		return nil
	}

	if err = db.View(func(tx *bolt.Tx) error {
		var length = len(keys)

		b := tx.Bucket([]byte(bucketName))
		// Go through the loop of keys.
		for i := 0; i < length; i++ {
			key := keys[i]
			valuestr := b.Get([]byte(key))
			resultKeyValues[key] = string(valuestr)
		}
		return nil
	}); err != nil {
		log.ErrLog(err)
	}
	return resultKeyValues
}

// GetAllKeyValuesFromBolt function can get all key and values from bucketName in the dbFileName file.
func GetAllKeyValuesFromBolt(dbFileName string, bucketName string) map[string]string {
	keyvalueMap := make(map[string]string, 25)
	db, err := bolt.Open(dbFileName, 0444, nil)
	defer db.Close()

	if err != nil {
		log.ErrLog(err)
	}

	if err = db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucketName))
		// Go through the loop of keys.
		b.ForEach(func(k, v []byte) error {
			keyvalueMap[string(k)] = string(v)
			return nil
		})
		return nil
	}); err != nil {
		log.ErrLog(err)
	}
	return keyvalueMap
}

// GetKeyValuesFromBoltByNum can get key-values by number from boltdb.
func GetKeyValuesFromBoltByNum(dbFileName string, bucketName string, num int) (map[string]string, int) {
	i := 0
	keyvalueMap := make(map[string]string, 25)

	allkeyvalueMap := make(map[string]string, 25)
	allkeyvalueMap = GetAllKeyValuesFromBolt(dbFileName, bucketName)

	for k, value := range allkeyvalueMap {
		keyvalueMap[k] = value
		if i++; i >= num {
			break
		}
	}

	return keyvalueMap, i
}

// DelKeyonBolt afford a function whick can delete one key-vlaue in the bucket of db file.
func DelKeyonBolt(dbFileName string, bucketName string, key string) error {
	db, err := bolt.Open(dbFileName, 0666, nil)
	defer db.Close()

	if err != nil {
		log.ErrLog(err)
	}

	if err = db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucketName))
		err := b.Delete([]byte(key))
		return err
	}); err != nil {
		log.ErrLog(err)
	}

	return nil
}

// DelKeysonBolt afford a function whick can delete some key-values in the bucket of db file.
func DelKeysonBolt(dbFileName string, bucketName string, keys []string) error {
	db, err := bolt.Open(dbFileName, 0666, nil)
	defer db.Close()

	if err != nil {
		log.ErrLog(err)
	}

	if err = db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucketName))
		for i := 0; i < len(keys); i++ {
			key := keys[i]
			err := b.Delete([]byte(key))
			return err
		}
		return nil
	}); err != nil {
		log.ErrLog(err)
	}

	return err
}

// DelBucketonBolt can delete a bucket in this db file.
func DelBucketonBolt(dbFileName string, bucketName string) error {
	db, err := bolt.Open(dbFileName, 0666, nil)
	defer db.Close()

	if err != nil {
		log.ErrLog(err)
	}

	if err = db.Update(func(tx *bolt.Tx) error {
		err := tx.DeleteBucket([]byte(bucketName))
		return err
	}); err != nil {
		log.ErrLog(err)
	}

	return err
}

// GetAllKeys can return the []string of all keys.
func GetAllKeys(dbFileName string, BucketName string) []string {
	var keys []string
	db, err := bolt.Open(dbFileName, 0444, nil)
	defer db.Close()

	if err != nil {
		log.ErrLog(err)
	}

	db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(BucketName))
		c := b.Cursor()

		for k, _ := c.First(); k != nil; k, _ = c.Next() {
			keys = append(keys, string(k))
		}
		return nil
	})

	return keys
}

// **************************** //

// BucketCreateChildBucket meaning create child bucket from this bucket directly.
// Maybe not nessiary.
func BucketCreateChildBucket(b *bolt.Bucket, childBucketName string) (*bolt.Bucket, error) {
	childBucket, err := b.CreateBucket([]byte(childBucketName))
	if err != nil {
		return nil, err
	}
	return childBucket, nil
}

// BucketGetChildBucket return a bucket under the provided bucket.
// when there is no such a bucket , return nil
func BucketGetChildBucket(b *bolt.Bucket, childBucketName string) (*bolt.Bucket, error) {
	// b.Bucket will return a new bucket if this childBucket is not exist.
	childBucket := b.Bucket([]byte(childBucketName))
	var err error
	// if can not find the child bucket, then create it now
	if childBucket == nil {
		childBucket, err = b.CreateBucket([]byte(childBucketName))
	}
	if err != nil {
		return nil, err
	}
	return childBucket, nil
}

// OpenBoltDB can get a bolt.DB object.
// return error and close connection when timeout
func OpenBoltDB(dbFileName string, code os.FileMode) (*bolt.DB, error) {
	var db *bolt.DB
	db, err := bolt.Open(dbFileName, code, &bolt.Options{Timeout: ConnectionTimeout})
	return db, err
}

// MustCloseBoltDB can check boltdb to close.
// Return true, nil. when sucessfully, return false, err when fail.
func MustCloseBoltDB(db *bolt.DB) (bool, error) {
	if err := db.Close(); err != nil {
		return false, err
	}
	return true, nil
}

// DBBucketAddKeyValue add key and vlaue in bucket.
// add which is exist has no error, just update it.
func DBBucketAddKeyValue(db *bolt.DB, bucketName, key, value string) error {
	// this one
	tx, err := (*db).Begin(true)
	defer tx.Rollback()

	if err != nil {
		return err
	}

	b := tx.Bucket([]byte(bucketName))
	if err = b.Put([]byte(key), []byte(value)); err != nil {
		return err
	}
	// commit change.
	if err = tx.Commit(); err != nil {
		return err
	}

	return nil
}

// DBBucketAddKeysValues add map[keys]values into bucket
func DBBucketAddKeysValues(db *bolt.DB, bucketName string, kvmap map[string]string) error {
	tx, err := (*db).Begin(true)
	defer tx.Rollback()

	if err != nil {
		return err
	}
	b := tx.Bucket([]byte(bucketName))
	for k, v := range kvmap {
		if err := b.Put([]byte(k), []byte(v)); err != nil {
			return err
		}
	}
	if err = tx.Commit(); err != nil {
		return err
	}

	return nil
}

// DBBucketDelKey can delete one key from bucket
func DBBucketDelKey(db *bolt.DB, bucketName, key string) error {
	tx, err := db.Begin(true)
	defer tx.Rollback()

	if err != nil {
		return err
	}
	b := tx.Bucket([]byte(bucketName))
	if err = b.Delete([]byte(key)); err != nil {
		return err
	}

	if err = tx.Commit(); err != nil {
		return err
	}

	return nil
}

// DBBucketDelKeys can delete a []string keys from bucket
func DBBucketDelKeys(db *bolt.DB, bucketName string, keys []string) error {
	tx, err := db.Begin(true)
	defer tx.Rollback()

	if err != nil {
		return err
	}
	b := tx.Bucket([]byte(bucketName))
	for _, k := range keys {
		if err = b.Delete([]byte(k)); err != nil {
			return err
		}
	}
	if err = tx.Commit(); err != nil {
		return err
	}

	return nil
}

// DBBucketGetKeysByNum can renturn all keys from bucket.
func DBBucketGetKeysByNum(db *bolt.DB, bucketName string, num int) ([]string, error) {
	if num == 0 {
		return []string{}, nil
	}

	var result []string
	tx, err := db.Begin(false)
	defer tx.Rollback()

	if err != nil {
		return nil, err
	}
	b := tx.Bucket([]byte(bucketName))
	c := b.Cursor()
	k, _ := c.First()
	result = append(result, string(k))
	for i := 1; i < num; i++ {
		k, _ := c.Next()
		result = append(result, string(k))
	}

	return result, nil
}

// DBBucketGetValueByKey can return value string from bucket.
func DBBucketGetValueByKey(db *bolt.DB, bucketName string, key string) (string, error) {
	tx, err := db.Begin(false)
	defer tx.Rollback()

	if err != nil {
		return "", err
	}
	b := tx.Bucket([]byte(bucketName))
	v := b.Get([]byte(key))

	return string(v), nil
}

// DBBucketGetKeyValues can return map[keys]values with int from bucket
func DBBucketGetKeyValues(db *bolt.DB, bucketName string, keys []string) (map[string]string, error) {
	result := make(map[string]string)
	tx, err := db.Begin(false)
	defer tx.Rollback()

	if err != nil {
		return nil, err
	}
	b := tx.Bucket([]byte(bucketName))
	for _, k := range keys {
		v := b.Get([]byte(k))
		result[k] = string(v)
	}

	return result, nil
}
