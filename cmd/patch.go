package cmd

import (
	"strings"
)

/*
	修改移动速度
*/
const MOVE_TIMES float32 = 1.5

func playerweighttable(items [][]Item_t) [][]Item_t {
	row_n := len(items)
	col_n := len(items[0])

	for i := 1; i < row_n; i++ {
		for j := 1; j < col_n; j++ {
			switch items[i][j].tp {
			case 5:
				val := items[i][j].val.(float32)
				if val > 0 {
					items[i][j].val = val * MOVE_TIMES
				}
			}
		}
	}

	return items
}

func rebootplayerweighttable(items [][]Item_t) [][]Item_t {
	row_n := len(items)
	col_n := len(items[0])

	for i := 1; i < row_n; i++ {
		for j := 1; j < col_n; j++ {
			switch items[i][j].tp {
			case 5:
				if items[0][j].val.(string) == "_MoveSpeed" {
					val := items[i][j].val.(float32)
					if val > 0 {
						items[i][j].val = val * MOVE_TIMES
					}
				}
			}
		}
	}

	return items
}

/*
	强化不破碎，最大将1级
*/
func enchanttable(items [][]Item_t) [][]Item_t {
	row_n := len(items)
	col_n := len(items[0])

	for i := 1; i < row_n; i++ {
		for j := 1; j < col_n; j++ {
			switch items[i][j].tp {
			case 3:
				if strings.Contains(items[0][j].val.(string), "Down") {
					val := items[i][j].val.(int32)
					if val > 1 {
						val = 1
					}
					items[i][j].val = val
				}
			case 4:
				if items[0][j].val.(string) == "_BreakRatio" {
					items[i][j].val = float32(0)
				}
			}
		}
	}

	return items
}

/*
	升一级 2倍 智力，力量、敏捷、耐力、技能点
*/
func playerleveltable(items [][]Item_t) [][]Item_t {
	row_n := len(items)
	col_n := len(items[0])

	for i := 1; i < row_n; i++ {
		for j := 1; j < col_n; j++ {
			switch items[i][j].tp {
			case 3:
				val := items[i][j].val.(int32)

				if items[0][j].val.(string) == "_SkillPoint" {
					items[i][j].val = val * 2
				} else if items[0][j].val.(string) == "_Strength" {
					items[i][j].val = val * 2
				} else if items[0][j].val.(string) == "_Agility" {
					items[i][j].val = val * 2
				} else if items[0][j].val.(string) == "_Intelligence" {
					items[i][j].val = val * 2
				} else if items[0][j].val.(string) == "_Stamina" {
					items[i][j].val = val * 2
				}
			}
		}
	}

	return items
}

/*
	副本随便进
*/
func stageentertable(items [][]Item_t) [][]Item_t {
	row_n := len(items)
	col_n := len(items[0])

	for i := 1; i < row_n; i++ {
		for j := 1; j < col_n; j++ {
			switch items[i][j].tp {
			case 3:
				val := items[i][j].val.(int32)

				if strings.Contains(items[0][j].val.(string), "PartyOneNum") {
					val = -1
					items[i][j].val = val
				} else if strings.Contains(items[0][j].val.(string), "MaxUsableCoin") {
					val = -1
					items[i][j].val = val
				} else if strings.Contains(items[0][j].val.(string), "NeedItemCount") {
					val = 0
					items[i][j].val = val
				} else if strings.Contains(items[0][j].val.(string), "UseCoin") {
					val = 1
					items[i][j].val = val
				}
			}
		}
	}

	return items
}

/*
	地图没有进入次数限制
*/
func maptable(items [][]Item_t) [][]Item_t {
	row_n := len(items)
	col_n := len(items[0])

	for i := 1; i < row_n; i++ {
		for j := 1; j < col_n; j++ {
			switch items[i][j].tp {
			case 3:
				val := items[i][j].val.(int32)

				if items[0][j].val.(string) == "_VipClear" {
					val = 1000
					items[i][j].val = val
				} else if items[0][j].val.(string) == "_MaxClearCount" {
					val = 1000
					items[i][j].val = val
				}
			}
		}
	}

	return items
}

/*
	通关经验降低100倍
*/
func stagerewardtable(items [][]Item_t) [][]Item_t {
	row_n := len(items)
	col_n := len(items[0])

	for i := 1; i < row_n; i++ {
		for j := 1; j < col_n; j++ {
			switch items[i][j].tp {
			case 3:
				val := items[i][j].val.(int32)

				if items[0][j].val.(string) == "_AwardEXP" {
					val /= 100
					if val < 1 {
						val = 1
					}
					items[i][j].val = val
				}
			}
		}
	}

	return items
}

/*
	怪物经验减半
*/
func monstertable(items [][]Item_t) [][]Item_t {
	row_n := len(items)
	col_n := len(items[0])

	for i := 1; i < row_n; i++ {
		for j := 1; j < col_n; j++ {
			switch items[i][j].tp {
			case 3:
				val := items[i][j].val.(int32)

				if items[0][j].val.(string) == "_DeadExperience" {
					items[i][j].val = val / 2
				} else if items[0][j].val.(string) == "_CompleteExperience" {
					items[i][j].val = val / 2
				}
			}
		}
	}

	return items
}

/*
	50% cd减免，，最高cd50秒
	技能修炼等级降低4级
*/
func skillleveltable_character(items [][]Item_t) [][]Item_t {
	row_n := len(items)
	col_n := len(items[0])

	for i := 1; i < row_n; i++ {
		for j := 1; j < col_n; j++ {
			switch items[i][j].tp {
			case 3:
				val := items[i][j].val.(int32)

				if items[0][j].val.(string) == "_LevelLimit" {
					if val > 10 {
						items[i][j].val = val - 4
					}
				} else if items[0][j].val.(string) == "_DelayTime" {
					if val > 0 {
						val /= 2

						if val > 50000 {
							val = 50000
						}

						items[i][j].val = val
					}
				}
			}
		}
	}

	return items
}

/*
	大招不共用冷却
*/
func skilltable_character(items [][]Item_t) [][]Item_t {
	row_n := len(items)
	col_n := len(items[0])

	for i := 1; i < row_n; i++ {
		for j := 1; j < col_n; j++ {
			switch items[i][j].tp {
			case 3:
				val := items[i][j].val.(int32)

				if items[0][j].val.(string) == "_SkillGroup" {
					val = 0
					items[i][j].val = val
				} else if items[0][j].val.(string) == "_GlobalCoolTimePvE" {
					if val > 50000 {
						val = 50000

						items[i][j].val = val
					}
				}
			}
		}
	}

	return items
}
