package db

var BucketTemplate = bucketTemplate
var RecordTemplate = recordTemplate

var SortRecords = sortRecords

func (db *BoltAPI) GetCurrentBucket() []string {
	return db.currentBucket
}

func (db *BoltAPI) ClearPath() {
	db.currentBucket = []string{}
}