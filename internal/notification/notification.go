package notification

import (
	"context"
	"encoding/json"
	"github.com/MatticNote/MatticNote/database"
	"github.com/google/uuid"
)

func makeNotification(
	targetUserUuid,
	fromUserUuid uuid.UUID,
	relateNoteUuid *uuid.UUID,
	notificationType string,
	metaData interface{}) error {
	var metaDataJson *[]byte = nil
	if metaData != nil {
		marshaled, err := json.Marshal(metaData)
		if err != nil {
			return err
		}
		metaDataJson = &marshaled
	}

	_, err := database.DBPool.Exec(
		context.Background(),
		"insert into notification(target_user, from_user, relate_note, type, metadata) values ($1, $2, $3, $4, $5);",
		targetUserUuid.String(),
		fromUserUuid.String(),
		relateNoteUuid,
		notificationType,
		metaDataJson,
	)
	if err != nil {
		return err
	}

	// TODO: websocketで該当ユーザーに送信

	return nil
}

func MakeFollowNotification(fromUserUuid, targetUserUuid uuid.UUID) error {
	return makeNotification(targetUserUuid, fromUserUuid, nil, "FOLLOWED", nil)
}
