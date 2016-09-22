package common

import (
	"bytes"
	"fmt"
)

type sqlBuilderWhere struct {
	k string
	v interface{}
}
type sqlBuilderLimit struct {
	pn int64
	ps int64
}
type SqlBuilder struct {
	selec   string
	from    string
	wheres  []sqlBuilderWhere
	limit   *sqlBuilderLimit
	orderBy string
	Args    []interface{}
}

func NewSqlBuilder() *SqlBuilder {
	return &SqlBuilder{}
}

func (sb *SqlBuilder) Select(s string) *SqlBuilder {
	sb.selec = s
	return sb
}

func (sb *SqlBuilder) From(s string) *SqlBuilder {
	sb.from = s
	return sb
}
func (sb *SqlBuilder) Where(k string, v interface{}) *SqlBuilder {
	if v != nil {
		sb.Args = append(sb.Args, v)
	}
	sb.wheres = append(sb.wheres, sqlBuilderWhere{k, v})
	return sb
}
func (sb *SqlBuilder) Limit(pn, ps int64) *SqlBuilder {
	sb.limit = &sqlBuilderLimit{pn, ps}
	return sb
}
func (sb *SqlBuilder) OrderBy(s string) *SqlBuilder {
	sb.orderBy = s
	return sb
}
func (sb *SqlBuilder) SelectSql() string {
	s := bytes.Buffer{}
	s.WriteString("select " + sb.selec)

	s.WriteString(sb.subSql())

	if sb.orderBy != "" {
		s.WriteString(fmt.Sprintf(" order by %s", sb.orderBy))
	}
	if sb.limit != nil {
		s.WriteString(fmt.Sprintf(" limit %d,%d", sb.limit.ps*(sb.limit.pn-1), sb.limit.ps*sb.limit.pn))
	}
	return s.String()
}
func (sb *SqlBuilder) subSql() string {
	s := bytes.Buffer{}
	s.WriteString(" from " + sb.from)
	first := true
	for _, wh := range sb.wheres {
		if first {
			s.WriteString(" where 1=1 ")
		}
		s.WriteString(" " + wh.k)
		first = false
	}

	return s.String()
}
func (sb *SqlBuilder) CountSql() string {
	s := bytes.Buffer{}
	s.WriteString("select count(*)")
	s.WriteString(sb.subSql())
	return s.String()
}
