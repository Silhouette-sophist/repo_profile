package graph_service

import (
	"context"
	"fmt"
	"strings"
	"sync"

	"github.com/Silhouette-sophist/repo_profile/internal/dal/neo4jdb"
	"github.com/Silhouette-sophist/repo_profile/internal/model"
	"github.com/Silhouette-sophist/repo_profile/internal/service"
	"github.com/Silhouette-sophist/repo_profile/zap_log"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
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
		uniqueIdMap, err := graphService.BatchCreateNodesWithMapping(ctx, declares)
		if err != nil {
			zap_log.CtxError(ctx, "Failed to create nodes", err, zap.Error(err))
			return
		}
		for id, data := range uniqueIdMap {
			fmt.Println(id, data)
		}
		// 补充边关系 uniqueIdMap

	}
}

// GraphService 处理图数据库操作
type GraphService struct {
	connector *neo4jdb.Neo4jConnector
	// 缓存uniqueId到elementId的映射
	idMap     sync.Map
	batchSize int
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

// BatchCreateNodesWithMapping 批量创建节点并返回uniqueId到elementId的映射
func (gs *GraphService) BatchCreateNodesWithMapping(ctx context.Context, nodes []interface{}) (map[string]string, error) {
	if len(nodes) == 0 {
		return map[string]string{}, nil
	}

	// 构建参数和Cypher语句
	params := make(map[string]interface{})
	var mergeClauses []string
	var returnClauses []string

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
			return nil, fmt.Errorf("unsupported node type: %T", node)
		}
		paramKey := fmt.Sprintf("node%d", i)
		nodeVar := fmt.Sprintf("n%d", i)
		params[paramKey] = nodeProps
		// MERGE并设置属性
		mergeClauses = append(mergeClauses,
			fmt.Sprintf(`MERGE (%s:%s {uniqueId: $%s.uniqueId}) 
				SET %s = $%s`, nodeVar, label, paramKey, nodeVar, paramKey))
		// 返回每个节点的uniqueId和elementId
		returnClauses = append(returnClauses,
			fmt.Sprintf("%s.uniqueId AS uniqueId%d, elementId(%s) AS elementId%d",
				nodeVar, i, nodeVar, i))
	}
	// 构建完整的Cypher查询
	cypher := strings.Join(mergeClauses, "\n") + "\nRETURN " + strings.Join(returnClauses, ", ")
	// 执行查询
	session := gs.connector.Driver.NewSession(ctx, neo4j.SessionConfig{})
	defer session.Close(ctx)
	result, err := session.Run(ctx, cypher, params)
	if err != nil {
		return nil, err
	}
	// 解析结果，构建映射
	idMapping := make(map[string]string)
	if record, err := result.Single(ctx); err == nil {
		for i := range nodes {
			uniqueIdKey := fmt.Sprintf("uniqueId%d", i)
			elementIdKey := fmt.Sprintf("elementId%d", i)

			if uniqueId, ok := record.Get(uniqueIdKey); ok {
				if elementId, ok := record.Get(elementIdKey); ok {
					idMapping[uniqueId.(string)] = elementId.(string)
					// 更新缓存
					gs.idMap.Store(uniqueId.(string), elementId.(string))
				}
			}
		}
	}

	return idMapping, nil
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
	if _, err := gs.BatchCreateNodesWithMapping(ctx, nodes); err != nil {
		return fmt.Errorf("failed to create nodes: %w", err)
	}
	// 再创建关系
	if err := gs.BatchCreateRelations(ctx, relations); err != nil {
		return fmt.Errorf("failed to create relations: %w", err)
	}
	return nil
}
