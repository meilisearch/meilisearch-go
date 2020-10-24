package meilisearch

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

func TestClientSettings_GetAll(t *testing.T) {
	var indexUID = "TestClientSettings_GetAll"

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
		RankingRules:          []string{"typo", "words", "proximity", "attribute", "wordsPosition", "exactness"},
		DistinctAttribute:     nil,
		SearchableAttributes:  []string{"*"},
		DisplayedAttributes:   []string{"*"},
		StopWords:             []string{},
		Synonyms:              map[string][]string(nil),
		AttributesForFaceting: []string{},
	}

	assert.Equal(t, *settingsRes, expected)
}

func TestClientSettings_UpdateAll(t *testing.T) {
	var indexUID = "TestClientSettings_UpdateAll"

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
			"car": {"automobile"},
		},
		AttributesForFaceting: []string{"title"},
	}

	updateIDRes, err := client.Settings(indexUID).UpdateAll(settings)

	if err != nil {
		t.Fatal(err)
	}

	_, _ = client.DefaultWaitForPendingUpdate(indexUID, updateIDRes)
}

func TestClientSettings_ResetAll(t *testing.T) {
	var indexUID = "TestClientSettings_ResetAll"

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

	_, _ = client.DefaultWaitForPendingUpdate(indexUID, updateIDRes)

}

func TestClientSettings_GetRankingRules(t *testing.T) {
	var indexUID = "TestClientSettings_GetRankingRules"

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

	assert.Equal(t, expected, *rankingRulesRes)
}

func TestClientSettings_UpdateRankingRules(t *testing.T) {
	var indexUID = "TestClientSettings_UpdateRankingRules"

	_, err := client.Indexes().Create(CreateIndexRequest{
		UID: indexUID,
	})

	if err != nil {
		t.Fatal(err)
	}

	rankingRules := []string{"typo", "words", "proximity", "attribute", "wordsPosition", "exactness"}

	updateIDRes, err := client.Settings(indexUID).UpdateRankingRules(rankingRules)

	if err != nil {
		t.Fatal(err)
	}

	_, _ = client.DefaultWaitForPendingUpdate(indexUID, updateIDRes)
}

func TestClientSettings_ResetRankingRules(t *testing.T) {
	var indexUID = "TestClientSettings_ResetRankingRules"

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

	_, _ = client.DefaultWaitForPendingUpdate(indexUID, updateIDRes)
}

func TestClientSettings_GetDistinctAttribute(t *testing.T) {
	var indexUID = "TestClientSettings_GetDistinctAttribute"

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

	assert.Equal(t, EmptyStr(), *distinctAttributeRes)
}

func TestClientSettings_UpdateDistinctAttribute(t *testing.T) {
	var indexUID = "TestClientSettings_UpdateDistinctAttribute"

	_, err := client.Indexes().Create(CreateIndexRequest{
		UID: indexUID,
	})

	if err != nil {
		t.Fatal(err)
	}

	updateIDRes, err := client.Settings(indexUID).UpdateDistinctAttribute("skuid")

	if err != nil {
		t.Fatal(err)
	}

	_, _ = client.DefaultWaitForPendingUpdate(indexUID, updateIDRes)

}

func TestClientSettings_ResetDistinctAttribute(t *testing.T) {
	var indexUID = "TestClientSettings_ResetDistinctAttribute"

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

	_, _ = client.DefaultWaitForPendingUpdate(indexUID, updateIDRes)

}

func TestClientSettings_GetSearchableAttributes(t *testing.T) {
	var indexUID = "TestClientSettings_GetSearchableAttributes"

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

	if len(*searchableAttributesRes) != 1 {
		t.Fatal("Wrong response for searchableAttributes")
	}

	searchableAttibutesString := *searchableAttributesRes
	expectedSearchableAttributesString := "*"

	if searchableAttibutesString[0] != expectedSearchableAttributesString {
		t.Fatalf(
			"Wrong response for searchableAttributes, expected %s, got %s\n",
			searchableAttibutesString,
			expectedSearchableAttributesString,
		)
	}
}

