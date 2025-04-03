package flamigo

import (
	"fmt"
	"hash/fnv"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func NewRandomID() string {
	return primitive.NewObjectID().Hex()
}

func NewHashId(id string) string {

	h := fnv.New32a()
	h.Write([]byte(id))
	return fmt.Sprintf("%d", h.Sum32())
}

func NewHashSeed(id string) int64 {
	hashId := NewHashId(id)
	hash := 0
	for _, char := range hashId {
		hash = (hash*31 + int(char)) % 1000 // Simple hash function with modulo to keep values manageable
	}
	return int64(hash)
}
