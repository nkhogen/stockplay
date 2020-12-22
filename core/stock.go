package core

import (
	"stockplay/util"
	"stockplay/common"
	"fmt"
	"github.com/VividCortex/ewma"
	"sort"
	"strings"
	"strconv"
)

const (
	storageFilePrefix = "stock-file"
	baseURL = "https://finviz.com/screener.ashx"
	maxDownloadSymbols = 1000
)

var (
	queryParams = map[string]string {
		"v": "111",
		//"s": "ta_topgainers",
		"o": "-change",
		"f": "geo_usa",
	}
)

type StockPlay struct {
	dataStore *util.DataStore
}

func NewStockPlay() *StockPlay {
	ds, err := util.NewDataStore(storageFilePrefix)
	if err != nil {
		panic(err)
	}
	return &StockPlay{dataStore: ds}
}

func (sp *StockPlay) Download() error {
	dataList, err := util.ReadEndpoint(baseURL, queryParams, maxDownloadSymbols)
	if err != nil {
		return err
	}
	err = sp.dataStore.Write(dataList)
	if err != nil {
		return err
	}
	return nil
}

func (sp *StockPlay) Start() error {
	dataLists, err := sp.dataStore.Loads()
	if err != nil {
		return err
	}
	symbolData := map[string][]*common.Data{}
	for _, dataList := range dataLists {
		for i := range dataList {
			data := dataList[i]
			symbolData[data.Symbol] = append(symbolData[data.Symbol], data)
		}
	}
	resultDataList := common.ResultDataList{}
	for symbol, dataList := range symbolData {
		sChange := ewma.NewMovingAverage()  //=> Returns a SimpleEWMA if called without params
		vChange := ewma.NewMovingAverage(2) //=> returns a VariableEWMA with a decay of 2 / (5 + 1)
		sPrice := ewma.NewMovingAverage()  //=> Returns a SimpleEWMA if called without params
		vPrice := ewma.NewMovingAverage(2) //=> returns a VariableEWMA with a decay of 2 / (5 + 1)
		loadIDs := []string{}
		resultData := &common.ResultData{}
		for i := range dataList {
			data := dataList[i]
			if data.Price < 0.01 {
			} else {
				resultData.Company = data.Company
				resultData.Industry = data.Industry
			}
			loadIDs = append(loadIDs, strconv.Itoa(data.LoadID))
			sChange.Add(data.Change)
			vChange.Add(data.Change)
			sPrice.Add(data.Price)
			vPrice.Add(data.Price)
		}
		resultData.Symbol = symbol
		resultData.ChangeSMA = sChange.Value()
		resultData.ChangeVMA = vChange.Value()
		resultData.PriceSMA = sPrice.Value()
		resultData.PriceVMA = vPrice.Value()
		resultData.Count = len(dataList)
		resultData.URL = fmt.Sprintf("https://finance.yahoo.com/quote/%s", symbol)
		resultData.LoadIDs = strings.Join(loadIDs, " ")
		resultDataList = append(resultDataList, resultData)
	}
	sort.Sort(resultDataList)
	err = util.GenerateHTMLFile(resultDataList)
	if err != nil {
		return err
	}
	return nil
}