func TestClientSettings_UpdateSearchableAttributes(t *testing.T) {
	var indexUID = "TestClientSettings_UpdateSearchableAttributes"

	_, err := client.Indexes().Create(CreateIndexRequest{
		UID: indexUID,
	})

	if err != nil {
		t.Fatal(err)
	}

	searchableAttributes := []string{"id", "title", "description"}

	updateIDRes, err := client.Settings(indexUID).UpdateSearchableAttributes(searchableAttributes)

	if err != nil {
		t.Fatal(err)
	}

	_, _ = client.DefaultWaitForPendingUpdate(indexUID, updateIDRes)
}

func TestClientSettings_ResetSearchableAttributes(t *testing.T) {
	var indexUID = "TestClientSettings_ResetSearchableAttributes"

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

	_, _ = client.DefaultWaitForPendingUpdate(indexUID, updateIDRes)

	searchableAttributesRes, err := client.Settings(indexUID).GetSearchableAttributes()

	if err != nil {
		t.Fatal(err)
	}

	if len(*searchableAttributesRes) != 1 {
		t.Fatal("Wrong response for searchableAttributes after reset")
	}

	searchableAttibutesString := *searchableAttributesRes
	expectedSearchableAttributesString := "*"

	if searchableAttibutesString[0] != expectedSearchableAttributesString {
		t.Fatalf(
			"Wrong response for searchableAttributes after reset, expected %s, got %s\n",
			searchableAttibutesString,
			expectedSearchableAttributesString,
		)
	}
}

func TestClientSettings_GetDisplayedAttributes(t *testing.T) {
	var indexUID = "TestClientSettings_GetDisplayedAttributes"

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

	if len(*displayedAttributesRes) != 1 {
		t.Fatal("Wrong response for displayedAttributes")
	}

	displayedAttributesString := *displayedAttributesRes
	expecteddisplayedAttrributesString := "*"

	if displayedAttributesString[0] != expecteddisplayedAttrributesString {
		t.Fatalf(
			"Wrong response for displayedAttributes, expected %s, got %s\n",
			displayedAttributesString,
			expecteddisplayedAttrributesString,
		)
	}
}

func TestClientSettings_UpdateDisplayedAttributes(t *testing.T) {
	var indexUID = "TestClientSettings_UpdateDisplayedAttributes"

	_, err := client.Indexes().Create(CreateIndexRequest{
		UID: indexUID,
	})

	if err != nil {
		t.Fatal(err)
	}

	displayedAttributes := []string{"id", "title", "description"}

	updateIDRes, err := client.Settings(indexUID).UpdateDisplayedAttributes(displayedAttributes)

	if err != nil {
		t.Fatal(err)
	}

	_, _ = client.DefaultWaitForPendingUpdate(indexUID, updateIDRes)
}

func TestClientSettings_ResetDisplayedAttributes(t *testing.T) {
	var indexUID = "TestClientSettings_ResetDisplayedAttributes"

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

	_, _ = client.DefaultWaitForPendingUpdate(indexUID, updateIDRes)

	displayedAttributesRes, err := client.Settings(indexUID).GetDisplayedAttributes()

	if err != nil {
		t.Fatal(err)
	}

	if len(*displayedAttributesRes) != 1 {
		t.Fatal("Wrong response for displayedAttributes after reset")
	}

	displayedAttributesString := *displayedAttributesRes
	expecteddisplayedAttrributesString := "*"

	if displayedAttributesString[0] != expecteddisplayedAttrributesString {
		t.Fatalf(
			"Wrong response for displayedAttributes after reset, expected %s, got %s\n",
			displayedAttributesString,
			expecteddisplayedAttrributesString,
		)
	}
}

