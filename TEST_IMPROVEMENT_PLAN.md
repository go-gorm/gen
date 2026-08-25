# gorm/gen 测试架构改造计划

> 分支：`test/improve-coverage`
> worktree：`/Users/bytedance/workspace/gen-test-coverage`
> 基线：`4be99c78cfef89247147a7a1725b88341ccbceba`
> 状态：**已在当前 worktree 实现；本地验收通过，待 GitHub Actions（含 MySQL 5.7）确认**

## 1. 结论与目标

这次工作不再按“哪个包覆盖率低就补哪个包”的方式推进，而是建立四类可区分的行为保证：

1. **root 单元契约**：安全校验、生成配置、manifest、并发池等 gen 自有逻辑。
2. **生成契约**：同一输入必须生成完整、可编译、与 golden 一致的代码。
3. **执行契约**：生成出来的 query 和 DIY SQL 必须在真实数据库上得到正确结果并传播错误。
4. **工具与工程契约**：gentool 参数/配置可测试，四个 Go module 和 CI 的验证范围准确可见。

覆盖率只用于发现盲区和观察增量，不作为本轮门禁。每个新增测试必须说明它保护的行为、可能捕获的回归，以及它不证明什么。

### 本轮成功标准

- 安全策略由硬编码的行为矩阵保护，不通过遍历当前实现来“自证正确”。
- 增量生成和 MergeQuery 至少有一次跨两轮真实文件生成的端到端回归。
- DB-independent generator fixture 能生成、做目录级 golden 比较，并由 Go 编译器实际编译。
- dialect-neutral DIY DSL 覆盖 `if/else`、`for`、`where`、`set`、`trim`，且 SQL 在 SQLite 中真实执行并断言结果。
- `DO` 的关键查询、聚合、批处理、写入、关系和事务路径有 SQLite 运行时契约。
- MySQL schema introspection/golden 与 MySQL 专属 SQL 继续由 MySQL job 负责，不冒充 SQLite 可验证的能力。
- root、tests、gentool、examples 四个 module 都有明确且正确的验证命令。

## 2. 已核验的当前状态

### 2.1 Go module 与实际职责

| Module | 当前最低版本/工具链 | 当前职责 | 当前验证 |
|---|---|---|---|
| root `gorm.io/gen` | Go 1.22 | 生成器、运行时 DO、field、parser/model | `go test -race ./...` |
| `tests` | Go 1.22，toolchain 1.24.4，replace root | MySQL schema 生成、golden、少量 CRUD | MySQL 5.7/latest job |
| `tools/gentool` | Go 1.25，无 replace | 独立 CLI，使用已发布的 gen | 只有 build |
| `examples` | Go 1.22，toolchain 1.24.4，replace root | 示例和 checked-in generated code | 只有 build |

root 的 `go test ./...` 不会进入三个嵌套 module。所有计划和验收命令必须按 module 分开。

### 2.2 root 覆盖率快照

以下数据由基线提交上的 `go test -coverprofile` 复核：

| 包 | 覆盖率 | 判断 |
|---|---:|---|
| `gorm.io/gen` | 21.8% | DB 执行路径和若干公开输出边界缺失 |
| `field` | 33.5% | 大量同构薄包装为 0%，不应逐函数追数 |
| `helper` | 0.0% | 有明确纯逻辑缺口 |
| `internal/diagnostic` | 58.0% | 已有错误定位测试，公开 JSON 输出尚未覆盖 |
| `internal/generate` | 31.4% | DSL 构建已有部分测试，仍需生成/执行契约兜底 |
| `internal/model` | 14.3% | nullable/default/tag/name 等纯逻辑缺口明显 |
| `internal/parser` | 40.2% | 有直接测试，不是“完全无测试” |
| `internal/utils/pools` | 0.0% | 生成器并发路径直接依赖，值得补测 |

### 2.3 已确认的测试架构问题

