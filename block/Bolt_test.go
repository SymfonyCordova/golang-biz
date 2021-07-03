package block

import (
	"fmt"
	"github.com/boltdb/bolt"
	"testing"
)

func TestBoltDb(t *testing.T) {
	//1 打开数据库
	db,err := bolt.Open("test.db",0660,nil)
	if err != nil {
		//log.Panic(fmt.Sprintf("打开数据库失败:%s", err))
		t.Error(fmt.Sprintf("打开数据库失败:%s", err))
	}
	defer func(db *bolt.DB) {
		err := db.Close()
		if err != nil {
			t.Error(fmt.Sprintf("关闭数据库失败:%s", err))
		}
	}(db)

	//tx是事务
	err = db.Update(func(tx *bolt.Tx) error {
		//找到抽屉bucket(如果没有就创建)
		bucket := tx.Bucket([]byte("b1"))
		if bucket == nil{
			//没有抽屉,我们需要创建
			bucket, err = tx.CreateBucket([]byte("b1"))
			if err != nil {
				//log.Panic(fmt.Sprintf("创建bucket(b1)失败:%s", err))
				t.Error(fmt.Sprintf("创建bucket(b1)失败:%s", err))
			}
		}

		//准备写数据
		//err = bucket.Put([]byte("11111"), []byte("hello"))
		//err = bucket.Put([]byte("22222"), []byte("world"))
		//strings test.db => 11111hello22222world
		return nil
	})

	//读数据
	db.View(func(tx *bolt.Tx) error {
		//1 找到抽屉 没有的话直接报错退出
		bucket := tx.Bucket([]byte("b1"))
		if bucket == nil{
			//log.Panic("bucket不应该为空请检查")
			t.Error("bucket不应该为空请检查")
		}

		//2 直接读取数据
		v1 := bucket.Get([]byte("11111"))
		v2 := bucket.Get([]byte("2222"))

		fmt.Printf("v1 %s\n",v1)
		fmt.Printf("v2 %s\n",v2)
		return nil
	})
}
