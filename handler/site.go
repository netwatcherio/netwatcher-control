package handler

/*func createSite(c *mongo.Database) (bool, error) {
	var siteD = models.Site{
		ID:      primitive.NewObjectID(),
		Name:    "Test Site",
		Members: nil,
	}

	mar, err := bson.Marshal(siteD)
	if err != nil {
		log.Errorf("1 %s", err)
		return false, err
	}
	var b *bson.D
	err = bson.Unmarshal(mar, &b)
	if err != nil {
		log.Errorf("2 %s", err)
		return false, err
	}
	result, err := c.Collection("sites").InsertOne(context.TODO(), b)
	if err != nil {
		log.Errorf("3 %s", err)
		return false, err
	}

	fmt.Printf(" with _id: %v\n", result.InsertedID)
	return true, nil
}

func getSite(id primitive.ObjectID, db *mongo.Database) (*models.Site, error) {
	// if hash is blank, search for pin matching with blank hash
	// if none exist, return error
	// if match, return new agent, and new hash, then let another function update the hash?
	// if hash is included, search both, and return nil for hash, and false for new if verified
	// if hash is included and none match, return err
	var filter = bson.D{{"_id", id}}

	cursor, err := db.Collection("sites").Find(context.TODO(), filter)
	if err != nil {
		return nil, err
	}
	var results []bson.D
	if err = cursor.All(context.TODO(), &results); err != nil {
		return nil, err
	}

	//fmt.Println(results)

	if len(results) > 1 {
		return nil, errors.New("multiple sites match when using id")
	}

	if len(results) == 0 {
		return nil, errors.New("no sites match when using id")
	}

	doc, err := bson.Marshal(&results[0])
	if err != nil {
		log.Errorf("1 %s", err)
		return nil, err
	}

	var site *models.Site
	err = bson.Unmarshal(doc, &site)
	if err != nil {
		log.Errorf("2 %s", err)
		return nil, err
	}

	return site, nil
}*/
