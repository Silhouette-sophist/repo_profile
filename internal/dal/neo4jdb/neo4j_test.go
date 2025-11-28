package neo4jdb

import (
	"context"
	"testing"
)

func TestInit(t *testing.T) {
	ctx := context.Background()
	connector, err := NewNeo4jConnector(ctx)
	if err != nil {
		t.Fatal(err)
	}
	first, err := connector.CreatePersonNode("one", 23)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("first node created: %v", first)
	second, err := connector.CreatePersonNode("two", 32)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("second node created: %v", second)
	if err = connector.CreateFriendship("one", "two", 2023); err != nil {
		t.Fatal(err)
	}
	persons, err := connector.GetAllPersons()
	if err != nil {
		t.Fatal(err)
	}
	for _, person := range persons {
		t.Logf("person: %v", person)
	}
}