func TestClientSettings_GetStopWords(t *testing.T) {
	var indexUID = "TestClientSettings_GetStopWords"

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

func TestClientSettings_UpdateStopWords(t *testing.T) {
	var indexUID = "TestClientSettings_UpdateStopWords"

	_, err := client.Indexes().Create(CreateIndexRequest{
		UID: indexUID,
	})

	if err != nil {
		t.Fatal(err)
	}

	stopWords := []string{"a", "the"}

	updateIDRes, err := client.Settings(indexUID).UpdateStopWords(stopWords)

	if err != nil {
		t.Fatal(err)
	}

	_, _ = client.DefaultWaitForPendingUpdate(indexUID, updateIDRes)

}

func TestClientSettings_ResetStopWords(t *testing.T) {
	var indexUID = "TestClientSettings_ResetStopWords"

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

	_, _ = client.DefaultWaitForPendingUpdate(indexUID, updateIDRes)

}

func TestClientSettings_GetSynonyms(t *testing.T) {
	var indexUID = "TestClientSettings_GetSynonyms"

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

func TestClientSettings_UpdateSynonyms(t *testing.T) {
	var indexUID = "TestClientSettings_UpdateSynonyms"

	_, err := client.Indexes().Create(CreateIndexRequest{
		UID: indexUID,
	})

	if err != nil {
		t.Fatal(err)
	}

	synonyms := map[string][]string{
		"car": {"automobile"},
	}

	updateIDRes, err := client.Settings(indexUID).UpdateSynonyms(synonyms)

	if err != nil {
		t.Fatal(err)
	}

	_, _ = client.DefaultWaitForPendingUpdate(indexUID, updateIDRes)
}

func TestClientSettings_ResetSynonyms(t *testing.T) {
	var indexUID = "TestClientSettings_ResetSynonyms"

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

	_, _ = client.DefaultWaitForPendingUpdate(indexUID, updateIDRes)

}

func TestClientSettings_GetAttributesForFaceting(t *testing.T) {
	var indexUID = "TestClientSettings_GetAttributesForFaceting"

	_, err := client.Indexes().Create(CreateIndexRequest{
		UID: indexUID,
	})

	if err != nil {
		t.Fatal(err)
	}

	AttributesForFacetingRes, err := client.Settings(indexUID).GetAttributesForFaceting()

	if err != nil {
		t.Fatal(err)
	}

	if reflect.DeepEqual(*AttributesForFacetingRes, nil) {
		t.Fatal("getAttributesForFaceting: Error getting attributesForFaceting on empty index")
	}
}

func TestClientSettings_UpdateAttributesForFaceting(t *testing.T) {
	var indexUID = "TestClientSettings_UpdateAttributesForFaceting"

	attributesForFaceting := []string{"tag", "title"}

	_, err := client.Indexes().Create(CreateIndexRequest{
		UID: indexUID,
	})

	if err != nil {
		t.Fatal(err)
	}

	updateIDRes, err := client.Settings(indexUID).UpdateAttributesForFaceting(attributesForFaceting)

	if err != nil {
		t.Fatal(err)
	}

	_, _ = client.DefaultWaitForPendingUpdate(indexUID, updateIDRes)
	r, _ := client.Settings(indexUID).GetAttributesForFaceting()
	if !reflect.DeepEqual(*r, attributesForFaceting) {
		fmt.Println(*r)
		t.Fatal("updateAttributesForFaceting: Error getting attributesForFaceting after update")
	}
}

func TestClientSettings_ResetAttributesForFaceting(t *testing.T) {
	var indexUID = "TestClientSettings_ResetAttributesForFaceting"

	attributesForFaceting := []string{"tag", "title"}

	_, err := client.Indexes().Create(CreateIndexRequest{
		UID: indexUID,
	})

	if err != nil {
		t.Fatal(err)
	}

	updateIDRes, err := client.Settings(indexUID).UpdateAttributesForFaceting(attributesForFaceting)

	if err != nil {
		t.Fatal(err)
	}

	_, _ = client.DefaultWaitForPendingUpdate(indexUID, updateIDRes)
	r, _ := client.Settings(indexUID).GetAttributesForFaceting()
	if !reflect.DeepEqual(*r, attributesForFaceting) {
		fmt.Println(*r)
		t.Fatal("resetAttributesForFaceting: Error getting attributesForFaceting after update")
	}

	updateIDRes, err = client.Settings(indexUID).ResetAttributesForFaceting()

	if err != nil {
		t.Fatal(err)
	}

	_, _ = client.DefaultWaitForPendingUpdate(indexUID, updateIDRes)
	r, _ = client.Settings(indexUID).GetAttributesForFaceting()
	if reflect.DeepEqual(*r, nil) {
		t.Fatal("resetAttributesForFaceting: Error getting attributesForFaceting after reset")
	}
}
