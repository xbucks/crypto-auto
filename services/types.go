package services

import (
	"reflect"
	"sync"

	"github.com/hirokimoto/crypto-auto/utils"
)

type Token struct {
	sync.Mutex
	target  string // stable, unstable
	updown  string // up, down
	name    string
	address string
	price   string
	change  string
	min     string
	max     string
	period  string
	swaps   []utils.Swap
}

func (c *Token) Get() string {
	c.Lock()
	defer c.Unlock()
	return c.name + " " + c.address + " " + c.price
}

type Tokens struct {
	sync.Mutex
	data     []Token
	progress int
	total    int
}

func (c *Tokens) Add(pair *Token) {
	c.Lock()
	defer c.Unlock()
	c.data = append(c.data, *pair)
}

func (c *Tokens) Get() []Token {
	c.Lock()
	defer c.Unlock()
	return c.data
}

func (c *Tokens) GetItem(index int, key string) string {
	c.Lock()
	defer c.Unlock()
	r := reflect.ValueOf(c.data[index])
	f := reflect.Indirect(r).FieldByName(key)
	return f.String()
}

func (c *Tokens) GetLength() int {
	c.Lock()
	defer c.Unlock()
	length := len(c.data)
	return length
}

func (c *Tokens) GetProgress() int {
	c.Lock()
	defer c.Unlock()
	return c.progress
}

func (c *Tokens) GetTotal() int {
	c.Lock()
	defer c.Unlock()
	return c.total
}

func (c *Tokens) SetProgress(p int) {
	c.Lock()
	defer c.Unlock()
	c.progress = p
}

func (c *Tokens) SetTotal(p int) {
	c.Lock()
	defer c.Unlock()
	c.total = p
}
