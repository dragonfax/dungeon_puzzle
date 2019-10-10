package main

type Direction int

const (
	LEFT Direction = iota
	RIGHT
	UP
	DOWN
)

func swipe(d Direction) {
	for z := 0; z <= MAX_X; z++ {
		swipeMonsters(d, z)
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

func isHorizontal(d Direction) bool {
	return d == LEFT || d == RIGHT
}

func isIncreasing(d Direction) bool {
	return d == RIGHT || d == DOWN
}

func extractAxis(d Direction, z int) []*PlacedEntity {
	monstersInRow := make([]*PlacedEntity, 0)
	for w := 0; w <= MAX_X; w++ {
		var x, y int
		if isHorizontal(d) {
			x = w
			y = z
		} else {
			x = z
			y = w
		}
		monstersInCell := otherEntitiesAt(character, x, y)
		if len(monstersInCell) != 0 {
			// there should never be more than one
			monstersInRow = append(monstersInRow, monstersInCell[0])
		}
	}
	// now we have the short list of monsters.
	return monstersInRow
}

func swipeMonsters(d Direction, z int) {
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

	monstersInRow := extractAxis(d, z)

	// if right, reverse list
	if isIncreasing(d) {
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
					// TODO actually remove from monsters global
					removeMonster(otherMonster)
				}
			}
		}
	}

	if isIncreasing(d) {
		monstersInRow = reverseEntityList(monstersInRow)
	}

	// redistribute the x values
	createAxis(d, z, monstersInRow)
}

func createAxis(d Direction, z int, monstersInRow []*PlacedEntity) {
	prefix := 0
	if isIncreasing(d) {
		prefix = CELLS_PER_BOARD - len(monstersInRow)
	}
	for i, monster := range monstersInRow {
		if isHorizontal(d) {
			monster.X = prefix + i
		} else {
			monster.Y = prefix + i
		}
	}
}
