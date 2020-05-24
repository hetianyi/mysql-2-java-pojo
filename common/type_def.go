package common

// 配置定义
type Config struct {
	Host                   string            `yaml:"host"`
	Port                   int               `yaml:"port"`
	User                   string            `yaml:"user"`
	Password               string            `yaml:"password"`
	DB                     string            `yaml:"database"`
	Options                map[string]string `yaml:"connection-options"`
	Author                 string            `yaml:"author"`
	Version                string            `yaml:"version"`
	Package                string            `yaml:"package"`
	AddSerializeAnnotation bool              `yaml:"add-serialize-annotation"`
	DateFormat             string            `yaml:"date-format"`
	IgnoreEmptyField       bool              `yaml:"ignore-empty-field"`
	UseLombok              bool              `yaml:"use-lombok"`
	UseMybatisPlus         bool              `yaml:"mybatis-plus"`
	BeanSuffix             string            `yaml:"bean-suffix"`
	Tables                 []string          `yaml:"tables"`
	IgnoreTablePrefix      []string          `yaml:"ignore-tab-prefix"`
	IgnoreTableSuffix      []string          `yaml:"ignore-tab-suffix"`
}

// 列定义
type Column struct {
	ColumnName    string `gorm:"column:COLUMN_NAME"`
	ColumnType    string `gorm:"column:DATA_TYPE"`
	IsId          bool   `gorm:"column:"`
	AutoIncrement bool   `gorm:"column:EXTRA"`
	Comment       string `gorm:"column:COLUMN_COMMENT"`
}

// 列定义
type ColumnTmp struct {
	ColumnName    string `gorm:"column:COLUMN_NAME"`
	ColumnType    string `gorm:"column:DATA_TYPE"`
	IsId          string `gorm:"column:COLUMN_KEY"`
	AutoIncrement string `gorm:"column:EXTRA"`
	Comment       string `gorm:"column:COLUMN_COMMENT"`
}
