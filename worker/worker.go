package worker

import (
	"bytes"
	"container/list"
	"github.com/hetianyi/gox"
	"github.com/hetianyi/gox/convert"
	"github.com/hetianyi/gox/file"
	"github.com/hetianyi/gox/logger"
	"github.com/hetianyi/mysql-2-java-pojo/common"
	"regexp"
	"strings"
	"time"
)

var config *common.Config
var outputDir string
var classFileBaseDir string

func Start(o string, c *common.Config) {
	outputDir = o
	config = c

	var buffer bytes.Buffer
	buffer.WriteString(c.User)
	buffer.WriteString(":")
	buffer.WriteString(c.Password)
	buffer.WriteString("@tcp(")
	buffer.WriteString(c.Host)
	buffer.WriteString(":")
	buffer.WriteString(convert.IntToStr(c.Port))
	buffer.WriteString(")/")
	buffer.WriteString(c.DB)

	if len(c.Options) > 0 {
		index := 0
		for k, v := range c.Options {
			if index == 0 {
				buffer.WriteString("?")
			} else {
				buffer.WriteString("&")
			}
			buffer.WriteString(k)
			buffer.WriteString("=")
			buffer.WriteString(v)
			index++
		}
	}

	connectString := buffer.String()

	logger.Info("connection string: ", connectString)

	InitMysqlClientConnection(connectString)

	ls := getTables(c.DB)

	allTableInfos := make(map[string]map[string]common.Column)
	orderedColumns := make(map[string][]string)

	gox.WalkList(ls, func(item interface{}) bool {
		m := item.(map[string]string)
		r, o := parseTable(c.DB, m)
		allTableInfos[m["TABLE_NAME"]] = r
		orderedColumns[m["TABLE_NAME"]] = o
		return false
	})

	createPackage(outputDir, c.Package)

	gox.WalkList(ls, func(item interface{}) bool {
		m := item.(map[string]string)
		createBean(m, allTableInfos[m["TABLE_NAME"]], orderedColumns)
		return false
	})
}

func getTables(dbName string) *list.List {
	s := "select TABLE_NAME, TABLE_TYPE, TABLE_COMMENT from information_schema.`TABLES` a where a.TABLE_SCHEMA = '" + dbName + "'"
	var err error
	var ls = list.New()
	gox.Try(func() {
		rows, err := db.Raw(s).Rows()
		if err != nil {
			err = transformNotFoundErr(err)
			return
		}
		for rows.Next() {
			var tn, tt, tc string
			if err = rows.Scan(&tn, &tt, &tc); err != nil {
				logger.Fatal(err)
			}
			ls.PushBack(map[string]string{
				"TABLE_NAME":    tn,
				"TABLE_TYPE":    tt,
				"TABLE_COMMENT": tc,
			})
		}

	}, func(e interface{}) {
		logger.Error(e)
		err = e.(error)
	})
	if err != nil {
		logger.Fatal(err)
	}
	return ls
}

func parseTable(dbName string, tableInfo map[string]string) (map[string]common.Column, []string) {
	tn := tableInfo["TABLE_NAME"]
	tt := tableInfo["TABLE_TYPE"]
	if isBaseTable(tt) {
		tableSchema := "select a.* from information_schema.`COLUMNS` a where a.TABLE_SCHEMA='" + dbName + "' and a.TABLE_NAME = '" + tn + "'"

		var err error
		var ret []common.ColumnTmp
		var ls = make(map[string]common.Column)
		var orderedColumns []string
		gox.Try(func() {
			if err := db.Raw(tableSchema).Find(&ret).Error; err != nil {
				logger.Fatal(err)
			}
			if err != nil {
				err = transformNotFoundErr(err)
				return
			}

			orderedColumns = make([]string, len(ret))
			for i, v := range ret {
				c := common.Column{
					ColumnName:    v.ColumnName,
					ColumnType:    v.ColumnType,
					Comment:       v.Comment,
					AutoIncrement: strings.Contains(v.AutoIncrement, "auto_increment"),
					IsId:          strings.Contains(v.IsId, "PRI"),
				}
				ls[v.ColumnName] = c
				orderedColumns[i] = v.ColumnName
			}

		}, func(e interface{}) {
			logger.Error(e)
			err = e.(error)
		})
		if err != nil {
			logger.Fatal(err)
		}
		return ls, orderedColumns
	}
	return nil, nil
}

func createPackage(outputDir, packageName string) {
	if err := file.CreateDirs(outputDir); err != nil {
		logger.Fatal(err)
	}
	packageDirs := strings.Split(packageName, ".")
	nextDir := outputDir
	for _, v := range packageDirs {
		if err := file.CreateDirs(nextDir + "/" + v); err != nil {
			logger.Fatal(err)
		}
		nextDir += "/" + v
	}
	classFileBaseDir = nextDir
}

