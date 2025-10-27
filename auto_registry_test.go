package gen

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestAutoRegistryInitGeneration(t *testing.T) {
	// 创建临时目录
	tempDir := t.TempDir()

	// 创建测试数据库
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to create test database: %v", err)
	}

	// 创建简单的测试表
	err = db.Exec(`CREATE TABLE users (id INTEGER PRIMARY KEY, name TEXT)`).Error
	if err != nil {
		t.Fatalf("failed to create users table: %v", err)
	}

	err = db.Exec(`CREATE TABLE orders (id INTEGER PRIMARY KEY, amount DECIMAL)`).Error
	if err != nil {
		t.Fatalf("failed to create orders table: %v", err)
	}

	tests := []struct {
		name           string
		configTables   []string
		expectInitFunc map[string]bool // table -> shouldHaveInit
	}{
		{
			name:         "AllTables",
			configTables: []string{}, // 空数组表示所有表
			expectInitFunc: map[string]bool{
				"users":  true,
				"orders": true,
			},
		},
		{
			name:         "OnlyUsersTable",
			configTables: []string{"users"},
			expectInitFunc: map[string]bool{
				"users":  true,
				"orders": false,
			},
		},
		{
			name:         "BothTables",
			configTables: []string{"users", "orders"},
			expectInitFunc: map[string]bool{
				"users":  true,
				"orders": true,
			},
		},
		{
			name:         "NoTables",
			configTables: []string{"nonexistent"},
			expectInitFunc: map[string]bool{
				"users":  false,
				"orders": false,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			modelDir := filepath.Join(tempDir, tt.name, "model")

			// 配置生成器
			g := NewGenerator(Config{
				OutPath:      filepath.Join(tempDir, tt.name, "query"),
				ModelPkgPath: modelDir,
				Mode:         WithDefaultQuery | WithAutoRegistry,
			})

			// 根据测试配置启用自动注册
			if len(tt.configTables) == 0 {
				g.WithAutoRegistry() // 不传参数，所有表
			} else {
				g.WithAutoRegistry(tt.configTables...) // 传入指定表名
			}

			g.UseDB(db)
			g.GenerateAllTable()
			g.Execute()

			// 验证每个表的 init 函数生成情况
			for tableName, shouldHaveInit := range tt.expectInitFunc {
				t.Run(tableName, func(t *testing.T) {
					checkInitFunction(t, modelDir, tableName, shouldHaveInit)
				})
			}

			// 验证注册表文件是否生成
			registryFile := filepath.Join(modelDir, "gen.go")
			if _, err := os.Stat(registryFile); os.IsNotExist(err) {
				t.Errorf("registry file %s should exist", registryFile)
			}
		})
	}
}

// checkInitFunction 检查指定表的模型文件是否包含正确的 init 函数
func checkInitFunction(t *testing.T, modelDir, tableName string, shouldHaveInit bool) {
	fileName := filepath.Join(modelDir, tableName+".gen.go")

	// 检查文件是否存在
	if _, err := os.Stat(fileName); os.IsNotExist(err) {
		t.Errorf("model file %s does not exist", fileName)
		return
	}

	// 读取文件内容
	content, err := os.ReadFile(fileName)
	if err != nil {
		t.Errorf("failed to read file %s: %v", fileName, err)
		return
	}

	fileContent := string(content)

	// 检查是否包含 init 函数
	hasInitFunc := strings.Contains(fileContent, "func init() {")
	hasRegisterCall := strings.Contains(fileContent, "RegisterModel(")

	if shouldHaveInit {
		if !hasInitFunc {
			t.Errorf("file %s should contain 'func init() {' but doesn't", fileName)
		}
		if !hasRegisterCall {
			t.Errorf("file %s should contain 'RegisterModel(' call but doesn't", fileName)
		}

		// 验证 RegisterModel 调用格式
		if hasInitFunc && hasRegisterCall {
			expectedModelName := getExpectedModelName(tableName)
			expectedCall := "RegisterModel(&" + expectedModelName + "{}, TableName" + expectedModelName + ")"

			if !strings.Contains(fileContent, expectedCall) {
				t.Errorf("file %s should contain %s", fileName, expectedCall)
				t.Logf("Actual file content:\n%s", fileContent)
			}
		}
	} else {
		if hasInitFunc && hasRegisterCall {
			t.Errorf("file %s should not contain init function with RegisterModel call", fileName)
		}
	}
}

// getExpectedModelName 根据表名获取期望的模型名
func getExpectedModelName(tableName string) string {
	switch tableName {
	case "users":
		return "User"
	case "orders":
		return "Order"
	default:
		// 简单的首字母大写
		return strings.Title(tableName)
	}
}

// TestShouldEnableAutoRegistry 测试表过滤逻辑
func TestShouldEnableAutoRegistry(t *testing.T) {
	tests := []struct {
		name           string
		configuredList []string
		tableName      string
		expected       bool
	}{
		{
			name:           "EmptyList_AllTablesEnabled",
			configuredList: []string{},
			tableName:      "users",
			expected:       true,
		},
		{
			name:           "TableInList_ShouldEnable",
			configuredList: []string{"users", "orders"},
			tableName:      "users",
			expected:       true,
		},
		{
			name:           "TableNotInList_ShouldDisable",
			configuredList: []string{"users", "orders"},
			tableName:      "products",
			expected:       false,
		},
		{
			name:           "SingleTable_Match",
			configuredList: []string{"users"},
			tableName:      "users",
			expected:       true,
		},
		{
			name:           "SingleTable_NoMatch",
			configuredList: []string{"users"},
			tableName:      "orders",
			expected:       false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := &Generator{
				Config: Config{
					RegistryTableList: tt.configuredList,
				},
			}

			result := g.shouldEnableAutoRegistry(tt.tableName)
			if result != tt.expected {
				t.Errorf("shouldEnableAutoRegistry(%s) = %v, want %v",
					tt.tableName, result, tt.expected)
			}
		})
	}
}

// TestWithAutoRegistryConfig 测试配置方法
func TestWithAutoRegistryConfig(t *testing.T) {
	tests := []struct {
		name         string
		tableNames   []string
		expectedList []string
		expectedMode GenerateMode
	}{
		{
			name:         "NoTables",
			tableNames:   []string{},
			expectedList: []string{},
			expectedMode: WithAutoRegistry,
		},
		{
			name:         "SingleTable",
			tableNames:   []string{"users"},
			expectedList: []string{"users"},
			expectedMode: WithAutoRegistry,
		},
		{
			name:         "MultipleTables",
			tableNames:   []string{"users", "orders", "products"},
			expectedList: []string{"users", "orders", "products"},
			expectedMode: WithAutoRegistry,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := Config{}

			// 调用 WithAutoRegistry 方法
			cfg.WithAutoRegistry(tt.tableNames...)

			// 验证配置结果
			if cfg.Mode&WithAutoRegistry == 0 {
				t.Error("WithAutoRegistry mode should be enabled")
			}

			if len(cfg.RegistryTableList) != len(tt.expectedList) {
				t.Errorf("RegistryTableList length = %d, want %d",
					len(cfg.RegistryTableList), len(tt.expectedList))
			}

			for i, expected := range tt.expectedList {
				if i >= len(cfg.RegistryTableList) || cfg.RegistryTableList[i] != expected {
					t.Errorf("RegistryTableList[%d] = %s, want %s",
						i, cfg.RegistryTableList[i], expected)
				}
			}
		})
	}
}
