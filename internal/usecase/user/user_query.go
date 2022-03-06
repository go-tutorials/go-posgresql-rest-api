package user

import (
	"fmt"
	"github.com/lib/pq"
	"strings"
)

func BuildQuery (sm interface{}) (str string, param []interface{}) {
	str = `select * from users`
	u := sm.(*UserFilter)
	var where []string

	i := 1
	if u.Interests != nil && len(u.Interests) > 0  {
		param = append(param, pq.Array(u.Interests))
		where = append(where, fmt.Sprintf(` interests && $%d`, i))
		i++
	}
	if u.Skills != nil  && len(u.Skills) > 0 {
		var skills []string
		for _, value := range u.Skills {
			param = append(param, value)
			skills = append(skills, fmt.Sprintf(` $%d <@ ANY(skills)`, i))
			i ++
		}
		where = append(where, fmt.Sprintf(`(%s)`, strings.Join(skills, " or")))
	}
	if u.Settings != nil {
		param = append(param, u.Settings)
		where = append(where, fmt.Sprintf(` settings && $%d`, i))
		i++
	}
	if u.Achievements != nil && len(u.Achievements) > 0 {
		var achievements []string
		for _, value := range u.Achievements {
			param = append(param, value)
			achievements = append(achievements, fmt.Sprintf(` $%d <@ ANY(achievements)`,i))
			i++
		}
		where = append(where, fmt.Sprintf(`(%s)`, strings.Join(achievements, " or")))
	}

	if len(where) > 0 {
		str = str + ` where` + fmt.Sprintf(` %s`, strings.Join(where, " and"))
	}

	fmt.Println(str)
	fmt.Println(param)
	return
}
