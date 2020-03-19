package meilisearch

import (
	"reflect"
	"testing"
)

func TestClientSettings_GetAll(t *testing.T) {
	var indexUID = "TestClientSettings_GetAll"

	var client = NewClient(Config{
		Host: "http://localhost:7700",
	})

	_, err := client.Indexes().Create(CreateIndexRequest{
		UID: indexUID,
	})

	if err != nil {
		t.Fatal(err)
	}

	settingsRes, err := client.Settings(indexUID).GetAll()

	if err != nil {
		t.Fatal(err)
	}

	expected := Settings{
		RankingRules:         []string{"typo", "words", "proximity", "attribute", "wordsPosition", "exactness"},
		DistinctAttribute:    nil,
		SearchableAttributes: []string{},
		DisplayedAttributes:  []string{},
		StopWords:            []string{},
		Synonyms:             map[string][]string{},
		AcceptNewFields:      true,
	}

	if !reflect.DeepEqual(*settingsRes, expected) {
		t.Fatalf("The response body is not good %v vs %v", settingsRes, expected)
	}
}

func TestClientSettings_AddOrUpdateAll(t *testing.T) {
	var indexUID = "TestClientSettings_AddOrUpdateAll"

	var client = NewClient(Config{
		Host: "http://localhost:7700",
	})

	_, err := client.Indexes().Create(CreateIndexRequest{
		UID: indexUID,
	})

	if err != nil {
		t.Fatal(err)
	}

	settings := Settings{
		RankingRules:         []string{"typo", "words", "proximity", "attribute", "wordsPosition", "exactness"},
		SearchableAttributes: []string{"id", "title", "description"},
		DisplayedAttributes:  []string{"id", "title", "description"},
		StopWords:            []string{"a", "the"},
		Synonyms: map[string][]string{
			"car": []string{"automobile"},
		},
		AcceptNewFields: false,
	}

	updateIDRes, err := client.Settings(indexUID).AddOrUpdateAll(settings)

	if err != nil {
		t.Fatal(err)
	}

	client.AwaitAsyncUpdateID(indexUID, updateIDRes)
}

func TestClientSettings_ResetAll(t *testing.T) {
	var indexUID = "TestClientSettings_ResetAll"

	var client = NewClient(Config{
		Host: "http://localhost:7700",
	})

	_, err := client.Indexes().Create(CreateIndexRequest{
		UID: indexUID,
	})

	if err != nil {
		t.Fatal(err)
	}

	updateIDRes, err := client.Settings(indexUID).ResetAll()

	if err != nil {
		t.Fatal(err)
	}

	client.AwaitAsyncUpdateID(indexUID, updateIDRes)

}

func TestClientSettings_GetRankingRules(t *testing.T) {
	var indexUID = "TestClientSettings_GetRankingRules"

	var client = NewClient(Config{
		Host: "http://localhost:7700",
	})

	_, err := client.Indexes().Create(CreateIndexRequest{
		UID: indexUID,
	})

	if err != nil {
		t.Fatal(err)
	}

	rankingRulesRes, err := client.Settings(indexUID).GetRankingRules()

	if err != nil {
		t.Fatal(err)
	}

	expected := []string{"typo", "words", "proximity", "attribute", "wordsPosition", "exactness"}

	if !reflect.DeepEqual(*rankingRulesRes, expected) {
		t.Fatal("The response body is not good")
	}
}

func TestClientSettings_SetRankingRules(t *testing.T) {
	var indexUID = "TestClientSettings_SetRankingRules"

	var client = NewClient(Config{
		Host: "http://localhost:7700",
	})

	_, err := client.Indexes().Create(CreateIndexRequest{
		UID: indexUID,
	})

	if err != nil {
		t.Fatal(err)
	}

	rankingRules := []string{"typo", "words", "proximity", "attribute", "wordsPosition", "exactness"}

	updateIDRes, err := client.Settings(indexUID).SetRankingRules(rankingRules)

	if err != nil {
		t.Fatal(err)
	}

	client.AwaitAsyncUpdateID(indexUID, updateIDRes)
}

func TestClientSettings_ResetRankingRules(t *testing.T) {
	var indexUID = "TestClientSettings_ResetRankingRules"

	var client = NewClient(Config{
		Host: "http://localhost:7700",
	})

	_, err := client.Indexes().Create(CreateIndexRequest{
		UID: indexUID,
	})

	if err != nil {
		t.Fatal(err)
	}

	updateIDRes, err := client.Settings(indexUID).ResetRankingRules()

	if err != nil {
		t.Fatal(err)
	}

	client.AwaitAsyncUpdateID(indexUID, updateIDRes)
}

func TestClientSettings_GetDistinctAttribute(t *testing.T) {
	var indexUID = "TestClientSettings_GetDistinctAttribute"

	var client = NewClient(Config{
		Host: "http://localhost:7700",
	})

	_, err := client.Indexes().Create(CreateIndexRequest{
		UID: indexUID,
	})

	if err != nil {
		t.Fatal(err)
	}

	distinctAttributeRes, err := client.Settings(indexUID).GetDistinctAttribute()

	if err != nil {
		t.Fatal(err)
	}

	if reflect.DeepEqual(*distinctAttributeRes, nil) {
		t.Fatal("The response body is not good")
	}
}