- `tests/generate_test.go` 在包初始化阶段无条件把 `GORM_DIALECT` 改成 `mysql`，外部传入 SQLite 无效。
- `tests/tables.sql` 是 MySQL DDL，SQLite 无法解析；`RunMigrations` 只记录错误后继续，存在假通过风险。
- `tests` 根包在 `init()` 中连库、迁移和生成代码，任何该包测试都会承担这些隐式副作用。
- `.gen/` 被 gitignore，不能作为其他测试的编译期依赖；`.expect/` 才是 checked-in fixture。
- 当前 DIY grammar fixture 包含故意不完整、丢弃 error、缺列或 MySQL 专属的 SQL，不能整体拿来做运行时测试。
- 从仓库根执行 `go test ./tests/...` 会跨 module 失败，必须使用 `go -C tests test ...` 或先进入目录。
- tests module 默认覆盖率不会计入 root 的 `do.go`；集成覆盖必须显式列出 root package pattern，避免把 `gorm.io/gen/tests/fixture/query` 误算进 root 覆盖率。
- `.golangci.yml` 设置了 `run.tests: false`，因此 lint 命令不会检查新增测试文件。

### 2.4 明确不是本轮缺口的项目

- `internal/utils/common.go` 只有 package 声明，没有可测试逻辑。
- `internal/template` 主要保存生成代码的字符串模板；应由 golden、编译和执行契约验证，不为其增加直接覆盖测试。
- `examples` 是可编译示例集合，没有 `*_test.go` 不等于没有验证；保持 compile smoke 即可。
- `field` 中大量结构相同的操作符包装不逐个补测，只选择有状态、类型转换或历史风险的代表路径。

## 3. 目标测试架构

| 层 | 建议位置 | DB | 保护的契约 | 明确不证明 |
|---|---|---|---|---|
| L0 快速单元 | root 各 package | 无 | 安全策略、纯函数、manifest、并发和错误序列化 | SQL 能在数据库执行 |
| L1 DB-independent generation | `tests/generation/` | 无 | fixture 能生成、目录一致、generated code 可编译 | 方言 schema introspection |
| L2 SQLite runtime | `tests/runtime/` | SQLite 临时文件 | generated API、DIY DSL、DO 的真实执行语义 | MySQL 专属语法 |
| L3 MySQL contract | `tests` 根包 | MySQL 5.7/latest | schema introspection、现有 golden、MySQL 专属 SQL | SQLite/其他方言兼容性 |
| L4 CLI | `tools/gentool` | 纯单元 + SQLite smoke | 参数/YAML/方言选择/完整 CLI 编排 | root 当前源码行为；该 module 无 replace |
| L5 CI/coverage | GitHub Actions | 分层 | 所有 module 的状态和覆盖范围准确展示 | 单一覆盖率代表整体质量 |

### 3.1 建议新增目录

```text
tests/
├── internal/golden/              # 目录比较器及其单测
├── internal/runtimefixture/      # fixture 的唯一生成编排
├── fixture/
│   ├── model/                    # 独立模型与可执行 DIY interface
│   └── query/                    # checked-in、可 import 的 generated fixture
├── cmd/generate-runtime-fixture/ # 可重复执行的 fixture 生成命令
├── generation/                   # 无 DB 的生成 + golden contract
└── runtime/                      # SQLite 实际执行测试
```

保留现有 `tests/.expect/dal_*` 和 tests 根包，不在第一轮大规模搬目录。现有根包继续是明确的 MySQL contract；SQLite job 只运行 `generation`、`runtime` 等隔离子包，不加载根包的全局状态。

## 4. 工作包设计

### A. root 高价值单元测试

#### A1. `sec_check.go` 安全契约

新增同包 `sec_check_test.go`，硬编码预期策略，不遍历 `banClauses` 生成测试数据。

必须覆盖：

- 默认拒绝的顶层 clause：`VALUES/SELECT/FROM/WHERE/GROUP BY/ORDER BY/LIMIT/UPDATE/SET/DELETE`。
- 未识别 expression、空列表、多个 condition 遇到首个错误即停止。
- 自定义 `ClauseChecker` 的 handled、fallback、custom error 三条路径。
- `Insert`：Raw table、空 modifier、`IGNORE`、priority 与 `IGNORE` 组合、非法 token、大小写和空白。
- `Locking`：合法 Strength/Options、Raw table、非法值。
- `OnConflict`：普通 assignment 与 `clause.Expr`/`*clause.Expr`。
- hints 和 dbresolver 的允许路径。

以下不是直接写 characterization test 固化现状，而是实现前必须确认的安全决策：

1. 默认 checker 是否应继续放行任意未列入 ban list 的 `clause.Interface`。
2. `OnConflict.OnConstraint`、`Columns.Raw`、`Where`、`TargetWhere` 是否需要额外拒绝或递归检查。
3. 单独的 `LOW_PRIORITY/DELAYED/HIGH_PRIORITY` 是否属于允许 modifier。

