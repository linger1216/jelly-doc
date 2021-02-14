package core

import (
	"bytes"
	"fmt"
	jsoniter "github.com/json-iterator/go"
	"github.com/lib/pq"
	"github.com/linger1216/go-utils/db/postgres"
	"github.com/linger1216/go-utils/snowflake"
	"github.com/linger1216/jelly-doc/src/server/pb"
	"strings"
	"time"
)

func NewBasisApiDDL() *BasisApiDDL {
	ret := &BasisApiDDL{Name: "api"}
	ret.columns = append(ret.columns, &MetaColumn{Name: "id", Type: "character varying", Primary: true})
	ret.columns = append(ret.columns, &MetaColumn{Name: "name", Type: "character varying", Index: true})
	ret.columns = append(ret.columns, &MetaColumn{Name: "description", Type: "character varying"})
	ret.columns = append(ret.columns, &MetaColumn{Name: "member_ids", Type: "character varying[]"})

	ret.columns = append(ret.columns, &MetaColumn{Name: "method", Type: "character varying"})
	ret.columns = append(ret.columns, &MetaColumn{Name: "url", Type: "character varying"})
	ret.columns = append(ret.columns, &MetaColumn{Name: "headers", Type: "character varying"})
	ret.columns = append(ret.columns, &MetaColumn{Name: "path_params", Type: "character varying"})
	ret.columns = append(ret.columns, &MetaColumn{Name: "url_params", Type: "character varying"})
	ret.columns = append(ret.columns, &MetaColumn{Name: "auth", Type: "character varying"})
	ret.columns = append(ret.columns, &MetaColumn{Name: "body", Type: "character varying"})

	ret.columns = append(ret.columns, &MetaColumn{Name: "timeout", Type: "int"})
	ret.columns = append(ret.columns, &MetaColumn{Name: "directories", Type: "character varying[]"})

	ret.columns = append(ret.columns, &MetaColumn{Name: "create_time", Type: "bigint", Default: `(date_part('epoch'::text, now()))::bigint`})
	ret.columns = append(ret.columns, &MetaColumn{Name: "update_time", Type: "bigint", Default: `(date_part('epoch'::text, now()))::bigint`})
	return ret
}

func (m *BasisApiDDL) Upsert(apis ...*pb.Api) (string, []interface{}) {
	cols := strings.Split(m.ColumnsString(), ",")
	size := len(apis)
	values := make([]string, 0, size)
	args := make([]interface{}, 0, size*len(cols))
	for i, v := range apis {
		if len(v.Id) == 0 {
			v.Id = snowflake.Generate()
		}

		var createTime, updateTime int64
		if v.CreateTime == 0 {
			createTime = time.Now().Unix()
		} else {
			createTime = v.CreateTime
		}

		if v.UpdateTime == 0 {
			updateTime = time.Now().Unix()
		} else {
			updateTime = v.UpdateTime
		}

		var headers, urlParams, pathParams, auth []byte
		var err error
		if len(v.Headers) > 0 {
			headers, err = jsoniter.ConfigFastest.Marshal(v.Headers)
			if err != nil {
				return "", nil
			}
		}

		if len(v.UrlParams) > 0 {
			urlParams, err = jsoniter.ConfigFastest.Marshal(v.UrlParams)
			if err != nil {
				return "", nil
			}
		}

		if len(v.PathParams) > 0 {
			pathParams, err = jsoniter.ConfigFastest.Marshal(v.PathParams)
			if err != nil {
				return "", nil
			}
		}

		if v.Auth != nil {
			auth, err = jsoniter.ConfigFastest.Marshal(v.Auth)
			if err != nil {
				return "", nil
			}
		}

		values = append(values, postgres.ValuesPlaceHolder(i*len(cols), len(cols)))
		args = append(args, v.Id, v.Name, v.Description, pq.Array(v.MemberIds),
			v.Method, v.Url, headers, pathParams, urlParams, auth, v.Body, v.Timeout, pq.Array(v.Directories),
			createTime, updateTime)
	}

	query := fmt.Sprintf(`insert into %s (%s) values %s %s`, m.Table(), m.ColumnsString(),
		strings.Join(values, ","), m.OnConflictDDL())
	return query, args
}

func (m *BasisApiDDL) List(in *pb.ListApiRequest) (string, []interface{}) {
	firstCond := true
	var buffer bytes.Buffer
	if in.Header > 0 {
		buffer.WriteString(fmt.Sprintf("select count(1) from %s", m.Table()))
	} else {
		buffer.WriteString(fmt.Sprintf("select %s from %s", m.Select(), m.Table()))
	}

	if len(in.Names) > 0 {
		query := fmt.Sprintf("%s name in (%s)", postgres.CondSql(firstCond), postgres.SqlStringIn(in.Names...))
		buffer.WriteString(query)
		firstCond = false
	}

	if in.Header == 0 {
		query := fmt.Sprintf(" offset %d limit %d", in.CurrentPage*in.PageSize, in.PageSize)
		buffer.WriteString(query)
	}

	buffer.WriteString(";")
	return buffer.String(), nil
}

