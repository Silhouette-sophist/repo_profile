package neo4j

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

const (
	uri      = "bolt://localhost:7687"
	userName = "neo4j"
	password = "chen150928"
)

// Neo4jConnector 管理 Neo4j 连接
type Neo4jConnector struct {
	driver neo4j.DriverWithContext
	ctx    context.Context
}

// NewNeo4jConnector 创建新的 Neo4j 连接器
func NewNeo4jConnector(ctx context.Context) (*Neo4jConnector, error) {
	// 创建驱动
	driver, err := neo4j.NewDriverWithContext(uri, neo4j.BasicAuth(userName, password, ""))
	if err != nil {
		return nil, fmt.Errorf("创建驱动失败: %w", err)
	}
	// 验证连接
	err = driver.VerifyConnectivity(ctx)
	if err != nil {
		return nil, fmt.Errorf("连接验证失败: %w", err)
	}
	return &Neo4jConnector{
		driver: driver,
		ctx:    ctx,
	}, nil
}

// Close 关闭连接
func (nc *Neo4jConnector) Close() error {
	return nc.driver.Close(nc.ctx)
}

// CreatePersonNode 创建人物节点
func (nc *Neo4jConnector) CreatePersonNode(name string, age int) (neo4j.Node, error) {
	session := nc.driver.NewSession(nc.ctx, neo4j.SessionConfig{})
	defer session.Close(nc.ctx)

	result, err := session.ExecuteWrite(nc.ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		query := `
			CREATE (p:Person {name: $name, age: $age, createdAt: $createdAt})
			RETURN p
		`
		params := map[string]any{
			"name":      name,
			"age":       age,
			"createdAt": time.Now().UTC(),
		}
		result, err := tx.Run(nc.ctx, query, params)
		if err != nil {
			return nil, err
		}
		record, err := result.Single(nc.ctx)
		if err != nil {
			return nil, err
		}
		saveRecord, success := record.Get("p")
		if !success {
			return nil, errors.New("not found")
		}
		return saveRecord, nil
	})
	if err != nil {
		return neo4j.Node{}, fmt.Errorf("创建节点失败: %w", err)
	}
	node, ok := result.(neo4j.Node)
	if !ok {
		return neo4j.Node{}, fmt.Errorf("结果不是节点类型")
	}
	return node, nil
}

// CreateFriendship 创建朋友关系
func (nc *Neo4jConnector) CreateFriendship(person1Name, person2Name string, since int) error {
	session := nc.driver.NewSession(nc.ctx, neo4j.SessionConfig{})
	defer session.Close(nc.ctx)

	if _, err := session.ExecuteWrite(nc.ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		query := `
			MATCH (p1:Person {name: $person1}), (p2:Person {name: $person2})
			CREATE (p1)-[:FRIEND_OF {since: $since}]->(p2)
			RETURN count(*)
		`
		params := map[string]any{
			"person1": person1Name,
			"person2": person2Name,
			"since":   since,
		}
		result, err := tx.Run(nc.ctx, query, params)
		if err != nil {
			return nil, err
		}
		return result.Consume(nc.ctx)
	}); err != nil {
		return fmt.Errorf("创建关系失败: %w", err)
	}
	return nil
}

// GetPersonByName 根据名称查询人物
func (nc *Neo4jConnector) GetPersonByName(name string) (map[string]any, error) {
	session := nc.driver.NewSession(nc.ctx, neo4j.SessionConfig{})
	defer session.Close(nc.ctx)

	result, err := session.ExecuteRead(nc.ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		query := `
			MATCH (p:Person {name: $name})
			RETURN p.name AS name, p.age AS age, p.createdAt AS createdAt
		`

		params := map[string]any{"name": name}

		result, err := tx.Run(nc.ctx, query, params)
		if err != nil {
			return nil, err
		}

		record, err := result.Single(nc.ctx)
		if err != nil {
			return nil, err
		}

		name, _ := record.Get("name")
		age, _ := record.Get("age")
		createdAt, _ := record.Get("createdAt")
		return map[string]any{
			"name":      name,
			"age":       age,
			"createdAt": createdAt,
		}, nil
	})

	if err != nil {
		return nil, fmt.Errorf("查询人物失败: %w", err)
	}

	person, ok := result.(map[string]any)
	if !ok {
		return nil, fmt.Errorf("结果格式错误")
	}

	return person, nil
}

// GetFriends 获取某人的所有朋友
func (nc *Neo4jConnector) GetFriends(name string) ([]map[string]any, error) {
	session := nc.driver.NewSession(nc.ctx, neo4j.SessionConfig{})
	defer session.Close(nc.ctx)

	result, err := session.ExecuteRead(nc.ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		query := `
			MATCH (p:Person {name: $name})-[:FRIEND_OF]->(friend:Person)
			RETURN friend.name AS name, friend.age AS age
		`
		params := map[string]any{"name": name}
		result, err := tx.Run(nc.ctx, query, params)
		if err != nil {
			return nil, err
		}
		var friends []map[string]any
		for result.Next(nc.ctx) {
			record := result.Record()
			name, _ := record.Get("name")
			age, _ := record.Get("age")
			friends = append(friends, map[string]any{
				"name": name,
				"age":  age,
			})
		}
		return friends, nil
	})
	if err != nil {
		return nil, fmt.Errorf("查询朋友失败: %w", err)
	}
	friends, ok := result.([]map[string]any)
	if !ok {
		return nil, fmt.Errorf("结果格式错误")
	}
	return friends, nil
}

// GetAllPersons 获取所有人的列表
func (nc *Neo4jConnector) GetAllPersons() ([]map[string]any, error) {
	session := nc.driver.NewSession(nc.ctx, neo4j.SessionConfig{})
	defer session.Close(nc.ctx)

	result, err := session.ExecuteRead(nc.ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		query := `
			MATCH (p:Person)
			RETURN p.name AS name, p.age AS age
			ORDER BY p.name
		`
		result, err := tx.Run(nc.ctx, query, nil)
		if err != nil {
			return nil, err
		}

		var persons []map[string]any
		for result.Next(nc.ctx) {
			record := result.Record()
			name, _ := record.Get("name")
			age, _ := record.Get("age")
			persons = append(persons, map[string]any{
				"name": name,
				"age":  age,
			})
		}
		return persons, nil
	})
	if err != nil {
		return nil, fmt.Errorf("查询所有人失败: %w", err)
	}
	persons, ok := result.([]map[string]any)
	if !ok {
		return nil, fmt.Errorf("结果格式错误")
	}
	return persons, nil
}