若决定收紧策略，安全测试和对应源码修复放在同一 PR；不提交预期失败的测试，也不静默改变公开兼容性。

断言以 error 是否存在、error 类型/稳定哨兵为主，避免绑定完整错误文案。

#### A2. `helper` 字符串与 JOIN builder

新增同包测试，覆盖：

- `IfClause` 的 nil/空、真假混合、多个空结果和既有空格契约。
- `WhereClause`/`SetClause` 的 nil、空字符串、AND/OR/XOR/逗号、大小写和首尾清理。
- `trimLeft`/`trimRight`/`trimAll` 的边界，保证只移除语法 token，不误删普通单词。
- `JoinWhereBuilder`/`JoinSetBuilder`/`JoinTrimAllBuilder` 的空与非空输入。
- `JoinTblExpr.Build` 的 nil builder、JOIN type、table expression、ON、USING 多列分支。
- `CheckObject` 的 StructName、Field.Name、Field.Type 和合法对象。

不设“helper 必须 80%”目标；以公开行为和主要分支完整为验收。

#### A3. `internal/utils/pools` 并发契约

覆盖：

- 正容量 pool 的 `Size/Num/Wait/Done`。
- 达到容量后下一次 `Wait` 阻塞，归还 token 后继续。
- `WaitAll` 和 `AsyncWaitAll` 只在全部 token 归还后完成。
- 负 size 的当前 no-op 语义。
- 所有用例在 `go test -race` 下运行，不使用固定长 sleep；用 channel 和短超时只作为失败保护。

`size == 0` 当前会使 `Wait` 阻塞，不直接调用造成挂死；先确认它应表示“禁止并发”还是“无限制”，再决定是否修源码和增加契约。

#### A4. diagnostic JSON 公共输出

给 `WriteDiagnosticJSON` 增加：

- nil error 输出 `null` 加换行。
- `diagnostic.Error` 输出结构化字段。
- 普通 error 输出稳定 fallback JSON。
- writer 返回错误时正确透传。

#### A5. `internal/model` 纯逻辑

优先覆盖容易改变生成结果的分支：

- `Config.Preprocess/GetNames/GetSchemaName/GetModelMethods`。
- option 的 modify/filter/create/method 分类。
- `Column.GetDataType/ToField` 的 nullable、coverable、unsigned、deleted_at、scan type。
- default tag 的零值、时间值、空白字符串和 `created_at/updated_at` 特例。
- 单行/多行 comment 清洗和 tag escaping。
- `GroupByColumn` 的 nil index、多列和 priority。
- `SQLBuffer` 空白归一化与 `Dump` reset。

优先使用小型 fake `gorm.ColumnType/gorm.Index`，不连数据库。

#### A6. Incremental/MergeQuery manifest 端到端

现有测试只覆盖局部 helper，需要补真实生成生命周期：

- 首次生成写文件和 manifest，第二次相同输入不覆写。
- 用户修改已记录文件后的现有保护语义。
- A 集合生成后以 B 集合执行 `MergeQuery`，最终 `gen.go` 与 manifest 同时保留 A+B。
- `MergeQuery` mode 不兼容时返回明确错误，原文件不被破坏。
- manifest 不存在、字段缺失、损坏 JSON 的行为。
- 多 model 并发输出在 `-race` 下稳定，manifest 文件集合和 hash 完整。
- 新增/删除生成文件时的预期清理策略；若当前没有清理能力，先把现状和非目标写清楚。

所有文件测试使用 `t.TempDir()`，断言内容与 manifest，不依赖 mtime 精度。

#### A7. 延后或按回归补充

- `condition.go` 已有较高间接覆盖，只补 nil、unsupported condition、error propagation 等未覆盖分支。
- `field` 只补 relation option 累积、serializer、dynamic column/subquery 等高状态路径，不穷举同构数字/时间函数。
- parser/generator fuzz 暂不立项；先把真实 bug 和 grammar fixture 变成稳定 seed，再评估纯 `Section.BuildSQL` 入口的 seeded fuzz。

### B. 测试基础设施

#### B1. 标准库目录 golden comparator

在 `tests/internal/golden` 实现 `CompareDirs(want, got) error`：

- 双向遍历，能发现缺失文件、多余文件和文件/目录类型不一致。
- 以相对路径排序后比较，输出确定性错误。
- 内容不同时至少给出文件名和首个差异行；不引入第三方库。
- 忽略项必须显式传入或硬编码为最小集合，例如运行时 DB 文件，不能静默忽略任意文件。

