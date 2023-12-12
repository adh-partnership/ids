package database

func (db *Database) Migrate(entities ...interface{}) error {
	for _, entity := range entities {
		if err := db.DB.AutoMigrate(entity); err != nil {
			return err
		}
	}

	return nil
}
