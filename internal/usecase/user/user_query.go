package user

import (
	"fmt"
	q "github.com/core-go/sql"
	"github.com/lib/pq"
	"strings"
)

func BuildQuery(filter interface{}) (query string, params []interface{}) {
	query = `select * from users`
	s := filter.(*UserFilter)
	var where []string

	i := 1
	if s.Interests != nil && len(s.Interests) > 0  {
		params = append(params, pq.Array(s.Interests))
		where = append(where, fmt.Sprintf(`interests && %s`, q.BuildDollarParam(i)))
		i++
	}
	if s.Skills != nil  && len(s.Skills) > 0 {
		var skills []string
		for _, value := range s.Skills {
			params = append(params, value)
			skills = append(skills, fmt.Sprintf(`%s <@ ANY(skills)`, q.BuildDollarParam(i)))
			i ++
		}
		where = append(where, fmt.Sprintf(`(%s)`, strings.Join(skills, " or ")))
	}
	if s.Settings != nil {
		params = append(params, s.Settings)
		where = append(where, fmt.Sprintf(`settings && %s`, q.BuildDollarParam(i)))
		i++
	}
	if s.Achievements != nil && len(s.Achievements) > 0 {
		var achievements []string
		for _, value := range s.Achievements {
			params = append(params, value)
			achievements = append(achievements, fmt.Sprintf(`%s <@ ANY(achievements)`, q.BuildDollarParam(i)))
			i++
		}
		where = append(where, fmt.Sprintf(`(%s)`, strings.Join(achievements, " or ")))
	}
	if len(where) > 0 {
		query = query + ` where ` + strings.Join(where, " and ")
	}
	return
}