比较器自身用 `t.TempDir()` 覆盖：相同目录、内容不同、want only、got only、嵌套目录、空文件。

`matchGeneratedFile` 改用该 helper，删除外部 `diff`、context timeout 和进程调用。验收通过 helper 单测完成，不手工篡改仓库 golden。

#### B2. 收敛 MySQL 根包的隐式初始化

将 tests 根包的全局 `init()` 改为显式 `TestMain`/setup：

- 不再在源码中覆盖 `GORM_DIALECT`。
- 根包明确只接受 MySQL；CI 必须显式提供 dialect/DSN。
- 连接、Ping、migration、generation 任一步失败都使测试失败，不能只打印后继续。
- setup 返回 cleanup，关闭连接并清理本轮生成的临时产物。
- generation 尽量输出到每个测试的临时目录；必须共享的部分明确串行。

这一步只整理测试 harness，不把现有 MySQL golden 改造成 SQLite golden。

### C. DB-independent 可执行 fixture 与生成契约

#### C1. fixture 设计

在独立 package 定义 `User/Company/Order` 等最小模型，避免 generated query 反向导入测试包造成 import cycle。DIY interface 只保留语义完整的方法。

每个 DSL 控制结构至少有一个可执行代表：

- `{{if}}/{{else}}`：可选过滤与明确的两条结果分支。
- `{{for}}`：列表 OR 条件；空列表也必须产生合法 SQL 或由接口契约明确拒绝。
- `{{where}}`：多个可选条件，验证前导 AND/OR 清理。
- map loop：循环体显式写 AND，以结果语义断言，不断言 map 顺序。
- `{{set}}`：按可选字段形成单次 UPDATE，最后带稳定主键条件。
- `{{trim}}`：清理尾部 OR/逗号，空/非空输入都有定义。
- dialect-neutral INSERT：不包含 `ON DUPLICATE KEY UPDATE`。

运行时 fixture 方法原则上返回 error；若特意验证无 error 返回签名，必须用“预期命中非空数据”避免 SQL 错误被空结果掩盖。

#### C2. generation contract

`tests/generation` 使用 struct + interface 生成到 `t.TempDir()`：

1. 与可被 runtime 直接 import 的 `fixture/query` 做目录级比较，避免 golden 与可执行副本漂移。
2. runtime package 显式 import checked-in query，使 CI 编译 generated code。
3. golden 变化必须由生成命令产生并在 PR 中单独说明；不手改 `.gen.go`。

新增 runtime fixture 本身属于有意新增 golden，不受旧计划“不改任何 `.expect`”的限制；现有 `dal_*` golden 只有生成逻辑确实变化时才更新。

### D. SQLite runtime contract

#### D1. 数据库生命周期

`tests/runtime` 提供 `newSQLiteDB(t)`：

- 每个测试使用 `t.TempDir()` 下独立数据库文件，避免 `:memory:` 多连接看不到同一 schema。
- 使用 fixture model `AutoMigrate` 或最小方言无关 DDL，迁移失败立即 `t.Fatal`。
- 启用 foreign keys，关闭底层 `sql.DB` 注册到 `t.Cleanup`。
- 使用 `q := query.Use(db)`，不调用全局 `SetDefault`。
- 只有数据库和 query instance 完全隔离的测试才允许 `t.Parallel()`。

#### D2. DIY DSL 执行矩阵

对 C1 的每个方法：

- 预置能区分分支的确定性数据。
- 同时覆盖非空和关键空输入。
- 断言返回集合、字段值、rows affected 和 error。
- 对 map 只断言语义，不断言生成顺序。
- 对会丢弃 error 的签名必须断言预期命中，禁止以“结果为空”为成功条件。

#### D3. generated DO 运行时矩阵

按 gen 自有适配行为组织，不为上游 GORM 的每个 wrapper 写一条测试：

1. **查询分页**：`Order + Limit + Offset + Find`，断言稳定顺序和边界。
2. **聚合投影**：`Distinct/Count/Pluck/Scan`，以及一条 `Group + Having`。
3. **批处理**：`FindInBatches` 的批次号、批大小、完整遍历和 callback error 提前停止。
4. **写入结果**：`Update/Updates/UpdateColumn/Delete` 的值、rows affected 和 error。
5. **写入安全**：无 WHERE 更新/删除的行为，避免测试为了通过而关闭 GORM global-update 防线。
6. **Scopes**：scope 组合后结果正确且不污染原 query。
7. **关系**：一条 `Preload` 和一条 relation `Joins`，覆盖 generated relation field 到 GORM 的适配。
8. **事务**：commit、rollback 和 callback error；使用 query instance，不依赖全局 Q。
9. **错误传播**：取消 context、无效目标或 callback error 能从 generated API 返回。

