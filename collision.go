package main

func otherEntitiesAt(entity *PlacedEntity, x, y int) []*PlacedEntity {
	collisions := entitiesAt(x, y)
	for i, e := range collisions {
		if e == entity {
			// remove that one
			collisions = append(collisions[:i], collisions[i+1:]...)
			return collisions
		}
	}
	return collisions
}

func entitiesAt(x, y int) []*PlacedEntity {
	collisions := make([]*PlacedEntity, 0)

	for _, other := range monsters {
		if other.X == x && other.Y == y {
			collisions = append(collisions, other)
		}
	}

	if character.X == x && character.Y == y {
		collisions = append(collisions, character)
	}

	return collisions
}
