package services

import (
	"reflect"
	"sync"
)

type Token struct {
	sync.Mutex
	name    string
	address string
	price   string
	change  string
	min     string
	max     string
	period  string
}

func (c *Token) Get() string {
	c.Lock()
	defer c.Unlock()
	return c.name + " " + c.address + " " + c.price
}

type Tokens struct {
	sync.Mutex
	data []Token
}

func (c *Tokens) Add(pair Token) {
	c.Lock()
	defer c.Unlock()
	c.data = append(c.data, pair)
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
