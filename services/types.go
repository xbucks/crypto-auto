package services

import (
	"fmt"
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

func (c *Tokens) GetByString() []string {
	c.Lock()
	defer c.Unlock()
	t := []string{}
	for _, token := range c.data {
		t = append(t, fmt.Sprintf("%s %s %s", token.name, token.address, token.price))
	}
	return t
}

func (c *Tokens) GetLength() int {
	c.Lock()
	defer c.Unlock()
	length := len(c.data)
	return length
}
