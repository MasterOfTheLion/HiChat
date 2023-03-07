package dao

import (
	"HiChat/global"
	"HiChat/models"
	"errors"
	"go.uber.org/zap"
)

func FriendList(userId uint) (*[]models.UserBasic, error) {
	relation := make([]models.Relation, 0)
	if tx := global.DB.Where("owner_id = ? and type=1", userId).Find(&relation); tx.RowsAffected == 0 {
		zap.S().Info("未查询到Relation数据")
		return nil, errors.New("未查到好友关系")
	}

	userID := make([]uint, 0)
	for _, v := range relation {
		userID = append(userID, v.TargetID)
	}

	user := make([]models.UserBasic, 0)
	if tx := global.DB.Where("id in ?", userID).Find(&user); tx.RowsAffected == 0 {
		zap.S().Info("未查询到Relation好友关系")
		return nil, errors.New("未查到好友")
	}
	return &user, nil
}

func AddFriend(userID, TargetId uint) (int, error) {
	if userID == TargetId {
		return -2, errors.New("userID和TargetId相等")
	}
	targetUser, err := FindUserID(TargetId)
	if err != nil {
		return -1, errors.New("未查询用户")
	}
	if targetUser.ID == 0 {
		zap.S().Info("未查询用户")
		return -1, errors.New("未查询用户")
	}

	relation := models.Relation{}

	if tx := global.DB.Where("owner_id = ? and target_id = ? and type = 1", userID, TargetId).First(&relation); tx.RowsAffected == 1 {
		zap.S().Info("该好友存在")
		return 0, errors.New("好友已经存在")
	}

	if tx := global.DB.Where("owner_id = ? and target_id = ? and type = 1", TargetId, userID).First(&relation); tx.RowsAffected == 1 {
		zap.S().Info("该好友存在")
		return 0, errors.New("好友已经存在")
	}

	tx := global.DB.Begin()

	relation.OwnerId = userID
	relation.TargetID = targetUser.ID
	relation.Type = 1

	if t := tx.Create(&relation); t.RowsAffected == 0 {
		zap.S().Info("创建失败")

		tx.Rollback()
		return -1, errors.New("创建好友记录失败")
	}

	relation = models.Relation{}
	relation.OwnerId = TargetId
	relation.TargetID = userID
	relation.Type = 1

	if t := tx.Create(&relation); t.RowsAffected == 0 {
		zap.S().Info("创建失败")

		tx.Rollback()
		return -1, errors.New("创建好友记录失败")
	}

	tx.Commit()
	return 1, nil
}
