package graph_service

import (
	"context"
	"fmt"
	"strings"

	"github.com/Silhouette-sophist/repo_profile/internal/dal/neo4jdb"
	"github.com/Silhouette-sophist/repo_profile/internal/model"
	"github.com/Silhouette-sophist/repo_profile/internal/service"
	"github.com/Silhouette-sophist/repo_profile/zap_log"
	"go.uber.org/zap"
)

func TransferGraph(ctx context.Context, repoPath string) {
	repoModules, err := service.ParseRepo(repoPath)
	if err != nil {
		zap_log.CtxError(ctx, "Failed to parse repo", err, zap.Error(err))
		return
	}
	graphService, err := NewGraphService(ctx)
	if err != nil {
		zap_log.CtxError(ctx, "Failed to create graph service", err, zap.Error(err))
		return
	}
	for _, module := range repoModules {
		zap_log.CtxInfo(ctx, "parse module", zap.String("module", module.Path))
		declares := make([]interface{}, 0)
		for _, infos := range module.PkgFuncMap {
			for _, info := range infos {
				receiver := ""
				if info.Receiver != nil {
					receiver = info.Receiver.Name
				}
				function := &model.AstFunction{
					Declaration: model.Declaration{
						Package:   info.Pkg,
						Name:      info.Name,
						File:      &info.RFilePath,
						StartLine: &info.StartPosition.Line,
						EndLine:   &info.EndPosition.Line,
						Content:   info.Content,
						UniqueId:  generateUniqueId(info.Pkg, receiver, info.Name, "AstFunction"),
					},
				}
				if info.Receiver != nil {
					function.Receiver = &info.Receiver.BaseType
				}
				declares = append(declares, function)
			}
		}
		for _, infos := range module.PkgStructMap {
			for _, info := range infos {
				declares = append(declares, &model.AstStruct{
					Declaration: model.Declaration{
						Package:   info.Pkg,
						Name:      info.Name,
						File:      &info.RFilePath,
						StartLine: &info.StartPosition.Line,
						EndLine:   &info.EndPosition.Line,
						Content:   info.Content,
						UniqueId:  generateUniqueId(info.Pkg, "", info.Name, "AstStruct"),
					},
				})
			}
		}
		for _, infos := range module.PkgVarMap {
			for _, info := range infos {
				declares = append(declares, &model.AstVariable{
					Declaration: model.Declaration{
						Package:  info.Pkg,
						Name:     info.Name,
						File:     &info.RFilePath,
						Content:  info.Content,
						UniqueId: generateUniqueId(info.Pkg, "", info.Name, "AstVariable"),
					},
				})
			}
		}
		if err := graphService.BatchCreateNodes(ctx, declares); err != nil {
			zap_log.CtxError(ctx, "Failed to create nodes", err, zap.Error(err))
			return
		}
	}
}

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
		case *model.AstFunction:
			label = "AstFunction"
			nodeProps = map[string]interface{}{
				"package":   n.Package,
				"name":      n.Name,
				"content":   n.Content,
				"file":      n.File,
				"startLine": n.StartLine,
				"endLine":   n.EndLine,
				"uniqueId":  n.UniqueId,
				"receiver":  n.Receiver,
			}
		case *model.AstStruct:
			label = "AstStruct"
			nodeProps = map[string]interface{}{
				"package":   n.Package,
				"name":      n.Name,
				"content":   n.Content,
				"file":      n.File,
				"startLine": n.StartLine,
				"endLine":   n.EndLine,
				"uniqueId":  n.UniqueId,
			}
		case *model.AstVariable:
			label = "AstVariable"
			nodeProps = map[string]interface{}{
				"package":   n.Package,
				"name":      n.Name,
				"content":   n.Content,
				"file":      n.File,
				"startLine": n.StartLine,
				"endLine":   n.EndLine,
				"uniqueId":  n.UniqueId,
			}
		default:
			return fmt.Errorf("unsupported node type: %T", node)
		}
		paramKey := fmt.Sprintf("node%d", i)
		nodeVar := fmt.Sprintf("n%d", i) // 使用唯一的节点变量名
		params[paramKey] = nodeProps
		// 使用唯一的变量名避免冲突
		cypherParts = append(cypherParts,
			fmt.Sprintf(`MERGE (%s:%s {uniqueId: $%s.uniqueId}) 
                SET %s = $%s`, nodeVar, label, paramKey, nodeVar, paramKey))
	}
	cypher := strings.Join(cypherParts, "\n")
	_, err := gs.connector.ExecuteCypher(ctx, cypher, params)
	return err
}

// BatchCreateRelations 批量创建关系
func (gs *GraphService) BatchCreateRelations(ctx context.Context, relations []model.BaseRelation) error {
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
func (gs *GraphService) BatchImport(ctx context.Context, nodes []interface{}, relations []model.BaseRelation) error {
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
