package neo4jdb

import (
	"context"
	"fmt"

	"github.com/Silhouette-sophist/repo_profile/zap_log"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"go.uber.org/zap"
)

const (
	uri      = "bolt://localhost:7687"
	userName = "neo4j"
	password = "chen150928"
)

// Neo4jConnector 管理 Neo4j 连接
type Neo4jConnector struct {
	Driver neo4j.DriverWithContext
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
		Driver: driver,
		ctx:    ctx,
	}, nil
}

// Close 关闭连接
func (nc *Neo4jConnector) Close() error {
	return nc.Driver.Close(nc.ctx)
}

// ExecuteCypher 执行cypher语句
func (nc *Neo4jConnector) ExecuteCypher(ctx context.Context, cypher string, params map[string]any) ([]map[string]any, error) {
	session := nc.Driver.NewSession(ctx, neo4j.SessionConfig{})
	defer session.Close(ctx)
	run, err := session.Run(ctx, cypher, params)
	if err != nil {
		zap_log.CtxError(ctx, "Failed to run query", err, zap.Error(err))
		return nil, err
	}
	collectRecords, err := run.Collect(ctx)
	if err != nil {
		zap_log.CtxError(ctx, "Failed to run query", err, zap.Error(err))
		return nil, err
	}
	data := make([]map[string]any, len(collectRecords))
	for _, record := range collectRecords {
		data = append(data, record.AsMap())
	}
	return data, nil
}
