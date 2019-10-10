package main

type Direction int

const (
	LEFT Direction = iota
	RIGHT
	UP
	DOWN
)

func swipe(d Direction) {
	if d == LEFT {
		for y := 0; y <= MAX_Y; y++ {
			swipeMonsters(LEFT, y)
		}
	}
	if d == RIGHT {
		for y := 0; y <= MAX_Y; y++ {
			swipeMonsters(RIGHT, y)
		}
	}
}

func reverseEntityList(list []*PlacedEntity) []*PlacedEntity {
	length := len(list)
	result := make([]*PlacedEntity, length, length)
	for i, entity := range list {
		result[length-1-i] = entity
	}
	return result
}

func swipeMonsters(d Direction, y int) {
	/*
		for each monster in the row
		left to right.
		we move it as far left as we can.
		we check for the monster to its left
		if its the same type
		we merge it in.

		if we didn't merge, we start over with that next left monster
		if we did merge, we go on to the 3rd monster in the row.
		...

	*/

	monstersInRow := make([]*PlacedEntity, 0)
	for x := 0; x <= MAX_X; x++ {
		monstersInCell := otherEntitiesAt(character, x, y)
		if len(monstersInCell) != 0 {
			monstersInRow = append(monstersInRow, monstersInCell[0])
		}
	}
	// now we have the short list of monsters.

	// if right, reverse list
	if d == RIGHT {
		monstersInRow = reverseEntityList(monstersInRow)
	}

	for i := 0; i < len(monstersInRow); i++ {
		monster := monstersInRow[i]
		// is there another monster?
		if i+1 < len(monstersInRow) {
			otherMonster := monstersInRow[i+1]
			// are they the same type
			if otherMonster.Sprite == monster.Sprite {
				// merge them
				if upgrade(monster) {
					// remove second monster
					monstersInRow = append(monstersInRow[:i+1], monstersInRow[i+2:]...)
				}
			}
		}
	}

	if d == RIGHT {
		monstersInRow = reverseEntityList(monstersInRow)
	}

	// redistribute the x values
	prefix := 0
	if d == RIGHT {
		prefix = CELLS_PER_BOARD - len(monstersInRow)
	}
	for i, monster := range monstersInRow {
		monster.X = prefix + i
	}
}
