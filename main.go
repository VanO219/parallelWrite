package main

import (
	"bytes"
	"context"
	"database/sql"
	"fmt"
	"github.com/pkg/errors"
	"log"
	"math/rand"
	"sync"
	"text/template"
	"time"
)

var (
	db                         = DB{}
	wg                         sync.WaitGroup
	goroutineQuantity          = int64(0)
	goroutineIterationQuantity = int64(0)
	transactionsQuantity       = int64(0)
	letterRunes                = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
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
	t0 := time.Now()
	log.Println("старт горутин")
	for i := int64(0); i < goroutineQuantity; i++ {
		go func() {
			err = writeData(ctx, db.DB1, goroutineIterationQuantity, transactionsQuantity, i)
			if err != nil {
				log.Println(err)
				return
			}
		}()
	}

	wg.Wait()
	t1 := time.Now()
	fmt.Println("Время выполнения программы", t1.Sub(t0))
}

func writeData(ctx context.Context, db *sql.DB, iterationQuantity int64, transactionsQuantity int64, id int64) (err error) {
	defer wg.Done()
	defer func() { err = errors.Wrap(err, "main.writeData") }()
	log.Println(id, "goroutine start")
	for i := int64(0); i < iterationQuantity; i++ {
		select {
		case <-ctx.Done():
			return
		default:
			str, err := genData(transactionsQuantity)
			if err != nil {
				return err
			}
			_, err = db.Exec(str)
			if err != nil {
				err = errors.Wrap(err, "failed query exec")
				return
			}
		}
	}
	log.Println(id, "goroutine stop")
	return err
}

func genData(transactionsQuantity int64) (query string, err error) {
	defer func() { err = errors.Wrap(err, "main.genData") }()

	type Data struct {
		Model   string
		Company string
		Price   int64
		Date    time.Time
	}

	var d Data
	var strQuery = []string{}

	for i := int64(0); i < transactionsQuantity; i++ {
		d = Data{}
		month := rand.Intn(12-1) + 1
		day := rand.Intn(16-1) + 1
		hour := rand.Intn(24)
		d.Date = time.Date(2021, time.Month(month), day, hour, 0, 0, 0, time.UTC)
		d.Model = randomString(40)
		d.Company = randomString(50)
		d.Price = rand.Int63n(5000000)

		if i != transactionsQuantity-1 {
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
