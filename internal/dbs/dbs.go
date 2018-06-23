// Package dbs is a wrapper for Package db
// Package provides functions for working with map[string]*db.BoltAPI
package dbs

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/ShoshinNikita/boltBrowser/internal/db"
)

// allDB keeps all opened databases. string – the path to the db
var allDB map[string]*db.BoltAPI

// DBInfo consist main info about db
type DBInfo struct {
	Name   string `json:"name"`
	DBPath string `json:"dbPath"`
	Size   int64  `json:"size"`
}

// Init – initializing allDB
func Init() {
	allDB = make(map[string]*db.BoltAPI)
}

// OpenDB is a wrapper for *BoltAPI.Open()
func OpenDB(dbPath string) (dbName string, code int, err error) {
	// Check if db was opened
	if _, ok := allDB[dbPath]; ok {
		return "", http.StatusBadRequest, errors.New("This DB was already opened")
	}

	newDB, err := db.Open(dbPath)
	if err != nil {
		return "", http.StatusInternalServerError, err
	}

	allDB[dbPath] = newDB

	return newDB.Name, http.StatusOK, nil
}

// CreateDB is a wrapper for *BoltAPI.Create()
func CreateDB(path string) (dbName, dbPath string, code int, err error) {
	newDB, err := db.Create(path)
	if err != nil {
		return "", "", http.StatusInternalServerError, err
	}

	allDB[newDB.DBPath] = newDB

	return newDB.Name, newDB.DBPath, http.StatusCreated, nil
}

// CloseDB is a wrapper for *BoltAPI.Close()
func CloseDB(dbPath string) (code int, err error) {
	if _, ok := allDB[dbPath]; !ok {
		return http.StatusBadRequest, errors.New("There's no any db with such path (" + dbPath + ")")
	}

	dbName := allDB[dbPath].Name
	allDB[dbPath].Close()
	delete(allDB, dbPath)

	fmt.Printf("[INFO] DB \"%s\" (%s) was closed\n", dbName, dbPath)
	return http.StatusOK, nil
}

// NextBucket is a wrapper for *BoltAPI.Next()
func NextBucket(dbPath, bucket string) (data db.Data, code int, err error) {
	if _, ok := allDB[dbPath]; !ok {
		return data, http.StatusBadRequest, errors.New("There's no any db with such path (" + dbPath + ")")
	}

	data, err = allDB[dbPath].Next(bucket)
	if err != nil {
		return data, http.StatusInternalServerError, err
	}

	return data, http.StatusOK, nil
}

// PrevBucket is a wrapper for *BoltAPI.Back()
func PrevBucket(dbPath string) (data db.Data, code int, err error) {
	if _, ok := allDB[dbPath]; !ok {
		return data, http.StatusBadRequest, errors.New("There's no any db with such path (" + dbPath + ")")
	}

	data, err = allDB[dbPath].Back()
	if err != nil {
		return data, http.StatusInternalServerError, err
	}

	return data, http.StatusOK, nil
}

// GetRoot is a wrapper for *BoltAPI.GetRoot()
func GetRoot(dbPath string) (data db.Data, code int, err error) {
	if _, ok := allDB[dbPath]; !ok {
		return data, http.StatusBadRequest, errors.New("There's no any db with such path (" + dbPath + ")")
	}

	data, err = allDB[dbPath].GetRoot()
	if err != nil {
		return data, http.StatusInternalServerError, err
	}

	return data, http.StatusOK, nil
}

// GetDBsList returns all opened BoltAPI
func GetDBsList() (list []DBInfo) {
	for _, v := range allDB {
		info := DBInfo{Name: v.Name, DBPath: v.DBPath, Size: v.Size}
		list = append(list, info)
	}

	return list
}

// GetCurrent is a wrapper for *BoltAPI.GetCurrent()
func GetCurrent(dbPath string) (info DBInfo, data db.Data, code int, err error) {
	if _, ok := allDB[dbPath]; !ok {
		return info, data, http.StatusBadRequest, errors.New("There's no any db with such path (" + dbPath + ")")
	}

	info.DBPath = dbPath
	info.Name = allDB[dbPath].Name
	info.Size = allDB[dbPath].Size

	data, err = allDB[dbPath].GetCurrent()
	if err != nil {
		return info, data, http.StatusInternalServerError, err
	}

	return info, data, http.StatusOK, nil
}

