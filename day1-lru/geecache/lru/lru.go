package lru

import (
	"container/list"
	"sync"
)

/*Notes
类型断言 var a = x.(T)，如果x不是nil，并且可以转换为T类型，则断言成功，并返回T类型的变量x。
如果断言失败则会panic。

双向链表中的元素定义
type Element struct {
	next, prev *Element //指向前后元素的指针
	list       *List    //元素所属链表的指针
	Value      any      //存储值，any是空接口interface{}的别名，空接口可以表示任意数据类型
}
*/

var wg sync.WaitGroup
var rwm sync.RWMutex

// Cache is a LRU cache. It is not safe for concurrent access.
type Cache struct {
	maxBytes  int64                         //允许使用的最大内存
	nbytes    int64                         //当前已使用内存
	ll        *list.List                    //双向链表指针
	cache     map[string]*list.Element      //key:string,value:指向链表元素的指针
	OnEvicted func(key string, value Value) // 可选，删除记录时的回调函数
}

//定义链表元素的结构为entry
type entry struct {
	key   string
	value Value
}

// Value use Len to count how many bytes it takes
type Value interface {
	Len() int // Implement in lru_test.go
}

// New is the Constructor of Cache 初始化一个缓存结构体
func New(maxBytes int64, onEvicted func(string, Value)) *Cache {
	return &Cache{
		maxBytes:  maxBytes,
		ll:        list.New(),                     //内置方法，初始化一个双向链表
		cache:     make(map[string]*list.Element), //初始化map
		OnEvicted: onEvicted,
	}
}

// Add adds a value to the cache. 向缓存中添加一个键值对
func (c *Cache) Add(key string, value Value) {

	if ele, ok := c.cache[key]; ok {
		//访问命中
		c.ll.MoveToFront(ele) //将元素移动到头部
		//计算内存占用
		kv := ele.Value.(*entry) //ele.Value为any类型，转换为entry结构体指针
		c.nbytes += int64(value.Len()) - int64(kv.value.Len())
		kv.value = value
	} else {
		//在头部插入一个新元素
		ele := c.ll.PushFront(&entry{key, value}) //返回新元素指针
		c.cache[key] = ele
		c.nbytes += int64(len(key)) + int64(value.Len())
	}
	//处理超出内存上限的情况
	//先添加后删除 or 先删除后添加
	for c.maxBytes != 0 && c.maxBytes < c.nbytes {
		c.RemoveOldest()
	}
}

// Get look ups a key's value 获取元素值
func (c *Cache) Get(key string) (value Value, ok bool) {
	if ele, ok := c.cache[key]; ok {
		//访问命中，元素移到头部
		c.ll.MoveToFront(ele)
		kv := ele.Value.(*entry)
		return kv.value, true
	}
	return
}

// RemoveOldest removes the oldest item  移除一个最近最少访问的元素
func (c *Cache) RemoveOldest() {
	ele := c.ll.Back()
	if ele != nil {
		//Remove()是container/list的内置方法，删除链表中的元素
		c.ll.Remove(ele)
		kv := ele.Value.(*entry)
		delete(c.cache, kv.key) //内置delet()，对map执行删除操作
		//更新占用内存
		c.nbytes -= int64(len(kv.key)) + int64(kv.value.Len())
		//移除时调用回调函数
		if c.OnEvicted != nil {
			c.OnEvicted(kv.key, kv.value)
		}
	}
}

// Len the number of cache entries
func (c *Cache) Len() int {
	//Len()是container/list的内置方法，返回双向链表的元素个数
	return c.ll.Len()
}