func TestClientSettings_SetDistinctAttribute(t *testing.T) {
	var indexUID = "TestClientSettings_SetDistinctAttribute"

	var client = NewClient(Config{
		Host: "http://localhost:7700",
	})

	_, err := client.Indexes().Create(CreateIndexRequest{
		UID: indexUID,
	})

	if err != nil {
		t.Fatal(err)
	}

	updateIDRes, err := client.Settings(indexUID).SetDistinctAttribute("skuid")

	if err != nil {
		t.Fatal(err)
	}

	client.AwaitAsyncUpdateID(indexUID, updateIDRes)

}

func TestClientSettings_ResetDistinctAttribute(t *testing.T) {
	var indexUID = "TestClientSettings_ResetDistinctAttribute"

	var client = NewClient(Config{
		Host: "http://localhost:7700",
	})

	_, err := client.Indexes().Create(CreateIndexRequest{
		UID: indexUID,
	})

	if err != nil {
		t.Fatal(err)
	}

	updateIDRes, err := client.Settings(indexUID).ResetDistinctAttribute()

	if err != nil {
		t.Fatal(err)
	}

	client.AwaitAsyncUpdateID(indexUID, updateIDRes)

}

func TestClientSettings_GetSearchableAttributes(t *testing.T) {
	var indexUID = "TestClientSettings_GetSearchableAttributes"

	var client = NewClient(Config{
		Host: "http://localhost:7700",
	})

	_, err := client.Indexes().Create(CreateIndexRequest{
		UID: indexUID,
	})

	if err != nil {
		t.Fatal(err)
	}

	searchableAttributesRes, err := client.Settings(indexUID).GetSearchableAttributes()

	if err != nil {
		t.Fatal(err)
	}

	if len(*searchableAttributesRes) > 0 {
		t.Fatal("The response body is not good")
	}
}

func TestClientSettings_SetSearchableAttributes(t *testing.T) {
	var indexUID = "TestClientSettings_SetSearchableAttributes"

	var client = NewClient(Config{
		Host: "http://localhost:7700",
	})

	_, err := client.Indexes().Create(CreateIndexRequest{
		UID: indexUID,
	})

	if err != nil {
		t.Fatal(err)
	}

	searchableAttributes := []string{"id", "title", "description"}

	updateIDRes, err := client.Settings(indexUID).SetSearchableAttributes(searchableAttributes)

	if err != nil {
		t.Fatal(err)
	}

	client.AwaitAsyncUpdateID(indexUID, updateIDRes)
}

func TestClientSettings_ResetSearchableAttributes(t *testing.T) {
	var indexUID = "TestClientSettings_ResetSearchableAttributes"

	var client = NewClient(Config{
		Host: "http://localhost:7700",
	})

	_, err := client.Indexes().Create(CreateIndexRequest{
		UID: indexUID,
	})

	if err != nil {
		t.Fatal(err)
	}

	updateIDRes, err := client.Settings(indexUID).ResetSearchableAttributes()

	if err != nil {
		t.Fatal(err)
	}

	client.AwaitAsyncUpdateID(indexUID, updateIDRes)
}

func TestClientSettings_GetDisplayedAttributes(t *testing.T) {
	var indexUID = "TestClientSettings_GetDisplayedAttributes"

	var client = NewClient(Config{
		Host: "http://localhost:7700",
	})

	_, err := client.Indexes().Create(CreateIndexRequest{
		UID: indexUID,
	})

	if err != nil {
		t.Fatal(err)
	}

	displayedAttributesRes, err := client.Settings(indexUID).GetDisplayedAttributes()

	if err != nil {
		t.Fatal(err)
	}

	if len(*displayedAttributesRes) > 0 {
		t.Fatal("The response body is not good")
	}
}

func TestClientSettings_SetDisplayedAttributes(t *testing.T) {
	var indexUID = "TestClientSettings_SetDisplayedAttributes"

	var client = NewClient(Config{
		Host: "http://localhost:7700",
	})

	_, err := client.Indexes().Create(CreateIndexRequest{
		UID: indexUID,
	})

	if err != nil {
		t.Fatal(err)
	}

	displayedAttributes := []string{"id", "title", "description"}

	updateIDRes, err := client.Settings(indexUID).SetDisplayedAttributes(displayedAttributes)

	if err != nil {
		t.Fatal(err)
	}

	client.AwaitAsyncUpdateID(indexUID, updateIDRes)
}

func TestClientSettings_ResetDisplayedAttributes(t *testing.T) {
	var indexUID = "TestClientSettings_ResetDisplayedAttributes"

	var client = NewClient(Config{
		Host: "http://localhost:7700",
	})

	_, err := client.Indexes().Create(CreateIndexRequest{
		UID: indexUID,
	})

	if err != nil {
		t.Fatal(err)
	}

	updateIDRes, err := client.Settings(indexUID).ResetDisplayedAttributes()

	if err != nil {
		t.Fatal(err)
	}

	client.AwaitAsyncUpdateID(indexUID, updateIDRes)
}

