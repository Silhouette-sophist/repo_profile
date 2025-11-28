package neo4j

import (
	"context"
	"fmt"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

var Neo4jDriver neo4j.DriverWithContext

func init() {
	uri := "bolt://localhost:7687"
	userName := "neo4j"
	password := "chen150928"
	ctx := context.Background()
	// 1.创建链接
	withContext, err := neo4j.NewDriverWithContext(uri, neo4j.BasicAuth(userName, password, ""))
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	// 2.测试连接
	if err := withContext.VerifyConnectivity(ctx); err != nil {
		fmt.Println(err)
		panic(err)
	}
	fmt.Println("Connectivity OK")
	// 3.持有连接
	Neo4jDriver = withContext
}
