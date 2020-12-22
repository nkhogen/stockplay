package util

import (
	"bytes"
	"runtime"
	"path/filepath"
	"stockplay/common"
	"strconv"
	"io/ioutil"
	"fmt"
)

func GetDataFilepath(filename string) string {
	_, file, _, _ := runtime.Caller(0)
	return filepath.Join(filepath.Dir(file), "../data", filename)
}

func Shorten(str string) string {
	if len(str) > 40 {
		runes := []rune(str)
		str = string(runes[0:40])
	}
	return str
}

func GenerateHTMLFile(resultList common.ResultDataList) error {
	wb := bytes.Buffer{}
	wb.WriteString("<html>")
	wb.WriteString("<title>Stock Play</title>")
	wb.WriteString("<body>")
	wb.WriteString("<table cellpadding=\"10\">")
	wb.WriteString("<tr>")
	wb.WriteString("<td>Symbol</td>")
	wb.WriteString("<td>Company</td>")
	wb.WriteString("<td>Industry</td>")
	wb.WriteString("<td>Change SMA</td>")
	wb.WriteString("<td>Change VMA</td>")
	wb.WriteString("<td>Price SMA</td>")
	wb.WriteString("<td>Price VMA</td>")
	wb.WriteString("<td>Count</td>")
	wb.WriteString("<td>Load IDs</td>")
	wb.WriteString("</tr>")
	for _, result := range resultList {
		wb.WriteString("<tr>")
		wb.WriteString("<td><a href=")
		wb.WriteString(result.URL)
		wb.WriteString(" target=_blank>")
		wb.WriteString(result.Symbol)
		wb.WriteString("</a></td>")
		wb.WriteString("<td>")
		wb.WriteString(Shorten(result.Company))
		wb.WriteString("</td>")
		wb.WriteString("<td>")
		wb.WriteString(Shorten(result.Industry))
		wb.WriteString("</td>")
		wb.WriteString("<td>")
		wb.WriteString(fmt.Sprintf("%f%%",result.ChangeSMA))
		wb.WriteString("</td>")
		wb.WriteString("<td>")
		wb.WriteString(fmt.Sprintf("%f%%",result.ChangeVMA))
		wb.WriteString("</td>")
		wb.WriteString("<td>")
		wb.WriteString(fmt.Sprintf("%f",result.PriceSMA))
		wb.WriteString("</td>")
		wb.WriteString("<td>")
		wb.WriteString(fmt.Sprintf("%f",result.PriceVMA))
		wb.WriteString("</td>")
		wb.WriteString("<td>")
		wb.WriteString(strconv.Itoa(result.Count))
		wb.WriteString("</td>")
		wb.WriteString("<td>")
		wb.WriteString(result.LoadIDs)
		wb.WriteString("</td>")
		wb.WriteString("</tr>")
	}
	wb.WriteString("</table>")
	wb.WriteString("</body>")
	wb.WriteString("</html>")
	data := []byte(wb.String())
	return ioutil.WriteFile(GetDataFilepath("output.html"), data, 0755)
}