func createBean(tableInfo map[string]string, columns map[string]common.Column, orderedColumns map[string][]string) {

	contains := false
	hasValidTable := false
	for _, v := range config.Tables {
		if strings.TrimSpace(v) != "" {
			hasValidTable = true
			break
		}
	}
	if hasValidTable {
		for _, v := range config.Tables {
			if strings.TrimSpace(v) != "" && tableInfo["TABLE_NAME"] == v {
				contains = true
				break
			}
		}
		if !contains {
			return
		}
	}

	beanName := CamelIt(tableInfo["TABLE_NAME"], true) + config.BeanSuffix
	classFileName := beanName + ".java"

	out, err := file.CreateFile(classFileBaseDir + "/" + classFileName)
	if err != nil {
		logger.Fatal(err)
	}
	defer out.Close()

	logger.Info("export ", tableInfo["TABLE_NAME"], " -> ", beanName)

	// file content ---------------------------------------------
	packageLine := "package " + config.Package + ";\n\n"
	imports := make(map[string]interface{})
	var fileBody bytes.Buffer
	var classBody bytes.Buffer
	var getset bytes.Buffer

	fileBody.WriteString(packageLine)

	classBody.WriteString("\n/**\n * ")
	classBody.WriteString(tableInfo["TABLE_COMMENT"])
	classBody.WriteString("\n * <br>Table : ")
	classBody.WriteString(tableInfo["TABLE_NAME"])
	classBody.WriteString("\n * <br>Generated by 'github.com/hetianyi/mysql-2-java-pojo'")
	classBody.WriteString("\n * @author ")
	classBody.WriteString(config.Author)
	classBody.WriteString("\n * @version ")
	classBody.WriteString(config.Version)
	classBody.WriteString("\n * @date ")
	classBody.WriteString(gox.GetLongDateString(time.Now()))
	classBody.WriteString("\n */\n")

	if config.UseLombok {
		classBody.WriteString("@Data\n")
		classBody.WriteString("@AllArgsConstructor\n")
		classBody.WriteString("@NoArgsConstructor\n")
		imports["lombok.Data"] = nil
		imports["lombok.AllArgsConstructor"] = nil
		imports["lombok.NoArgsConstructor"] = nil
	}

	if config.UseMybatisPlus {
		classBody.WriteString("@TableName(\"" + tableInfo["TABLE_NAME"] + "\")\n")
		imports["com.baomidou.mybatisplus.annotation.TableName"] = nil
	}

	classBody.WriteString("public class " + beanName + " ")
	if config.AddSerializeAnnotation {
		classBody.WriteString("implements Serializable {\n\n")
		classBody.WriteString("    private static final long serialVersionUID = 1L;")
		imports["java.io.Serializable"] = nil
	} else {
		classBody.WriteString("{")
	}

	for _, cn := range orderedColumns[tableInfo["TABLE_NAME"]] {

		k := cn
		v := columns[k]

		fieldName := CamelIt(k, false)
		comment := v.Comment

		if comment != "" {
			classBody.WriteString("\n\n    /** ")
			if strings.Contains(comment, "\n") {
				classBody.WriteString("\n     * ")
				classBody.WriteString(strings.ReplaceAll(comment, "\n", "\n     * "))
				classBody.WriteString("\n     */")
			} else {
				classBody.WriteString(strings.ReplaceAll(comment, "\n", "\n    "))
				classBody.WriteString(" */")
			}
		} else {
			classBody.WriteString("\n")
		}

		if config.UseMybatisPlus {
			if v.IsId {
				classBody.WriteString("\n    @TableId(\"" + k + "\")")
				imports["com.baomidou.mybatisplus.annotation.TableId"] = nil
			} else {
				classBody.WriteString("\n    @TableField(\"" + k + "\")")
				imports["com.baomidou.mybatisplus.annotation.TableField"] = nil
			}
		}

		classBody.WriteString("\n    private ")
		imp, typeN := common.GetType(v.ColumnType)
		if imp != "" {
			imports[imp] = nil
		}
		classBody.WriteString(typeN)
		classBody.WriteString(" ")
		classBody.WriteString(fieldName)
		classBody.WriteString(";")

		getset.WriteString("\n    public ")
		getset.WriteString(typeN)
		getset.WriteString(" get")
		getset.WriteString(CamelIt(k, true))
		getset.WriteString("() {\n")
		getset.WriteString("        return this.")
		getset.WriteString(fieldName)
		getset.WriteString(";\n    }\n")
		getset.WriteString("    public void set")
		getset.WriteString(CamelIt(k, true))
		getset.WriteString("(")
		getset.WriteString(typeN)
		getset.WriteString(" ")
		getset.WriteString(fieldName)
		getset.WriteString(") {\n")
		getset.WriteString("        this.")
		getset.WriteString(fieldName)
		getset.WriteString(" = ")
		getset.WriteString(fieldName)
		getset.WriteString(";\n    }\n")
	}

	if !config.UseLombok {
		classBody.WriteString("\n\n    //---------------- Getters and Setters ----------------\\\\\n")
		classBody.WriteString(getset.String())
	}

	classBody.WriteString("\n}")

	fileBody.WriteString("\n")
	for k := range imports {
		fileBody.WriteString("import ")
		fileBody.WriteString(k)
		fileBody.WriteString(";\n")
	}

	fileBody.WriteString(classBody.String())
	out.Write(fileBody.Bytes())
}

// 是否真实表格
func isBaseTable(tabType string) bool {
	return "VIEW" != tabType
}

func CamelIt(input string, isTab bool) string {
	pattern := regexp.MustCompile("(_[^_])")
	s := pattern.ReplaceAllStringFunc(input, func(s string) string {
		return strings.ToUpper(s[1:2]) + s[2:]
	})
	if isTab {
		return strings.ToUpper(s[0:1]) + s[1:]
	}
	return s
}
