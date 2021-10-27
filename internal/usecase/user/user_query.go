package user

import (
	"fmt"
	"strings"
)

func BuildQuery (sm interface{}) (str string, param []interface{}) {
	str = `select * from users`
	u := sm.(*UserFilter)
	var where []string
	var orWhereAchievements []string
	var orWhereSkills []string
	i := 1
	if u.Interests != nil && len(u.Interests) > 0  {
		param = append(param, u.Interests)
		where = append(where, fmt.Sprintf(` interests && $%d`, i))
		i++
	}
	if u.Settings != nil {
		param = append(param, u.Settings)
		where = append(where, fmt.Sprintf(` settings && $%d`, i))
		i++
	}
	if u.Achievements != nil && len(u.Achievements) > 0 {
		for _, value := range u.Achievements {
			param = append(param, value)
			orWhereAchievements = append(orWhereAchievements, fmt.Sprintf(` $%d <@ ANY(achievements)`,i))
			i++
		}
	}
	if u.Skills != nil  && len(u.Skills) > 0 {
		for _, value := range u.Skills {
			param = append(param, value)
			orWhereSkills = append(orWhereSkills, fmt.Sprintf(` $%d <@ ANY(skills)`, i))
			i ++
		}
	}
	if len(orWhereAchievements) > 0 {
		where = append(where, fmt.Sprintf(`(%s)`, strings.Join(orWhereAchievements, " or")))
	}
	if len(orWhereSkills) > 0 {
		where = append(where, fmt.Sprintf(`(%s)`, strings.Join(orWhereSkills, " or")))
	}
	if len(where) > 0 {
		str = str + ` where` + fmt.Sprintf(` %s`, strings.Join(where, " and"))
	}

	//if u.Limit > 0 {
	//	str = str + fmt.Sprintf(` limit %d`, u.Limit)
	//}

	fmt.Println(str)
	fmt.Println(param)
	return
}
