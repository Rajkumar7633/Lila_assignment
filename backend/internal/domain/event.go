package domain

import "regexp"

var uuidPattern = regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)

func IsBot(userID string) bool {
	return !uuidPattern.MatchString(userID)
}

func EventCategory(eventType string) string {
	switch eventType {
	case "Position", "BotPosition":
		return "movement"
	case "Kill", "BotKill":
		return "kill"
	case "Killed", "BotKilled":
		return "death"
	case "KilledByStorm":
		return "storm"
	case "Loot":
		return "loot"
	default:
		return "other"
	}
}

func IsMovementEvent(eventType string) bool {
	return eventType == "Position" || eventType == "BotPosition"
}

func IsKillEvent(eventType string) bool {
	return eventType == "Kill" || eventType == "BotKill"
}

func IsDeathEvent(eventType string) bool {
	return eventType == "Killed" || eventType == "BotKilled" || eventType == "KilledByStorm"
}
