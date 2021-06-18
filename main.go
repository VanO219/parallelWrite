package main

import (
	"context"
	"log"
	"math/rand"
	"sync"
	"time"
)

var (
	db = DB{}
	wg sync.WaitGroup
	gorutineQuantity int64
	gorutineIterationQuantity int64
	transactionsQuantity int64
)

func main()  {
	err := db.Connect("ivan", "pass")
	if err != nil {
		log.Fatalln(err)
	}

	rand.Seed(time.Now().UnixNano())

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer func() {
		err := db.Close()
		if err != nil {
			log.Println(err)
		}
	}()
	defer cancel()
}

func genData() (query string)  {
	type Data struct {
		Model string
		Company string
		Price string
		Data string
	}
	return
}