func TestClientSettings_GetStopWords(t *testing.T) {
	var indexUID = "TestClientSettings_GetStopWords"

	var client = NewClient(Config{
		Host: "http://localhost:7700",
	})

	_, err := client.Indexes().Create(CreateIndexRequest{
		UID: indexUID,
	})

	if err != nil {
		t.Fatal(err)
	}

	stopWordsRes, err := client.Settings(indexUID).GetStopWords()

	if err != nil {
		t.Fatal(err)
	}

	if len(*stopWordsRes) > 0 {
		t.Fatal("The response body is not good")
	}

}

func TestClientSettings_SetStopWords(t *testing.T) {
	var indexUID = "TestClientSettings_SetStopWords"

	var client = NewClient(Config{
		Host: "http://localhost:7700",
	})

	_, err := client.Indexes().Create(CreateIndexRequest{
		UID: indexUID,
	})

	if err != nil {
		t.Fatal(err)
	}

	stopWords := []string{"a", "the"}

	updateIDRes, err := client.Settings(indexUID).SetStopWords(stopWords)

	if err != nil {
		t.Fatal(err)
	}

	client.AwaitAsyncUpdateID(indexUID, updateIDRes)

}

func TestClientSettings_ResetStopWords(t *testing.T) {
	var indexUID = "TestClientSettings_ResetStopWords"

	var client = NewClient(Config{
		Host: "http://localhost:7700",
	})

	_, err := client.Indexes().Create(CreateIndexRequest{
		UID: indexUID,
	})

	if err != nil {
		t.Fatal(err)
	}

	updateIDRes, err := client.Settings(indexUID).ResetStopWords()

	if err != nil {
		t.Fatal(err)
	}

	client.AwaitAsyncUpdateID(indexUID, updateIDRes)

}

func TestClientSettings_GetSynonyms(t *testing.T) {
	var indexUID = "TestClientSettings_GetSynonyms"

	var client = NewClient(Config{
		Host: "http://localhost:7700",
	})

	_, err := client.Indexes().Create(CreateIndexRequest{
		UID: indexUID,
	})

	if err != nil {
		t.Fatal(err)
	}

	synonymsRes, err := client.Settings(indexUID).GetSynonyms()

	if err != nil {
		t.Fatal(err)
	}

	if len(*synonymsRes) > 0 {
		t.Fatal("The response body is not good")
	}

}

func TestClientSettings_SetSynonyms(t *testing.T) {
	var indexUID = "TestClientSettings_SetSynonyms"

	var client = NewClient(Config{
		Host: "http://localhost:7700",
	})

	_, err := client.Indexes().Create(CreateIndexRequest{
		UID: indexUID,
	})

	if err != nil {
		t.Fatal(err)
	}

	synonyms := map[string][]string{
		"car": []string{"automobile"},
	}

	updateIDRes, err := client.Settings(indexUID).SetSynonyms(synonyms)

	if err != nil {
		t.Fatal(err)
	}

	client.AwaitAsyncUpdateID(indexUID, updateIDRes)
}

func TestClientSettings_ResetSynonyms(t *testing.T) {
	var indexUID = "TestClientSettings_ResetSynonyms"

	var client = NewClient(Config{
		Host: "http://localhost:7700",
	})

	_, err := client.Indexes().Create(CreateIndexRequest{
		UID: indexUID,
	})

	if err != nil {
		t.Fatal(err)
	}

	updateIDRes, err := client.Settings(indexUID).ResetSynonyms()

	if err != nil {
		t.Fatal(err)
	}

	client.AwaitAsyncUpdateID(indexUID, updateIDRes)

}

func TestClientSettings_GetAcceptNewFields(t *testing.T) {
	var indexUID = "TestClientSettings_GetAcceptNewFields"

	var client = NewClient(Config{
		Host: "http://localhost:7700",
	})

	_, err := client.Indexes().Create(CreateIndexRequest{
		UID: indexUID,
	})

	if err != nil {
		t.Fatal(err)
	}

	acceptNewFieldsRes, err := client.Settings(indexUID).GetAcceptNewFields()

	if err != nil {
		t.Fatal(err)
	}

	if !*acceptNewFieldsRes {
		t.Fatal("The response body is not good")
	}
}

func TestClientSettings_SetAcceptNewFields(t *testing.T) {
	var indexUID = "TestClientSettings_SetAcceptNewFields"

	var client = NewClient(Config{
		Host: "http://localhost:7700",
	})

	_, err := client.Indexes().Create(CreateIndexRequest{
		UID: indexUID,
	})

	if err != nil {
		t.Fatal(err)
	}

	updateIDRes, err := client.Settings(indexUID).SetAcceptNewFields(false)

	if err != nil {
		t.Fatal(err)
	}

	client.AwaitAsyncUpdateID(indexUID, updateIDRes)

}
