# 枚举类型支持

从 PostgreSQL 数据库中自动生成枚举类型常量。

## 功能介绍

当使用 PostgreSQL 数据库时，go-gorm/gen 现在可以自动识别枚举类型字段，并为其生成相应的 Go 常量。这使得在代码中使用枚举值变得更加类型安全和方便。

## 工作原理

1. 当从 PostgreSQL 数据库生成模型时，gen 会查询数据库中的枚举类型定义
2. 对于每个使用枚举类型的表字段，gen 会生成对应的常量定义
3. 常量名称会根据表名和枚举值自动格式化为符合 Go 命名规范的形式

## 示例

假设在 PostgreSQL 数据库中有以下定义：

```sql
-- 创建枚举类型
CREATE TYPE user_status AS ENUM ('active', 'inactive', 'pending');

-- 创建使用枚举类型的表
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    status user_status NOT NULL DEFAULT 'pending'
);
```

使用 gen 生成模型代码后，会自动生成如下常量：

```go
// Enum values for User table
const (
    // user_status enum values
    UserActive   = "active"
    UserInactive = "inactive"
    UserPending  = "pending"
)

// User struct
type User struct {
    ID     int64  `gorm:"column:id;primaryKey;autoIncrement:true" json:"id"`
    Name   string `gorm:"column:name;not null" json:"name"`
    Status string `gorm:"column:status;not null;default:pending" json:"status"`
}
```

## 使用方法

在代码中，你可以直接使用生成的常量：

```go
user := User{
    Name:   "New User",
    Status: UserActive, // 使用生成的枚举常量
}

// 查询所有活跃用户
activeUsers, err := query.User.Where(query.User.Status.Eq(UserActive)).Find()
```

## 注意事项

1. 此功能仅适用于 PostgreSQL 数据库，因为它依赖于 PostgreSQL 的枚举类型支持
2. 枚举字段在 Go 中仍然是 string 类型，但使用生成的常量可以提高代码的可读性和类型安全性
3. 如果修改了数据库中的枚举类型定义，需要重新生成模型代码以更新常量定义