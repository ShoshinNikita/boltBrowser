package db

var (
	BucketTemplate = bucketTemplate
	RecordTemplate = recordTemplate

	SortRecords = sortRecords
)

func (db *BoltAPI) GetCurrentBucketsPath() []string {
	return db.currentBucket
}

func (db *BoltAPI) ClearPath() {
	db.currentBucket = []string{}
}
