## 仓库画像

### 静态分析
#### 单点分析
基于treesitter完成不同文件的语法单元分析

#### 语法链路分析
基于ast+treesitter

`
在 `golang.org/x/tools/go/packages` 包中，`packages.Package` 结构体包含了代码分析所需的语法树、类型信息等关键数据。其中 `TypesInfo`、`Defs`、`Types` 是与类型检查相关的核心字段，作用如下：


### 1. `packages.Package.TypesInfo`
类型：`*types.Info`  
**作用**：存储当前包的**类型检查结果**，是连接抽象语法树（`ast`）与类型系统（`types`）的核心桥梁。

`types.Info` 结构体（来自标准库 `go/types`）包含多个映射表，记录了语法节点（如表达式、标识符）与类型信息的对应关系，常用字段包括：
- `Types`：`map[ast.Expr]types.TypeAndValue` —— 存储每个表达式（`ast.Expr`）对应的类型和值信息。
- `Defs`：`map[*ast.Ident]types.Object` —— 存储每个标识符（`ast.Ident`）在定义处对应的对象（如变量、函数、类型等）。
- `Uses`：`map[*ast.Ident]types.Object` —— 存储每个标识符在引用处对应的对象（即该标识符指向的定义）。
- `Implicits`：`map[ast.Node]types.Object` —— 存储隐式声明的对象（如匿名函数、`init` 函数等）。

简单说，`TypesInfo` 是类型检查器的“结果集”，所有语法节点的类型信息都通过它查询。


### 2. `packages.Package.Defs`
类型：`map[*ast.Ident]types.Object`  
**作用**：等价于 `TypesInfo.Defs`，是 `TypesInfo` 中 `Defs` 字段的快捷引用，直接存储**标识符定义处**对应的对象。

- 当你在语法树中找到一个标识符（`*ast.Ident`，如变量名、函数名），且该标识符是“定义”（而非“引用”）时，通过 `Defs[ident]` 可获取它对应的 `types.Object`（包含类型、作用域等信息）。
- 例如，对于代码 `var x int`，`x` 是定义处的标识符，`Defs[x]` 会返回一个 `*types.Var` 对象，其 `Type()` 方法返回 `int` 类型。

本质上，`pkg.Defs` 是 `pkg.TypesInfo.Defs` 的别名，方便直接访问。


### 3. `packages.Package.Types`
类型：`map[ast.Expr]types.TypeAndValue`  
**作用**：等价于 `TypesInfo.Types`，是 `TypesInfo` 中 `Types` 字段的快捷引用，存储**表达式**对应的类型和值信息。

- 对于任意表达式（`ast.Expr`，如 `x+1`、`obj.Field`、`funcCall()` 等），`Types[expr]` 会返回一个 `types.TypeAndValue` 结构体，其中：
    - `Type` 字段是表达式的静态类型（`types.Type`）。
    - `Value` 字段是编译期可确定的常量值（若表达式是常量）。

- 例如，对于代码 `a := 1 + 2`，表达式 `1 + 2` 的 `Types` 记录中，`Type` 是 `int`，`Value` 是 `3`。

同样，`pkg.Types` 是 `pkg.TypesInfo.Types` 的别名，方便直接访问表达式的类型。


### 总结关系
- `TypesInfo` 是总容器，包含 `Defs`、`Types` 等所有类型映射信息。
- `Defs` 是 `TypesInfo.Defs` 的快捷方式，专注于“标识符定义”与对象的映射。
- `Types` 是 `TypesInfo.Types` 的快捷方式，专注于“表达式”与类型/值的映射。

使用时，可根据场景直接选择 `pkg.Defs` 或 `pkg.Types`（更简洁），或通过 `pkg.TypesInfo` 访问更多细节（如 `Uses`、`Implicits`）。`

#### 编译链路分支分析
基于ssa分析

### 构图
#### 基于sqlite构图【关系型数据库】

#### 基于neo4j构图
- 函数invoke关系
- 函数refer包变量
- 函数associate类型
- 类型depend类型

### 检索能力
#### 关键词检索

#### 图关系检索

#### 向量检索

### agent调用
#### ls/view/grep/grob基于文本读取
https://github.com/charmbracelet/crush

#### tools加载