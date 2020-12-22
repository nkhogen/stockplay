package util
import (
	"io"
	"golang.org/x/net/html"
	"fmt"
	"net/http"
	"encoding/json"
	"strconv"
	"strings"
	"stockplay/common"
	"context"
)
/*
 1
 DISCB
 Discovery, Inc.
 Communication Services
 Entertainment
 USA
 377.39M
 23.40
 57.95
 68.84%
 218,052
 */

var (
	replacer = strings.NewReplacer(
		"M", "",
		",", "",
		"%", "",
		"_", "0",
		"-", "0",
		"B", "000",
	)
)
type InternalData common.Data

func (iData *InternalData) UnmarshalJSON(data []byte) error {
	fields := []string{}
	err := json.Unmarshal(data, &fields)
	if err != nil {
		return err
	}
	for i, v := range fields {
		switch i {
		case 0:
			v = replacer.Replace(v)
			iVal, err := strconv.Atoi(v)
			if err != nil {
				return err
			}
			iData.Serial = iVal
		case 1: 
			iData.Symbol = v
		case 2: 
			iData.Company = v 
		case 3: 
			iData.Sector = v
		case 4: 
			iData.Industry = v
		case 5: 
			iData.Country = v
		case 6:
			v = replacer.Replace(v)
			fVal, err := strconv.ParseFloat(v, 32)
			if err != nil {
				return err
			}
			iData.MCap = fVal
		case 7:
			v = replacer.Replace(v)
			fVal, err := strconv.ParseFloat(v, 32)
			if err != nil {
				return err
			}
			iData.PE = fVal
		case 8:
			v = replacer.Replace(v)
			fVal, err := strconv.ParseFloat(v, 32)
			if err != nil {
				return err
			}
			iData.Price = fVal
		case 9:
			v = replacer.Replace(v)
			fVal, err := strconv.ParseFloat(v, 32)
			if err != nil {
				return err
			}
			iData.Change = fVal
		case 10:
			v = replacer.Replace(v)
			fVal, err := strconv.ParseFloat(v, 32)
			if err != nil {
				return err
			}
			iData.Volume = fVal
		}
	}
	return nil
}

func getAttr(tokenizer *html.Tokenizer, key string) string {
	var attrVal string
	for {
		currKey, attr, hasMore := tokenizer.TagAttr()
		if string(currKey) == key {
			attrVal = string(attr)
		}
		if !hasMore {
			break
		}
	}
	return attrVal
}

func ReadEndpoint(url string, params map[string]string, maxDownloadSymbols int) ([]*common.Data, error) {
	dataList := []*common.Data{}
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return dataList, err
	}
	query := req.URL.Query()
	if params != nil {
		for key, val := range params {
			query.Set(key, val)
		}
	}
	pageIndex := 0
	totalCount := 0
	for {
		query.Set("r", strconv.Itoa(pageIndex))
		req.URL.RawQuery = query.Encode()
		resp, err := client.Do(req)
		if err != nil {
			return dataList, err
		}
		count := 0
		err = ParseHTML(context.TODO(), resp.Body, func(ctx context.Context, fields []string) (bool, error) {
			if maxDownloadSymbols > 0 && totalCount >= maxDownloadSymbols {
				return false, nil
			}
			b, _:= json.Marshal(fields)
			iData := &InternalData{}
			err := json.Unmarshal(b, iData)
			if err != nil {
				fmt.Printf("Error: %v", err.Error())
				return false, err
			}
			data := common.Data(*iData)
			dataListLen := len(dataList)
			// Last is duplicated
			if dataListLen > 0 && dataList[dataListLen-1].Serial == data.Serial{
				return true, nil
			}
			dataList = append(dataList, &data)
			count++
			totalCount++
			return true, err
		})
		if err != nil {
			fmt.Printf("Error: %+v\n", err.Error())
			return dataList, err
		}
		if count == 0 {
			fmt.Printf("\nTotal symbols: %d\n", totalCount)
			break
		}
		pageIndex = pageIndex + count + 1
	}
	return dataList, nil 
}

func ParseHTML(ctx context.Context, r io.Reader, callback func(ctx context.Context, fields[]string) (bool, error)) error {
	z := html.NewTokenizer(r)
	depth := 0
	inCol := false
	inRow := false
	inTable := false
	colDepth := 0
	rowDepth := 0
	tableDepth := 0
	fields := []string{}
	for {
		tt := z.Next()
		switch tt {
		case html.ErrorToken:
			err := z.Err()
			if err == io.EOF {
				return nil
			}
			return err
		case html.TextToken:
			if depth > 0  && inCol {
				textData := z.Text()
				fields = append(fields, string(textData))
			}
		case html.StartTagToken:
			tn, _ := z.TagName()
			key, attr, _ := z.TagAttr()
			tagName := string(tn)
			keyName := string(key)
			attrName := string(attr)
			if !inTable && tagName == "td" && keyName == "class" && attrName == "table-top" {
				tableDepth++
				inTable = true
			}
			if tagName == "tr" {
				if inTable {
					tableDepth++
					if inRow {
						rowDepth++
					} else {
						attr := getAttr(z, "class")
						if  attr == "table-dark-row-cp" || attr == "table-light-row-cp"{
							rowDepth++
							inRow = true
							fields = []string{}
						}
					}
				}
			} else if tagName == "td" {
				if inRow {
					if inCol {
						colDepth++
					} else {
						attr := getAttr(z, "class")
						if attr == "screener-body-table-nw" {
							colDepth++
							// Column started
							inCol = true
						}
					}
				}
			}
			if len(tn) == 1 && tn[0] == 'a' {
				depth++
			}
		case html.EndTagToken:
			tn, _ := z.TagName()
			tagName := string(tn)
			if tagName == "tr" {
				if inRow {
					rowDepth--
					if rowDepth == 0 {
						inRow = false
						if len(fields) > 0 {
							isContinue, err := callback(ctx, fields)
							if err != nil {
								return err
							}
							if !isContinue {
								return nil
							}
						}
					}
				}
				if inTable {
					tableDepth--
					if tableDepth == 0 {
						inTable = false
					}
				}
			} else if tagName == "td" {
				if inCol {
					colDepth--
					if colDepth == 0 {
						inCol = false
					}
				}
			}
			if len(tn) == 1 && tn[0] == 'a' {
				depth--
			}
		}
	}
	return nil
}