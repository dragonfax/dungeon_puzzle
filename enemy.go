package main

var ENEMY_ORDER []string = []string{
	"imp",
	"goblin",
	"masked_orc",
	"skelet",
	"muddy",
	"zombie",
	"ice_zombie",
	"orc_warrior",
	"orc_shaman",
	"wogol",
	"big_zombie",
	"ogre",
	"big_damon",
}

var monsters []*PlacedEntity

func findEmptyPosition() (x, y int) {
	for x := 0; x < MAX_X; x++ {
		for y := 0; x < MAX_Y; y++ {
			occupied := entitiesAt(x, y)
			if len(occupied) == 0 {
				return x, y
			}
		}
	}
	panic("board full")
}

func spawnMonster(monsters []*PlacedEntity) {
	x, y := findEmptyPosition()
	newMonster := &PlacedEntity{
		Sprite: spriteByName("ogre_idle_anim"),
		X:      x,
		Y:      y,
	}
	monsters = append(monsters, newMonster)
}
