package whale

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestExampleScrape(t *testing.T) {
	txList, err := FetchBscTransaction()
	for _, list := range txList {
		fmt.Printf("%+v\n", list)
	}
	assert.Equal(t, nil, err)
}

func TestTime(t *testing.T) {
	today := time.Now().Format("2006-01-02 15:04")
	fmt.Println(today)
}