### E. MySQL contract

现有 MySQL 5.7/latest job继续负责：

- 从真实 schema introspect model/query。
- 现有 `dal_*` golden 比较。
- 现有 CRUD 和 association transaction。

新增一个独立 MySQL runtime 表/fixture，覆盖 SQLite 无法证明的最小方言集合：

- `ON DUPLICATE KEY UPDATE`。
- `concat` 等确实存在于公开 fixture 的 MySQL 表达式。
- 如 generator 对 5.7 与 latest 产生差异，明确各自预期，不用 SQLite 守卫把测试静默跳过。

不要直接执行当前 `TestForInSet`、`AddUser` 等 grammar fixture；它们存在裸 `where`、缺少 `age` 列等已知不满足运行契约的问题。

### F. gentool

先做小型可测试性重构，再补测试：

```go
parseArgs(fs *flag.FlagSet, args []string) (*CmdParams, error)
loadConfig(path string) (*CmdParams, error)
dialectorFor(dbType DBType, dsn string) (gorm.Dialector, error)
run(args []string) error
```

要求：

- `loadConfig` 不调用 `log.Fatalf`；错误交给 `main` 决定退出。
- 不使用全局 `flag.CommandLine`，测试间不修改全局 `os.Args`。
- 明确 `-c` 是“配置文件独占”还是“CLI 覆盖 YAML”，保持或修改行为都必须有契约测试。
- `revise` 只测试它实际负责的默认值和 table 清理；逗号拆分属于参数解析测试。
- `dialectorFor` 测所有支持类型和非法类型，不连接外部数据库。
- `run` 用预先创建 schema 的 SQLite 临时文件做一次 smoke，断言输出文件存在且可解析/编译。
- gentool 使用 Go 1.25，并依赖已发布 gen；其结果不合并到 root 源码覆盖率。

## 5. PR 拆分与依赖

| PR | 内容 | 依赖 | 完成判据 |
|---|---|---|---|
| 1 | sec_check 安全契约；必要时同 PR 修安全缺口 | 无 | root unit/race 通过，安全决策写入 PR |
| 2 | helper、pools、diagnostic、model 单元测试 | 无 | 无 DB、无 sleep 型脆弱测试、race 通过 |
| 3 | Incremental/MergeQuery manifest E2E | 无 | 两轮真实生成、损坏输入、并发文件集合通过 |
| 4 | golden comparator + MySQL TestMain/harness | 无 | comparator 自测；现有 MySQL golden/CRUD 不回归 |
| 5 | executable fixture + DB-independent generation | PR 4 | 新 golden 由生成器产生、目录一致、编译通过 |
| 6 | SQLite DIY + DO runtime | PR 5 | 无容器执行，结果/error/rows 均有断言，race 通过 |
| 7 | gentool 可测试性重构和 SQLite smoke | 无 | gentool module test 通过，不访问外部 DB |
| 8 | CI 与分层 coverage artifact | PR 1-7 | 各 job 职责清晰、artifact 唯一、无门禁 |

PR 1、2、3、7 可以独立推进；PR 4 → 5 → 6 → 8 是主依赖链。不要把全部内容合成一个 coverage PR。

## 6. CI 设计

### 6.1 Jobs

1. **root-unit**：保留 Go 1.22-1.26 matrix，执行 `go test -race ./...`。
2. **sqlite-runtime**：在 canonical Go 版本执行 generation/runtime 子包，不启动容器。
3. **mysql-contract**：保留 MySQL 5.7/latest matrix，运行 tests 根包和 MySQL 专属用例。
4. **nested-modules**：
   - `go -C examples test ./...`，作为 compile smoke。
   - `go -C tools/gentool test ./...`，Go 1.25+。
5. **tidy**：继续逐 module 检查。

tests module 当前有 `toolchain go1.24.4`，不要把 SQLite job宣传成“真实 Go 1.22 集成验证”；最低版本兼容性仍由 root matrix 负责。是否移除 toolchain 指令是独立维护决策。

### 6.2 Coverage

只在一个 canonical Go 版本生成，避免 matrix artifact 冲突：