// GetNextRecords is a wrapper for *BoltAPI.NextRecords()
func GetNextRecords(dbPath string) (data db.Data, code int, err error) {
	if _, ok := allDB[dbPath]; !ok {
		return data, http.StatusBadRequest, errors.New("There's no any db with such path (" + dbPath + ")")
	}

	data, err = allDB[dbPath].NextRecords()
	if err != nil {
		return data, http.StatusInternalServerError, err
	}

	return data, http.StatusOK, nil
}

// GetPrevRecrods is a wrapper for *BoltAPI.PrevRecords()
func GetPrevRecrods(dbPath string) (data db.Data, code int, err error) {
	if _, ok := allDB[dbPath]; !ok {
		return data, http.StatusBadRequest, errors.New("There's no any db with such path (" + dbPath + ")")
	}

	data, err = allDB[dbPath].PrevRecords()
	if err != nil {
		return data, http.StatusInternalServerError, err
	}

	return data, http.StatusOK, nil
}

// Search is a wrapper for *BoltAPI.SearchRegexp() and *BoltAPI.Search()
func Search(dbPath, mode, text string) (records []db.Record, path string, recordsAmount int, code int, err error) {
	if _, ok := allDB[dbPath]; !ok {
		return records, "", 0, http.StatusBadRequest, errors.New("There's no any db with such path (" + dbPath + ")")
	}

	if mode == "regex" {
		records, path, recordsAmount, err = allDB[dbPath].SearchRegexp(text)
	} else {
		records, path, recordsAmount, err = allDB[dbPath].Search(text)
	}

	if err != nil {
		return records, "", 0, http.StatusInternalServerError, err
	}

	return records, path, recordsAmount, http.StatusOK, nil
}

// AddBucket is a wrapper for *BoltAPI.AddBucket()
func AddBucket(dbPath, bucketName string) (code int, err error) {
	if _, ok := allDB[dbPath]; !ok {
		return http.StatusBadRequest, errors.New("There's no any db with such path (" + dbPath + ")")
	}

	err = allDB[dbPath].AddBucket(bucketName)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	return http.StatusCreated, nil
}

// EditBucketName is a wrapper for *BoltAPI.EditBucketName()
func EditBucketName(dbPath, oldName, newName string) (code int, err error) {
	if _, ok := allDB[dbPath]; !ok {
		return http.StatusBadRequest, errors.New("There's no any db with such path (" + dbPath + ")")
	}

	err = allDB[dbPath].EditBucketName(oldName, newName)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	return http.StatusOK, nil
}

// DeleteBucket is a wrapper for *BoltAPI.DeleteBucket()
func DeleteBucket(dbPath, bucketName string) (code int, err error) {
	if _, ok := allDB[dbPath]; !ok {
		return http.StatusBadRequest, errors.New("There's no any db with such path (" + dbPath + ")")
	}

	err = allDB[dbPath].DeleteBucket(bucketName)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	return http.StatusOK, nil
}

// AddRecord is a wrapper for *BoltAPI.AddRecord()
func AddRecord(dbPath, key, value string) (code int, err error) {
	if _, ok := allDB[dbPath]; !ok {
		return http.StatusBadRequest, errors.New("There's no any db with such path (" + dbPath + ")")
	}

	err = allDB[dbPath].AddRecord(key, value)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	return http.StatusCreated, nil
}

// EditRecord is a wrapper for *BoltAPI.EditRecord()
func EditRecord(dbPath, oldKey, newKey, newValue string) (code int, err error) {
	if _, ok := allDB[dbPath]; !ok {
		return http.StatusBadRequest, errors.New("There's no any db with such path (" + dbPath + ")")
	}

	err = allDB[dbPath].EditRecord(oldKey, newKey, newValue)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	return http.StatusOK, nil
}

// DeleteRecord is a wrapper for *BoltAPI.DeleteRecord()
func DeleteRecord(dbPath, key string) (code int, err error) {
	if _, ok := allDB[dbPath]; !ok {
		return http.StatusBadRequest, errors.New("There's no any db with such path (" + dbPath + ")")
	}

	err = allDB[dbPath].DeleteRecord(key)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	return http.StatusOK, nil
}

// CloseDBs closes all databases
func CloseDBs() {
	for k := range allDB {
		allDB[k].Close()
		delete(allDB, k)
	}
	fmt.Println("[INFO] All databases were closed")
}
