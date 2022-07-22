package whale

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

type Symbol string

var StableCoin = map[string]bool{
	"WBNB":    true,
	"BUSD":    true,
	"USDC":    true,
	"USDT":    true,
	"BNB":     true,
	"BSC-USD": true,
	"FIST":    true,
	"ETH":     true,
	"Cake":    true,
}

const (
	ContractVerifyUrlPrefix = "https://api.gopluslabs.io/api/v1/token_security/56?contract_addresses="
)

//<tr>
//<td><span class="hash-tag text-truncate"><a target="_parent" href="/tx/0x7d2aa1585c34b3319da3b786dd707f21b9de644ec11ca7e944ea8be59e21f181">0x7d2aa1585c34b3319da3b786dd707f21b9de644ec11ca7e944ea8be59e21f181</a></span></td>
//<td class="text-nowrap"><span>11 secs ago</span></td>
//<td class="text-nowrap">
//90.671774000739 <a target="_parent" href="/address/0x55d398326f99059ff775485246999027b3197955">BSC-USD </a>
//</td>
//<td>
//<span class="btn btn-xs btn-icon btn-soft-success rounded-circle"><i class="fas fa-long-arrow-alt-right btn-icon__inner"></i></span>
//</td>
//<td class="text-nowrap">
//107.750257 <a target="_parent" href="/address/0xc9882def23bc42d53895b8361d0b1edc7570bc6a">FIST </a>
//</td>
//<td class="text-nowrap">
//0.8414994 BSC-USD
//</td>
//<td>
//<img src="/images/dex/dex-empty.png" data-toggle="tooltip" title="Unknown" class="u-xs-avatar" />
//</td>
//</tr>
//<tr>

func FetchBscTransaction() (txList []*BSCTransaction, err error) {
	// Request the HTML page.
	data, err := HttpGet("https://bscscan.com/dextracker?ps=100")
	if err != nil {
		log.Fatal(err)
	}

	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(bytes.NewBuffer(data))
	if err != nil {
		log.Fatal(err)
	}

	// Find the review items,tr代表一次交易
	doc.Find("tbody tr").Each(func(i int, s *goquery.Selection) {
		tx := &BSCTransaction{TxTime: time.Now()}
		s.Find("td").Each(func(i int, s *goquery.Selection) {
			// 获取交易hash
			if i == 0 {
				tx.TxId = s.Find("span a ").Text()
			}
			// 获取交易时间
			if i == 1 {
				tx.Time = s.Find("span").Text()
			}
			// 获取maker
			if i == 2 {
				dataArr := strings.Split(strings.TrimSpace(s.Text()), " ")
				if len(dataArr) > 0 {
					amount := strings.Replace(dataArr[0], ",", "", -1)
					tx.MakerAmount, err = strconv.ParseFloat(amount, 10)
				}
				if err != nil {
					fmt.Println("err is:", err)
				}
				a := s.Find("a")
				symbol := a.Text()
				tx.MakerSymbol = strings.TrimSpace(symbol)
				href, exist := a.Attr("href")
				if exist && href != "" {
					hrefArr := strings.Split(href, "/")
					tx.MakerContract = hrefArr[len(hrefArr)-1]
				}
			}

			if i == 4 {
				dataArr := strings.Split(strings.TrimSpace(s.Text()), " ")
				if len(dataArr) > 0 {
					amount := strings.Replace(dataArr[0], ",", "", -1)
					tx.TakerAmount, err = strconv.ParseFloat(amount, 10)
				}
				if err != nil {
					fmt.Println("err is:", err)
				}
				a := s.Find("a")
				symbol := a.Text()
				tx.TakerSymbol = strings.TrimSpace(symbol)
				href, exist := a.Attr("href")
				if exist && href != "" {
					hrefArr := strings.Split(href, "/")
					tx.TakerContract = hrefArr[len(hrefArr)-1]
				}
			}

			if i == 5 {
				data := strings.TrimSpace(s.Text())
				dataArr := strings.Split(data, " ")
				if len(dataArr) == 2 {
					tx.Price, err = strconv.ParseFloat(dataArr[0], 10)
					if err != nil {
						fmt.Println("err is:", err)
					}
					tx.PriceSymbol = dataArr[1]
				}
			}
			txList = append(txList, tx)
		})
	})

	// 删除非本次数据
	if len(txList) == 0 {
		fmt.Println("txList is empty")
		return
	}
	return
}

func VerifyContract(tx *BSCTransaction) (security bool, score *ContractVerifyScore, err error) {
	// 如果买入卖出都是稳定币,则跳过不需要统计
	if StableCoin[tx.MakerSymbol] && StableCoin[tx.TakerSymbol] {
		return
	}

	contract := tx.MakerContract
	if StableCoin[tx.MakerSymbol] {
		contract = tx.TakerContract
	}

	if IsContractHasVerified(contract) {
		return
	}

	url := ContractVerifyUrlPrefix + contract
	data, err := HttpGet(url)
	if err != nil {
		return
	}

	var result = new(ContractVerifyResult)
	err = json.Unmarshal(data, result)
	if err != nil {
		return
	}
	security, score, err = AnalysisContract(contract, result)
	return
}
