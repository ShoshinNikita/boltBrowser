package db

func (db *BoltAPI) clearPath() {
	db.currentBucket = []string{}
}