package util

import (
	"encoding/json"
	"stockplay/common"
	"fmt"
	"os"
)

const (
	MaxScans = 50
)
type DataStore struct {
	filepath string
}

func NewDataStore(filename string) (*DataStore, error) {
	filepath := GetDataFilepath(filename)
	fmt.Println(filepath)
	ds := &DataStore{filepath: filepath}
	return ds, nil
}

func (ds *DataStore) Write(dataList []*common.Data) error {
	for i := MaxScans-1; i >= 0 ; i-- {
		iFilepath := fmt.Sprintf("%s-%d.json", ds.filepath, i)
		_, err := os.Stat(iFilepath)
		if os.IsNotExist(err) {
			continue
		}
		newFilePath := fmt.Sprintf("%s-%d.json", ds.filepath, i + 1)
		if i + 1 < MaxScans {
			err = os.Rename(iFilepath, newFilePath)
			if err != nil {
				return err
			}
		}
		os.Remove(iFilepath)
	}
	file, err := os.OpenFile(fmt.Sprintf("%s-0.json", ds.filepath), os.O_APPEND | os.O_RDWR | os.O_CREATE, 0755)
	if err != nil {
		return err
	}
	defer file.Close() 
	encoder := json.NewEncoder(file)
	encoder.SetIndent(" ", " ")
	return encoder.Encode(dataList)
}

func (ds *DataStore) Loads() ([][]*common.Data, error) {
	dataLists := make([][]*common.Data, 0, MaxScans)
	for i := 0; i < MaxScans; i++ {
		dataList := []*common.Data{}
		iFilepath := fmt.Sprintf("%s-%d.json", ds.filepath, i)
		_, err := os.Stat(iFilepath)
		if os.IsNotExist(err) {
			dataLists = append(dataLists, dataList)
			continue
		}
		err = ds.Load(i, &dataList)
		if err != nil {
			return dataLists, err
		}
		dataLists = append(dataLists, dataList)
	}
	return dataLists, nil
}

func (ds *DataStore) Load(idx int, dataList *[]*common.Data) error {
	if dataList == nil {
		return nil
	}
	file, err := os.OpenFile(fmt.Sprintf("%s-%d.json", ds.filepath, idx), os.O_APPEND | os.O_RDWR | os.O_CREATE, 0755)
	if err != nil {
		return err
	}
	defer file.Close()
	decoder := json.NewDecoder(file)
	err = decoder.Decode(dataList)
	if err != nil {
		return err
	}
	for _, data := range *dataList {
		data.LoadID = idx
	}
	return nil
}