```bash
go test -covermode=atomic -coverprofile=coverage-root.out ./...
go -C tests test \
  -covermode=atomic \
  -coverpkg=gorm.io/gen,gorm.io/gen/field,gorm.io/gen/helper,gorm.io/gen/internal/... \
  -coverprofile=coverage-runtime.out \
  ./runtime/...
```

- 分别上传 `coverage-root-<go-version>` 和 `coverage-runtime-<go-version>`。
- 第一轮只上传 artifact，不接第三方、不设置百分比门禁。
- PR 说明列出变更前后相关 package/function 差异，不用总覆盖率掩盖关键路径。

### 6.3 格式与 lint

- 所有新增 Go 文件运行 `goimports`。
- 生产源码重构运行 `golangci-lint run --timeout 5m`。
- 当前 lint 配置不检查测试文件；测试质量由 goimports、编译、race、review 和精确断言保证。是否开启 `run.tests` 另行评估，避免一次引入全仓历史 lint 噪声。

## 7. 验收命令

### 无 DB

```bash
go test ./...
go test -race ./...
go -C tests test ./generation/... ./internal/golden/...
go -C tools/gentool test ./...
go -C examples test ./...
golangci-lint run --timeout 5m
git diff --check
```

### SQLite runtime

```bash
go -C tests test ./runtime/...
go -C tests test -race -count=1 ./runtime/...
```

### MySQL contract

先启动 `tests/docker-compose.yml` 中的数据库，再执行：

```bash
GITHUB_ACTION=true \
GORM_DIALECT=mysql \
GEN_DSN='<disposable-test-dsn>' \
./tests/test.sh
```

不得把真实 DSN 写入文档、测试或提交记录。

## 8. 测试编写约束

- 测试名称描述行为和条件，不使用只对应函数名的模糊名称。
- 纯逻辑优先表驱动；状态机、文件生命周期、事务不为了形式强行表驱动。
- 默认断言公开结果、error、rows affected 和文件集合，不断言私有调用次数。
- 不依赖 map 遍历顺序、目录遍历顺序、mtime 精度或固定 sleep。
- 数据库测试每例独立建库、迁移、seed、cleanup；迁移和清理错误不能吞掉。
- 不使用全局 `query.SetDefault` 驱动可并行测试。
- 不手改 generated `.gen.go`；新增 checked-in fixture 必须由生成流程产生并经 diff 审核。
- golden comparator 必须同时发现内容变化、缺失文件和多余文件。
- 测试发现源码 bug 时，先写能稳定复现的失败用例，再在同一逻辑 PR 修复；PR 中明确区分“测试补充”和“行为变化”。

## 9. 非目标

- 本轮不新增 PostgreSQL/SQL Server 容器矩阵。
- 不追求全仓、field 或 template 的任意覆盖率百分比。
- 不把现有 grammar fixture 全部改造成可执行 SQL。
- 不用 SQLite 结果宣称 MySQL schema introspection 已验证。
- 不在首轮接入 Codecov 或设置 coverage gate。
- 不借测试重构改变 gentool 配置优先级、安全策略或 generated public API，除非单独确认并在 PR 中明确说明。

## 10. 最终完成定义

只有同时满足以下条件，才能把整个计划标记为完成：

- 上述 PR 的实际代码均已合入，而不只是计划已写完。
- root、SQLite、MySQL、gentool、examples 的对应验证全部通过。
- 新增 fixture 确认由生成器产生且能编译、执行。
- CI 中能区分 unit、generation、SQLite runtime、MySQL contract 和 nested modules 的状态。
- coverage artifact 能准确说明采集范围，且未设置未经维护者同意的门禁。
- PR 文档明确列出测试代码、生产修复、generated 文件和兼容性变化。

当前 worktree 的实现与本地验收状态：

- root、generation/golden、SQLite runtime、gentool 和 examples 测试通过；对应 race 路径已运行。
- MySQL latest 容器 contract 通过；MySQL 5.7 留给 CI matrix 验证。
- root module 聚合 coverage 为 41.1%，SQLite runtime 对明确列出的 root packages 覆盖为 7.4%；只上传 artifact，不设门禁。
- gentool module lint 为 0 issue；root 与 tests module 仍有本轮之前即存在的 lint baseline，且本轮未改对应生产文件。
- GitHub Actions 尚未在本地 worktree 上实际运行，因此不把“CI 已通过”或“已合入”写成完成状态。
