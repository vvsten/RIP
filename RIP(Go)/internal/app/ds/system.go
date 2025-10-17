package ds

// SystemUsers - singleton для системных пользователей
type SystemUsers struct {
	CreatorID   int
	ModeratorID int
}

var systemUsers *SystemUsers

// GetSystemUsers - получение singleton системных пользователей
func GetSystemUsers() *SystemUsers {
	if systemUsers == nil {
		systemUsers = &SystemUsers{
			CreatorID:   1, // ID создателя (зафиксирован)
			ModeratorID: 2, // ID модератора (зафиксирован)
		}
	}
	return systemUsers
}

// GetCreatorID - получение ID создателя
func GetCreatorID() int {
	return GetSystemUsers().CreatorID
}

// GetModeratorID - получение ID модератора
func GetModeratorID() int {
	return GetSystemUsers().ModeratorID
}
