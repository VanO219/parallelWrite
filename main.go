package main

import (
	"bytes"
	"context"
	"fmt"
	"github.com/pkg/errors"
	"log"
	"math/rand"
	"sync"
	"text/template"
	"time"
)

var (
	db                        = DB{}
	wg                        sync.WaitGroup
	gorutineQuantity          int64
	gorutineIterationQuantity int64
	transactionsQuantity      int64
	letterRunes               = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
)

func main() {
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

func genData(trq int64) (query string, err error) {
	defer func() { err = errors.Wrap(err, "main.genData") }()

	type Data struct {
		Model   string
		Company string
		Price   int64
		Date    time.Time
	}

	var d Data
	var strQuery = []string{}

	for i := int64(0); i < trq; i++ {
		d = Data{}
		month := rand.Intn(12-1) + 1
		day := rand.Intn(16-1) + 1
		hour := rand.Intn(24)
		d.Date = time.Date(2021, time.Month(month), day, hour, 0, 0, 0, time.UTC)
		d.Model = randomString(40)
		d.Company = randomString(50)
		d.Price = rand.Int63n(5000000)

		if i != trq-1 {
			strQuery = append(strQuery, fmt.Sprintf("(%s, %s, %d, %s), ", d.Model, d.Company, d.Price, d.Date.Format(time.RFC3339)))
		} else {
			strQuery = append(strQuery, fmt.Sprintf("(%s, %s, %d, %s)", d.Model, d.Company, d.Price, d.Date.Format(time.RFC3339)))
		}
	}

	sqlstr := `
	INSERT INTO public.testtable  
	(model, company, price, date)
	VALUES{{range $i, $j := .}}{{$j}}{{end}}`

	tmp, err := template.New("").Parse(sqlstr)
	if err != nil {
		err = errors.Wrap(err, "failed template parse")
		return "", err
	}

	b := bytes.Buffer{}
	err = tmp.Execute(&b, strQuery)
	if err != nil {
		err = errors.Wrap(err, "failed template execute")
		return "", err
	}

	return b.String(), err
}

func randomString(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}
