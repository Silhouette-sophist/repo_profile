package graph_service

import (
	"context"

	"github.com/Silhouette-sophist/repo_profile/internal/dal/neo4jdb"
	"github.com/Silhouette-sophist/repo_profile/zap_log"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"go.uber.org/zap"
)

// ExecuteCypher 执行cypher语句
func ExecuteCypher(ctx context.Context, cypher string, params map[string]any) ([]map[string]any, error) {
	connector, err := neo4jdb.NewNeo4jConnector(ctx)
	if err != nil {
		zap_log.CtxError(ctx, "Failed to connect to Neo4j", err, zap.Error(err))
		return nil, err
	}
	defer connector.Close()
	session := connector.Driver.NewSession(ctx, neo4j.SessionConfig{})
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
