package graph_service

import (
	"context"
	"fmt"
	"strings"

	"github.com/Silhouette-sophist/repo_profile/internal/dal/neo4jdb"
)

type (
	Declaration struct {
		Package   string
		Name      string
		Content   string
		File      *string
		StartLine *int
		EndLine   *int
		UniqueId  string
	}
	AstFunction struct {
		Declaration
		Receiver *string
	}

	AstStruct struct {
		Declaration
	}

	AstVariable struct {
		Declaration
	}

	BaseRelation struct {
		SourceElementId string
		TargetElementId string
		RelationType    string
	}

	RelationType string
)

const (
	INVOKE     RelationType = "INVOKE"
	REFERENCE  RelationType = "REFERENCE"
	ASSOCIATE  RelationType = "ASSOCIATE"
	DEPENDENCE RelationType = "DEPENDENCE"
)

// GraphService 处理图数据库操作
type GraphService struct {
	connector *neo4jdb.Neo4jConnector
}

// NewGraphService 创建新的图服务实例
func NewGraphService(ctx context.Context) (*GraphService, error) {
	connector, err := neo4jdb.NewNeo4jConnector(ctx)
	if err != nil {
		return nil, err
	}
	return &GraphService{connector: connector}, nil
}

// Close 关闭数据库连接
func (gs *GraphService) Close() error {
	return gs.connector.Close()
}

func generateUniqueId(packageName, receiver, name, label string) string {
	if receiver == "" {
		receiver = "_"
	}
	return fmt.Sprintf("%s:%s:%s:%s", packageName, receiver, name, label)
}

// BatchCreateNodes 批量创建节点
func (gs *GraphService) BatchCreateNodes(ctx context.Context, nodes []interface{}) error {
	if len(nodes) == 0 {
		return nil
	}

	// 构建 CYPHER 查询
	var cypherParts []string
	params := make(map[string]interface{})

	for i, node := range nodes {
		var label string
		nodeProps := make(map[string]interface{})

		switch n := node.(type) {
		case *AstFunction:
			label = "AstFunction"
			receiver := ""
			if n.Receiver != nil {
				receiver = *n.Receiver
			}
			nodeProps = map[string]interface{}{
				"package":   n.Package,
				"name":      n.Name,
				"content":   n.Content,
				"file":      n.File,
				"startLine": n.StartLine,
				"endLine":   n.EndLine,
				"uniqueId":  generateUniqueId(n.Package, receiver, n.Name, label),
				"receiver":  n.Receiver,
			}

		case *AstStruct:
			label = "AstStruct"
			nodeProps = map[string]interface{}{
				"package":   n.Package,
				"name":      n.Name,
				"content":   n.Content,
				"file":      n.File,
				"startLine": n.StartLine,
				"endLine":   n.EndLine,
				"uniqueId":  generateUniqueId(n.Package, "", n.Name, label),
			}
		case *AstVariable:
			label = "AstVariable"
			nodeProps = map[string]interface{}{
				"package":   n.Package,
				"name":      n.Name,
				"content":   n.Content,
				"file":      n.File,
				"startLine": n.StartLine,
				"endLine":   n.EndLine,
				"uniqueId":  generateUniqueId(n.Package, "", n.Name, label),
			}
		default:
			return fmt.Errorf("unsupported node type: %T", node)
		}
		paramKey := fmt.Sprintf("node%d", i)
		params[paramKey] = nodeProps
		// 使用 MERGE 避免重复创建
		cypherParts = append(cypherParts,
			fmt.Sprintf(`MERGE (n:%s {uniqueId: $%s.uniqueId}) 
				SET n = $%s`, label, paramKey, paramKey))
	}
	cypher := strings.Join(cypherParts, "\n")
	_, err := gs.connector.ExecuteCypher(ctx, cypher, params)
	return err
}

// BatchCreateRelations 批量创建关系
func (gs *GraphService) BatchCreateRelations(ctx context.Context, relations []BaseRelation) error {
	if len(relations) == 0 {
		return nil
	}

	// 构建 CYPHER 查询
	var cypherParts []string
	params := make(map[string]interface{})

	for i, rel := range relations {
		paramKey := fmt.Sprintf("rel%d", i)
		params[paramKey] = rel

		// 匹配源节点和目标节点，创建关系
		cypherParts = append(cypherParts,
			fmt.Sprintf(`MATCH (source {uniqueId: $%s.SourceElementId}), 
				(target {uniqueId: $%s.TargetElementId}) 
				MERGE (source)-[:%s]->(target)`,
				paramKey, paramKey, rel.RelationType))
	}

	cypher := strings.Join(cypherParts, "\n")

	_, err := gs.connector.ExecuteCypher(ctx, cypher, params)
	return err
}

// BatchImport 批量导入节点和关系
func (gs *GraphService) BatchImport(ctx context.Context, nodes []interface{}, relations []BaseRelation) error {
	// 先创建节点
	if err := gs.BatchCreateNodes(ctx, nodes); err != nil {
		return fmt.Errorf("failed to create nodes: %w", err)
	}
	// 再创建关系
	if err := gs.BatchCreateRelations(ctx, relations); err != nil {
		return fmt.Errorf("failed to create relations: %w", err)
	}
	return nil
}
