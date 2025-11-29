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
		if err := HandleSingleModule(ctx, module, graphService); err != nil {
			return
		}
	}
}

func HandleSingleModule(ctx context.Context, module *service.ModuleInfo, graphService *GraphService) error {
	zap_log.CtxInfo(ctx, "parse module", zap.String("module", module.Path))
	// 1.先存储点关系
	uniqueIdMap, err := SaveModuleNodes(ctx, module, graphService)
	if err != nil {
		return err
	}
	// 2.补充边关系 uniqueIdMap
	if err = SaveModuleEdges(ctx, module, uniqueIdMap, graphService); err != nil {
		return err
	}
	return nil
}

func SaveModuleEdges(ctx context.Context, module *service.ModuleInfo, uniqueIdMap map[string]string, graphService *GraphService) error {
	relations := make([]model.BaseRelation, 0)
	for _, infos := range module.PkgFuncMap {
		for _, info := range infos {
			callerReceiver := ""
			if info.Receiver != nil {
				callerReceiver = info.Receiver.Name
			}
			callerUniqueId := generateUniqueId(info.Pkg, callerReceiver, info.Name, "AstFunction")
			callerElementId, ok := uniqueIdMap[callerUniqueId]
			if !ok {
				continue
			}
			if len(info.RelatedCallee) > 0 {
				calleeSet := make(map[string]struct{})
				for _, indices := range info.RelatedCallee {
					for _, index := range indices {
						calleeReceiver := ""
						if index.Receiver != nil {
							calleeReceiver = *index.Receiver
						}
						calleeUniqueId := generateUniqueId(index.Pkg, calleeReceiver, index.Name, "AstFunction")
						if _, ok := calleeSet[calleeUniqueId]; ok {
							continue
						}
						if calleeElementId, ok := uniqueIdMap[calleeUniqueId]; ok {
							relations = append(relations, model.BaseRelation{
								SourceElementId: callerElementId,
								TargetElementId: calleeElementId,
								RelationType:    model.INVOKE.String(),
							})
							fmt.Println("invoke caller", callerElementId, "callee:", calleeElementId)
						}
						calleeSet[calleeUniqueId] = struct{}{}
					}
				}
			}
			for _, indices := range info.RelatedPkgStruct {
				associateSet := make(map[string]struct{})
				for _, index := range indices {
					associateId := generateUniqueId(index.Pkg, "", index.Name, "AstStruct")
					if _, ok := associateSet[associateId]; ok {
						continue
					}
					if s, ok := uniqueIdMap[associateId]; ok {
						relations = append(relations, model.BaseRelation{
							SourceElementId: callerElementId,
							TargetElementId: s,
							RelationType:    model.ASSOCIATE.String(),
						})
						fmt.Println("associate caller", callerElementId, "callee:", s)
					}
					associateSet[associateId] = struct{}{}
				}
			}
			for _, indices := range info.RelatedPkgVar {
				referenceSet := make(map[string]struct{})
				for _, index := range indices {
					referenceId := generateUniqueId(index.Pkg, "", index.Name, "AstVariable")
					if _, ok := referenceSet[referenceId]; ok {
						continue
					}
					if s, ok := uniqueIdMap[referenceId]; ok {
						relations = append(relations, model.BaseRelation{
							SourceElementId: callerElementId,
							TargetElementId: s,
							RelationType:    model.REFERENCE.String(),
						})
						fmt.Println("reference caller", callerElementId, "callee:", s)
					}
					referenceSet[referenceId] = struct{}{}
				}
			}
		}
	}
	counts, err := graphService.BatchCreateRelations(ctx, relations)
	if err != nil {
		zap_log.CtxError(ctx, "Failed to create relations", err, zap.Error(err))
		return err
	}
	fmt.Println(counts)
	return nil
}

func SaveModuleNodes(ctx context.Context, module *service.ModuleInfo, graphService *GraphService) (map[string]string, error) {
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
		return nil, err
	}
	fmt.Printf("uniqueIdMap: %v\n", uniqueIdMap)
	return uniqueIdMap, nil
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
	fmt.Printf("BatchCreateNodesWithMapping len %d\n", len(nodes))
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

// BatchCreateRelations 使用 UNWIND 批量创建关系（更高性能）
func (gs *GraphService) BatchCreateRelations(ctx context.Context, relations []model.BaseRelation) (int, error) {
	fmt.Printf("BatchCreateRelations len %d\n", len(relations))
	if len(relations) == 0 {
		return 0, nil
	}

	session := gs.connector.Driver.NewSession(ctx, neo4j.SessionConfig{})
	defer session.Close(ctx)

	// 准备批量参数
	params := make([]map[string]interface{}, len(relations))
	for i, rel := range relations {
		params[i] = map[string]interface{}{
			"sourceId": rel.SourceElementId,
			"targetId": rel.TargetElementId,
			"relType":  rel.RelationType,
		}
		fmt.Printf("BatchCreateRelations rel %+v\n", rel)
	}

	// 按关系类型分组处理
	relGroups := make(map[string][]map[string]interface{})
	for _, param := range params {
		relType := param["relType"].(string)
		if _, ok := relGroups[relType]; !ok {
			relGroups[relType] = []map[string]interface{}{}
		}
		relGroups[relType] = append(relGroups[relType], param)
	}

	totalCreated := 0

	// 对每种关系类型执行批量创建
	for relType, batchParams := range relGroups {
		// 使用 UNWIND 进行批量操作（修正：使用 elementId 匹配节点）
		cypher := fmt.Sprintf(`
    UNWIND $batch AS rel
    MATCH (source) WHERE elementId(source) = rel.sourceId
    MATCH (target) WHERE elementId(target) = rel.targetId
    MERGE (source)-[r:%s]->(target)
    RETURN COUNT(r) as createdCount
    `, relType)

		// 执行查询
		result, err := session.Run(ctx, cypher, map[string]interface{}{
			"batch": batchParams,
		})
		if err != nil {
			return totalCreated, err
		}

		// 获取创建的关系数量
		if record, err := result.Single(ctx); err == nil {
			if count, ok := record.Get("createdCount"); ok {
				totalCreated += int(count.(int64))
			}
		}
	}
	fmt.Printf("Successfully created %d relations out of %d attempted\n", totalCreated, len(relations))
	return 0, nil
}

// BatchImport 批量导入节点和关系
func (gs *GraphService) BatchImport(ctx context.Context, nodes []interface{}, relations []model.BaseRelation) error {
	// 先创建节点
	if _, err := gs.BatchCreateNodesWithMapping(ctx, nodes); err != nil {
		return fmt.Errorf("failed to create nodes: %w", err)
	}
	// 再创建关系
	counts, err := gs.BatchCreateRelations(ctx, relations)
	if err != nil {
		return fmt.Errorf("failed to create relations: %w", err)
	}
	fmt.Printf("Importing %d relations %d\n", len(relations), counts)
	return nil
}
