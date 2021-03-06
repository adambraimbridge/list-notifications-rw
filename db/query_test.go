package db

import (
	"regexp"
	"testing"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"gopkg.in/mgo.v2/bson"
)

func init() {
	log.SetLevel(log.DebugLevel)
}

func TestShiftSince(t *testing.T) {
	cacheDelay = 42
	now := time.Now()

	shifted := shiftSince(now)
	assert.Equal(t, now.Add(-42*time.Second), shifted, "Should've been shifted by 42s!")
}

func TestCalculateTillDate(t *testing.T) {
	now := time.Now().UTC()
	cacheDelay = 48

	till := calculateTill(now)
	assert.Equal(t, now.Add(-48*time.Second), till, "Should've been shifted by 48s!")
}

func TestGetMatchNoOffset(t *testing.T) {
	since := time.Now().UTC()

	match := getMatch(0, since)

	data, _ := bson.MarshalJSON(match)
	regex := regexp.MustCompile(`\{"\$match":\{"lastModified":\{"\$gt":\{"\$date":".*"},"\$lt":\{"\$date":".*"}}}}`)
	assert.True(t, regex.MatchString(string(data)), "Query json should match!")
}

func TestGetMatchWithOffset(t *testing.T) {
	since, _ := time.Parse(time.RFC3339Nano, "2016-10-26T16:15:09.46Z")

	match := getMatch(60, since)

	data, _ := bson.MarshalJSON(match)
	regex := regexp.MustCompile(`\{"\$match":\{"lastModified":\{"\$gte":\{"\$date":".*"},"\$lte":\{"\$date":".*"}}}}`)
	assert.True(t, regex.MatchString(string(data)), "Query json should match!")
}

func TestQuery(t *testing.T) {
	since, _ := time.Parse(time.RFC3339Nano, "2016-10-26T16:15:09.46Z")
	maxLimit = 102

	query := generateQuery(50, since)

	regex := regexp.MustCompile(`\[\{"\$match":\{"lastModified":\{"\$gte":\{"\$date":".*"},"\$lte":\{"\$date":".*"}}}},\{"\$sort":\{"lastModified":-1}},\{"\$group":\{"_id":"\$uuid","eventType":\{"\$first":"\$eventType"},"lastModified":\{"\$first":"\$lastModified"},"publishReference":\{"\$first":"\$publishReference"},"title":\{"\$first":"\$title"},"uuid":\{"\$first":"\$uuid"}}},\{"\$sort":\{"lastModified":1,"uuid":1}},\{"\$skip":50},\{"\$limit":103}]`)
	data, _ := bson.MarshalJSON(query)
	assert.True(t, regex.MatchString(string(data)), "Query json should match!")
}

func TestFindNotificationQuery(t *testing.T) {
	query := findByTxID("tid_i-am-a-tid")

	data, _ := bson.MarshalJSON(query)
	assert.Contains(t, string(data), `{"publishReference":"tid_i-am-a-tid"}`, "Query json should match!")
}

func TestFindNotificationQueryByPartialTXID(t *testing.T) {
	query := findByPartialTxID("tid_i-am-a-tid")

	data, _ := bson.MarshalJSON(query)
	assert.Contains(t, string(data), `{"publishReference":{"$regex":"^tid_i-am-a-tid"}}`)
}