func (m *BasisApiDDL) Delete(ids ...string) (string, []interface{}) {
	query := fmt.Sprintf("delete from %s where %s in (%s);", m.Table(), "id", postgres.SqlStringIn(ids...))
	return query, nil
}

func (m *BasisApiDDL) Get(ids ...string) (string, []interface{}) {
	query := fmt.Sprintf("select %s from %s where %s in (%s);", m.Select(), m.Table(), "id", postgres.SqlStringIn(ids...))
	return query, nil
}

type MetaColumn struct {
	Name    string `json:"Name"`
	Type    string `json:"type"` // character varying, bigint, integer, geometry, character varying[], integer[]
	Primary bool   `json:"primary"`
	Index   bool   `json:"index"`
	Unique  bool   `json:"unique"`
	Default string `json:"default"`
}

func (m *MetaColumn) ColumnDDL() string {
	var primary string
	if m.Primary {
		primary = "primary key"
	}

	var defaultVal string
	if len(m.Default) > 0 {
		defaultVal = "default " + m.Default
	}

	return fmt.Sprintf("%s %s %s %s", m.Name, m.Type, primary, defaultVal)
}

func (m *MetaColumn) Select() string {
	switch m.Type {
	case `character varying`:
		return m.Name
	case `bigint`:
		return m.Name
	case `integer`, `int`:
		return m.Name
	case `geometry`:
		return fmt.Sprintf("st_asgeojson(%s) as %s", m.Name, m.Name)
	case `character varying[]`:
		return fmt.Sprintf("array_to_string(%s, ',') as %s", m.Name, m.Name)
	case `integer[]`, `int[]`:
		return fmt.Sprintf("array_to_string(%s, ',') as %s", m.Name, m.Name)
	}
	return ""
}

func (m *MetaColumn) IndexDDL(table string) string {
	if m.Primary {
		return ""
	}
	unique := ""
	if m.Unique {
		unique = "unique"
	}
	engine := ""
	switch m.Type {
	case "character varying", "bigint", "integer", "character varying[]", "integer[]":
		engine = fmt.Sprintf("btree(%s)", m.Name)
	case "geometry":
		engine = fmt.Sprintf("gist (geography(%s))", m.Name)
	}
	return fmt.Sprintf("create %s index if not exists %s_%s_index ON %s using %s;", unique, table, m.Name, table, engine)
}

type BasisApiDDL struct {
	Name    string
	columns []*MetaColumn
}

func (m *BasisApiDDL) Select() string {
	arr := make([]string, len(m.columns))
	for i := range m.columns {
		arr[i] = m.columns[i].Select()
	}
	return strings.Join(arr, ",")
}

func (m *BasisApiDDL) Table() string {
	return m.Name
}

func (m *BasisApiDDL) ColumnsString() string {
	arr := make([]string, len(m.columns))
	for i := range m.columns {
		arr[i] = m.columns[i].Name
	}
	return strings.Join(arr, ",")
}

func (m *BasisApiDDL) CreateTableDDL() string {
	var buf bytes.Buffer
	buf.WriteString(fmt.Sprintf("create table if not exists %s", m.Name))
	buf.WriteString("(\n")
	for i := range m.columns {
		buf.WriteString(m.columns[i].ColumnDDL())
		if i < len(m.columns)-1 {
			buf.WriteByte(',')
		}
		buf.WriteByte('\n')
	}
	buf.WriteString(");\n")
	return buf.String()
}

func (m *BasisApiDDL) IndexTableDDL() []string {
	arr := make([]string, 0)
	for _, v := range m.columns {
		if v.Index {
			arr = append(arr, v.IndexDDL(m.Name))
		}
	}
	return arr
}

func (m *BasisApiDDL) DBPrimaryColumn() *MetaColumn {
	for i := range m.columns {
		if m.columns[i].Primary {
			return m.columns[i]
		}
	}
	return nil
}

func (m *BasisApiDDL) OnConflictDDL() string {
	primaryColumn := m.DBPrimaryColumn()
	var buf bytes.Buffer
	buf.WriteString(fmt.Sprintf("on conflict (%s)\n", primaryColumn.Name))
	buf.WriteString("do update set\n")

	for i, v := range m.columns {
		if v.Primary || v.Name == "create_time" {
			continue
		}
		if v.Name == "update_time" {
			buf.WriteString(fmt.Sprintf("update_time = GREATEST(%s.update_time, excluded.update_time)", m.Name))
		} else {
			buf.WriteString(fmt.Sprintf("%s = excluded.%s", v.Name, v.Name))
		}
		if i < len(m.columns)-1 {
			buf.WriteString(",")
		} else {
			buf.WriteString(";")
		}
		buf.WriteString("\n")
	}
	return buf.String()

}
