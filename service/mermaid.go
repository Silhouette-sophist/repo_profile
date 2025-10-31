package service

import (
	"fmt"
	"strings"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

// DataItem 完全使用Neo4j官方类型，存储一条完整路径（n→...→m）
type DataItem struct {
	N neo4j.Node           `json:"n"`
	R []neo4j.Relationship `json:"r"`
	M neo4j.Node           `json:"m"`
}

// MermaidConverter Mermaid转换器
type MermaidConverter struct {
	nodes         map[int64]*neo4j.Node // 节点ID到节点的映射
	relationships []*neo4j.Relationship // 所有关系
	nodeNames     map[int64]string      // 节点ID到显示名称的映射
}

func NewMermaidConverter() *MermaidConverter {
	return &MermaidConverter{
		nodes:         make(map[int64]*neo4j.Node),
		relationships: make([]*neo4j.Relationship, 0),
		nodeNames:     make(map[int64]string),
	}
}

// 处理数据项
func (c *MermaidConverter) ProcessDataItems(dataItems []DataItem) {
	for _, item := range dataItems {
		// 处理起始节点 N
		c.addNode(&item.N)

		// 处理路径中的所有关系
		for i := range item.R {
			c.addRelationship(&item.R[i])
		}

		// 处理结束节点 M
		c.addNode(&item.M)
	}
}

// 添加节点到集合
func (c *MermaidConverter) addNode(node *neo4j.Node) {
	if _, exists := c.nodes[node.Id]; !exists {
		c.nodes[node.Id] = node
		c.nodeNames[node.Id] = c.generateNodeName(node)
	}
}

// 添加关系到集合
func (c *MermaidConverter) addRelationship(rel *neo4j.Relationship) {
	// 检查是否已存在相同的关系
	for _, existingRel := range c.relationships {
		if existingRel.StartId == rel.StartId && existingRel.EndId == rel.EndId {
			return
		}
	}
	c.relationships = append(c.relationships, rel)

	// 确保关系涉及的两个节点都在节点集合中
	// 注意：这里我们假设节点已经在ProcessDataItems中被添加
}

// 生成节点显示名称
func (c *MermaidConverter) generateNodeName(node *neo4j.Node) string {
	// 优先使用Name属性
	if name, ok := node.Props["Name"].(string); ok && name != "" {
		return name
	}

	// 其次使用UniqueId
	if uniqueId, ok := node.Props["UniqueId"].(string); ok && uniqueId != "" {
		// 从UniqueId中提取简化的名称
		parts := strings.Split(uniqueId, ":")
		if len(parts) > 0 {
			return parts[len(parts)-1]
		}
		return uniqueId
	}

	// 最后使用标签
	if len(node.Labels) > 0 {
		return node.Labels[0]
	}

	// 如果都没有，使用ID
	return fmt.Sprintf("Node_%d", node.Id)
}

// 获取节点类型用于样式
func (c *MermaidConverter) getNodeType(node *neo4j.Node) string {
	if nodeType, ok := node.Props["Type"].(string); ok {
		return nodeType
	}
	return "UNKNOWN"
}

// 生成Mermaid代码
func (c *MermaidConverter) GenerateMermaid() string {
	var builder strings.Builder

	// Mermaid头部
	builder.WriteString("graph TD\n")

	// 生成所有节点定义
	for id, node := range c.nodes {
		nodeName := c.nodeNames[id]
		nodeType := c.getNodeType(node)
		styleClass := c.getStyleClass(nodeType)

		// 添加节点信息
		displayName := nodeName
		if file, ok := node.Props["File"].(string); ok && file != "" {
			displayName = fmt.Sprintf("%s<br/>%s", nodeName, c.getShortFilePath(file))
		}

		builder.WriteString(fmt.Sprintf("    %d[\"%s\"]%s\n", id, displayName, styleClass))
	}

	builder.WriteString("\n")

	// 生成所有关系
	for _, rel := range c.relationships {

		// 添加关系
		builder.WriteString(fmt.Sprintf("    %d --> %d\n", rel.StartId, rel.EndId))
	}

	// 添加样式定义
	builder.WriteString(c.generateStyles())

	return builder.String()
}

// 获取短文件路径
func (c *MermaidConverter) getShortFilePath(fullPath string) string {
	parts := strings.Split(fullPath, "/")
	if len(parts) > 2 {
		return parts[len(parts)-2] + "/" + parts[len(parts)-1]
	}
	return fullPath
}

// 获取样式类
func (c *MermaidConverter) getStyleClass(nodeType string) string {
	switch nodeType {
	case "HTTP_API":
		return ":::http-api"
	case "FUNC":
		return ":::function"
	case "STRUCT":
		return ":::struct"
	case "INTERFACE":
		return ":::interface"
	default:
		return ":::default"
	}
}

// 生成样式定义
func (c *MermaidConverter) generateStyles() string {
	return `
    classDef http-api fill:#e1f5fe,stroke:#01579b,stroke-width:2px
    classDef function fill:#f3e5f5,stroke:#4a148c,stroke-width:1px
    classDef struct fill:#e8f5e8,stroke:#1b5e20,stroke-width:1px
    classDef interface fill:#fff3e0,stroke:#e65100,stroke-width:1px
    classDef default fill:#f5f5f5,stroke:#616161,stroke-width:1px
`
}
