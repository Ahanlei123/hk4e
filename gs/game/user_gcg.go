package game

import (
	"hk4e/common/constant"
	"hk4e/gs/model"
	"hk4e/pkg/logger"
	"hk4e/protocol/cmd"
	"hk4e/protocol/proto"

	pb "google.golang.org/protobuf/proto"
)

func (g *GameManager) GCGLogin(player *model.Player) {
	// player.SceneId = 1076
	// player.Pos.X = 8.974
	// player.Pos.Y = 0
	// player.Pos.Z = 9.373
	// GCG基础信息
	g.SendMsg(cmd.GCGBasicDataNotify, player.PlayerID, player.ClientSeq, g.PacketGCGBasicDataNotify(player))
	// GCG等级挑战解锁
	g.SendMsg(cmd.GCGLevelChallengeNotify, player.PlayerID, player.ClientSeq, g.PacketGCGLevelChallengeNotify(player))
	// GCG禁止的卡牌
	g.SendMsg(cmd.GCGDSBanCardNotify, player.PlayerID, player.ClientSeq, g.PacketGCGDSBanCardNotify(player))
	// GCG解锁或拥有的内容
	g.SendMsg(cmd.GCGDSDataNotify, player.PlayerID, player.ClientSeq, g.PacketGCGDSDataNotify(player))
	// GCG酒馆挑战数据
	g.SendMsg(cmd.GCGTCTavernChallengeDataNotify, player.PlayerID, player.ClientSeq, g.PacketGCGTCTavernChallengeDataNotify(player))
}

// GCGTavernInit GCG酒馆初始化
func (g *GameManager) GCGTavernInit(player *model.Player) {
	if player.SceneId == 1076 {
		// GCG酒馆信息通知
		g.SendMsg(cmd.GCGTCTavernInfoNotify, player.PlayerID, player.ClientSeq, g.PacketGCGTCTavernInfoNotify(player))
		// GCG酒馆NPC信息通知
		g.SendMsg(cmd.GCGTavernNpcInfoNotify, player.PlayerID, player.ClientSeq, g.PacketGCGTavernNpcInfoNotify(player))
	}
}

// GCGStartChallenge GCG开始挑战
func (g *GameManager) GCGStartChallenge(player *model.Player) {
	// GCG开始游戏通知
	// gcgStartChallengeByCheckRewardRsp := &proto.GCGStartChallengeByCheckRewardRsp{
	//	ExceededItemTypeList: make([]uint32, 0, 0),
	//	LevelId:              0,
	//	ExceededItemList:     make([]uint32, 0, 0),
	//	LevelType:            proto.GCGLevelType_GCG_LEVEL_TYPE_GUIDE_GROUP,
	//	ConfigId:             7066505,
	//	Retcode:              0,
	// }
	// g.SendMsg(cmd.GCGStartChallengeByCheckRewardRsp, player.PlayerID, player.ClientSeq, gcgStartChallengeByCheckRewardRsp)

	// 创建GCG游戏
	game := GCG_MANAGER.CreateGame(30101, []*model.Player{player})

	// GCG游戏简要信息通知
	GAME_MANAGER.SendMsg(cmd.GCGGameBriefDataNotify, player.PlayerID, player.ClientSeq, g.PacketGCGGameBriefDataNotify(player, proto.GCGGameBusinessType_GCG_GAME_BUSINESS_TYPE_GUIDE_GROUP, game))

	// 玩家进入GCG界面
	g.TeleportPlayer(player, constant.EnterReasonConst.DungeonEnter, 79999, new(model.Vector), 2162)
}

// GCGAskDuelReq GCG决斗请求
func (g *GameManager) GCGAskDuelReq(player *model.Player, payloadMsg pb.Message) {
	// 获取玩家所在的游戏
	game, ok := GCG_MANAGER.gameMap[player.GCGCurGameGuid]
	if !ok {
		g.CommonRetError(cmd.GCGAskDuelRsp, player, &proto.GCGAskDuelRsp{}, proto.Retcode_RET_GCG_GAME_NOT_RUNNING)
		return
	}
	// 获取玩家的操控者对象
	gameController := game.GetControllerByUserId(player.PlayerID)
	if gameController == nil {
		g.CommonRetError(cmd.GCGAskDuelRsp, player, &proto.GCGAskDuelRsp{}, proto.Retcode_RET_GCG_NOT_IN_GCG_DUNGEON)
		return
	}

	// 更改操控者加载状态
	gameController.loadState = ControllerLoadState_AskDuel

	// 计数器+1
	game.serverSeqCounter++
	// PacketGCGAskDuelRsp
	gcgAskDuelRsp := &proto.GCGAskDuelRsp{
		Duel: &proto.GCGDuel{
			ServerSeq: game.serverSeqCounter,
			// ShowInfoList 游戏内显示双方头像名字
			ShowInfoList:              make([]*proto.GCGControllerShowInfo, 0, len(game.controllerMap)),
			ForbidFinishChallengeList: nil,
			// CardList 卡牌列表
			CardList:            make([]*proto.GCGCard, 0, 0),
			Unk3300_BIANMOPDEHO: 1, // Unk
			CostRevise: &proto.GCGCostReviseInfo{ // 暂无数据
				CanUseHandCardIdList:  nil,
				SelectOnStageCostList: nil,
				PlayCardCostList:      nil,
				AttackCostList:        nil,
				IsCanAttack:           false,
			},
			GameId: 0, // 官服是0
			// FieldList 玩家牌盒信息 卡牌显示相关
			FieldList:           make([]*proto.GCGPlayerField, 0, len(game.controllerMap)),
			Unk3300_CDCMBOKBLAK: make([]*proto.Unk3300_ADHENCIFKNI, 0, len(game.controllerMap)),
			BusinessType:        0,
			IntentionList:       nil, // empty
			// ChallengeList 可能是挑战目标
			ChallengeList: []*proto.GCGDuelChallenge{
				// TODO 暂时写死
				{
					ChallengeId:   1,
					CurProgress:   906,
					TotalProgress: 0,
				},
				{
					ChallengeId:   1,
					CurProgress:   907,
					TotalProgress: 0,
				},
				{
					ChallengeId:   1,
					CurProgress:   901,
					TotalProgress: 0,
				},
				{
					ChallengeId:   1,
					CurProgress:   903,
					TotalProgress: 0,
				},
				{
					ChallengeId:   1,
					CurProgress:   904,
					TotalProgress: 0,
				},
				{
					ChallengeId:   1,
					CurProgress:   905,
					TotalProgress: 0,
				},
				{
					ChallengeId:   1,
					CurProgress:   908,
					TotalProgress: 0,
				},
				{
					ChallengeId:   1,
					CurProgress:   909,
					TotalProgress: 0,
				},
			},
			HistoryCardList:     make([]*proto.GCGCard, 0, len(game.historyCardList)),
			Round:               game.roundInfo.roundNum,
			ControllerId:        gameController.controllerId,
			HistoryMsgPackList:  game.historyMsgPackList,
			Unk3300_JHDDNKFPINA: 0,
			// CardIdList 游戏内的所有卡牌Id
			CardIdList:          make([]uint32, 0, 0),
			Unk3300_JBBMBKGOONO: 0, // Unk
			// 阶段数据
			Phase: &proto.GCGPhase{
				PhaseType:          game.roundInfo.phaseType,
				AllowControllerMap: make(map[uint32]uint32),
			},
		},
	}
	// 玩家信息列表
	for _, controller := range game.controllerMap {
		gcgControllerShowInfo := &proto.GCGControllerShowInfo{
			ControllerId:   controller.controllerId,
			ProfilePicture: &proto.ProfilePicture{},
		}
		// 如果为玩家则更改为玩家信息
		if controller.controllerType == ControllerType_Player {
			gcgControllerShowInfo.ProfilePicture.AvatarId = player.HeadImage
			gcgControllerShowInfo.ProfilePicture.AvatarId = player.AvatarMap[player.HeadImage].Costume
		}
		gcgAskDuelRsp.Duel.ShowInfoList = append(gcgAskDuelRsp.Duel.ShowInfoList)
	}
	// 玩家牌盒信息 卡牌显示相关
	for _, controller := range game.controllerMap {
		// FieldList 玩家牌盒信息 卡牌显示相关
		playerField := &proto.GCGPlayerField{
			Unk3300_IKJMGAHCFPM: 0,
			// 卡牌图片
			ModifyZoneMap:       make(map[uint32]*proto.GCGZone, len(controller.charCardList)),
			Unk3300_GGHKFFADEAL: 0,
			Unk3300_AOPJIOHMPOF: &proto.GCGZone{
				CardList: []uint32{},
			},
			Unk3300_FDFPHNDOJML: 0,
			// 卡牌技能?
			Unk3300_IPLMHKCNDLE: &proto.GCGZone{
				CardList: []uint32{}, // 官服CardList: []uint32{5},
			},
			Unk3300_EIHOMDLENMK: &proto.GCGZone{
				CardList: []uint32{},
			},
			WaitingList:         []*proto.GCGWaitingCharacter{},
			Unk3300_PBECINKKHND: 0,
			ControllerId:        controller.controllerId,
			// 卡牌位置
			Unk3300_INDJNJJJNKL: &proto.GCGZone{
				CardList: make([]uint32, 0, len(controller.charCardList)),
			},
			Unk3300_EFNAEFBECHD: &proto.GCGZone{
				CardList: []uint32{},
			},
			IsPassed:            false,
			IntentionList:       []*proto.GCGPVEIntention{},
			DiceSideList:        []proto.GCGDiceSideType{},
			DeckCardNum:         0,
			Unk3300_GLNIFLOKBPM: 0,
		}
		for _, info := range controller.charCardList {
			playerField.ModifyZoneMap[info.guid] = &proto.GCGZone{CardList: []uint32{}}
			playerField.Unk3300_INDJNJJJNKL.CardList = append(playerField.Unk3300_INDJNJJJNKL.CardList, info.guid)
		}
		// 添加完所有卡牌的位置之类的信息后添加这个牌盒
		gcgAskDuelRsp.Duel.FieldList = append(gcgAskDuelRsp.Duel.FieldList, playerField)
	}
	// 历史卡牌信息
	for _, cardInfo := range game.historyCardList {
		gcgAskDuelRsp.Duel.HistoryCardList = append(gcgAskDuelRsp.Duel.HistoryCardList, cardInfo.ToProto())
	}
	// 卡牌信息
	for _, controller := range game.controllerMap {
		// 角色牌以及手牌都要
		for _, cardInfo := range controller.charCardList {
			gcgAskDuelRsp.Duel.CardList = append(gcgAskDuelRsp.Duel.CardList, cardInfo.ToProto())
			gcgAskDuelRsp.Duel.CardIdList = append(gcgAskDuelRsp.Duel.CardIdList, cardInfo.cardId)
		}
		for _, cardInfo := range controller.handCardList {
			gcgAskDuelRsp.Duel.CardList = append(gcgAskDuelRsp.Duel.CardList, cardInfo.ToProto())
			gcgAskDuelRsp.Duel.CardIdList = append(gcgAskDuelRsp.Duel.CardIdList, cardInfo.cardId)
		}
	}
	// 阶段信息
	for _, controller := range game.controllerMap {
		// 操控者是否允许操作
		gcgAskDuelRsp.Duel.Phase.AllowControllerMap[controller.controllerId] = controller.allow
	}
	// Unk3300_CDCMBOKBLAK 你问我这是啥? 我也不知道
	for _, controller := range game.controllerMap {
		gcgAskDuelRsp.Duel.Unk3300_CDCMBOKBLAK = append(gcgAskDuelRsp.Duel.Unk3300_CDCMBOKBLAK, &proto.Unk3300_ADHENCIFKNI{
			BeginTime:    0,
			TimeStamp:    0,
			ControllerId: controller.controllerId,
		})
	}

	GAME_MANAGER.SendMsg(cmd.GCGAskDuelRsp, player.PlayerID, player.ClientSeq, gcgAskDuelRsp)
}

// GCGInitFinishReq GCG初始化完成请求
func (g *GameManager) GCGInitFinishReq(player *model.Player, payloadMsg pb.Message) {
	// 获取玩家所在的游戏
	game, ok := GCG_MANAGER.gameMap[player.GCGCurGameGuid]
	if !ok {
		g.CommonRetError(cmd.GCGInitFinishRsp, player, &proto.GCGInitFinishRsp{}, proto.Retcode_RET_GCG_GAME_NOT_RUNNING)
		return
	}
	// 获取玩家的操控者对象
	gameController := game.GetControllerByUserId(player.PlayerID)
	if gameController == nil {
		g.CommonRetError(cmd.GCGInitFinishRsp, player, &proto.GCGInitFinishRsp{}, proto.Retcode_RET_GCG_NOT_IN_GCG_DUNGEON)
		return
	}

	// 更改操控者加载状态
	gameController.loadState = ControllerLoadState_InitFinish

	GAME_MANAGER.SendMsg(cmd.GCGInitFinishRsp, player.PlayerID, player.ClientSeq, &proto.GCGInitFinishRsp{})

	// 检查所有玩家是否已加载完毕
	game.CheckAllInitFinish()
}

// GCGOperationReq GCG游戏客户端操作请求
func (g *GameManager) GCGOperationReq(player *model.Player, payloadMsg pb.Message) {
	req := payloadMsg.(*proto.GCGOperationReq)

	// 获取玩家所在的游戏
	game, ok := GCG_MANAGER.gameMap[player.GCGCurGameGuid]
	if !ok {
		g.CommonRetError(cmd.GCGOperationRsp, player, &proto.GCGOperationRsp{}, proto.Retcode_RET_GCG_GAME_NOT_RUNNING)
		return
	}
	// 获取玩家的操控者对象
	gameController := game.GetControllerByUserId(player.PlayerID)
	if gameController == nil {
		g.CommonRetError(cmd.GCGOperationRsp, player, &proto.GCGOperationRsp{}, proto.Retcode_RET_GCG_NOT_IN_GCG_DUNGEON)
		return
	}

	switch req.Op.Op.(type) {
	case *proto.GCGOperation_OpSelectOnStage:
		// 选择角色卡牌
		op := req.Op.GetOpSelectOnStage()
		// 操作者是否拥有该卡牌
		cardInfo := gameController.GetCharCardByGuid(op.CardGuid)
		if cardInfo == nil {
			GAME_MANAGER.CommonRetError(cmd.GCGOperationRsp, player, &proto.GCGOperationRsp{}, proto.Retcode_RET_GCG_SELECT_HAND_CARD_GUID_ERROR)
			return
		}
		// 操控者选择角色牌
		game.ControllerSelectChar(gameController, cardInfo, op.CostDiceIndexList)
	case *proto.GCGOperation_OpReroll:
		// 确认骰子重投
		op := req.Op.GetOpReroll()
		diceSideList, ok := game.roundInfo.diceSideMap[gameController.controllerId]
		if !ok {
			g.CommonRetError(cmd.GCGOperationRsp, player, &proto.GCGOperationRsp{}, proto.Retcode_RET_GCG_DICE_INDEX_INVALID)
			return
		}
		// 判断骰子索引是否有效
		for _, diceIndex := range op.DiceIndexList {
			if diceIndex > uint32(len(diceSideList)) {
				g.CommonRetError(cmd.GCGOperationRsp, player, &proto.GCGOperationRsp{}, proto.Retcode_RET_GCG_DICE_INDEX_INVALID)
				return
			}
		}
		// 操控者确认重投骰子
		game.ControllerReRollDice(gameController, op.DiceIndexList)
	case *proto.GCGOperation_OpAttack:
		// 角色使用技能
		op := req.Op.GetOpAttack()
		diceSideList, ok := game.roundInfo.diceSideMap[gameController.controllerId]
		if !ok {
			g.CommonRetError(cmd.GCGOperationRsp, player, &proto.GCGOperationRsp{}, proto.Retcode_RET_GCG_DICE_INDEX_INVALID)
			return
		}
		// 判断骰子索引是否有效
		for _, diceIndex := range op.CostDiceIndexList {
			if diceIndex > uint32(len(diceSideList)) {
				g.CommonRetError(cmd.GCGOperationRsp, player, &proto.GCGOperationRsp{}, proto.Retcode_RET_GCG_DICE_INDEX_INVALID)
				return
			}
		}
		// 操控者使用技能
		game.ControllerUseSkill(gameController, op.SkillId, op.CostDiceIndexList)
	default:
		logger.Error("gcg op is not handle, op: %T", req.Op.Op)
		return
	}
	// PacketGCGOperationRsp
	gcgOperationRsp := &proto.GCGOperationRsp{
		OpSeq: req.OpSeq,
	}
	GAME_MANAGER.SendMsg(cmd.GCGOperationRsp, player.PlayerID, player.ClientSeq, gcgOperationRsp)
}

// PacketGCGSkillPreviewNotify GCG游戏技能预览通知
func (g *GameManager) PacketGCGSkillPreviewNotify(controller *GCGController) *proto.GCGSkillPreviewNotify {
	selectedCharCard := controller.GetSelectedCharCard()
	// 确保玩家选择了角色牌
	if selectedCharCard == nil {
		logger.Error("selected char card is nil, cardGuid: %v", controller.selectedCharCardGuid)
		return new(proto.GCGSkillPreviewNotify)
	}
	// PacketGCGSkillPreviewNotify
	gcgSkillPreviewNotify := &proto.GCGSkillPreviewNotify{
		ControllerId: controller.controllerId,
		// 当前角色牌拥有的技能信息
		SkillPreviewList: make([]*proto.GCGSkillPreviewInfo, 0, len(selectedCharCard.skillList)),
		// 切换到其他角色牌的所需消耗信息
		ChangeOnstagePreviewList: make([]*proto.GCGChangeOnstageInfo, 0, 2), // 暂时写死
		PlayCardList:             make([]*proto.GCGSkillPreviewPlayCardInfo, 0, 0),
		OnstageCardGuid:          selectedCharCard.guid, // 当前被选择的角色牌guid
	}
	// SkillPreviewList
	for _, skillInfo := range selectedCharCard.skillList {
		gcgSkillPreviewInfo := &proto.GCGSkillPreviewInfo{
			ChangeOnstageCharacterList: nil,
			Unk3300_DAJFJEDNLKK:        nil,
			SkillId:                    skillInfo.skillId,
			// 技能造成的血量预览信息
			HpInfoMap:           make(map[uint32]*proto.GCGSkillPreviewHpInfo, 1),
			Unk3300_AGNONGELFGC: nil,
			ExtraInfo:           nil,
			ReactionInfoMap:     nil,
			// 技能对自身改变预览信息
			CardTokenChangeMap: make(map[uint32]*proto.GCGSkillPreviewTokenChangeInfo, 1),
		}
		// HpInfoMap
		// 暂时不知道3代表什么意思
		gcgSkillPreviewInfo.HpInfoMap[3] = &proto.GCGSkillPreviewHpInfo{
			ChangeType:    proto.GCGSkillHpChangeType_GCG_SKILL_HP_CHANGE_TYPE_DAMAGE,
			HpChangeValue: skillInfo.damage,
		}
		// CardTokenChangeMap
		// 暂时不知道1代表什么意思
		gcgSkillPreviewInfo.CardTokenChangeMap[1] = &proto.GCGSkillPreviewTokenChangeInfo{
			TokenChangeList: []*proto.GCGSkillPreviewTokenInfo{
				{
					// Token类型
					TokenType:           constant.GCGTokenConst.TOKEN_CUR_ELEM,
					Unk3300_MMIKPPJMHAD: 0,
					// 更改为的值
					Unk3300_IKICJMEFEON: selectedCharCard.tokenMap[constant.GCGTokenConst.TOKEN_CUR_ELEM] + 1,
				},
			},
		}
		gcgSkillPreviewNotify.SkillPreviewList = append(gcgSkillPreviewNotify.SkillPreviewList, gcgSkillPreviewInfo)
	}
	// ChangeOnstagePreviewList
	for _, cardInfo := range controller.charCardList {
		// 排除当前已选中的角色卡
		if cardInfo.guid == selectedCharCard.guid {
			continue
		}
		gcgChangeOnstageInfo := &proto.GCGChangeOnstageInfo{
			IsQuick:  false, // 是否为快速行动
			CardGuid: cardInfo.guid,
			// 切换角色预览
			ChangeOnstagePreviewInfo: &proto.GCGSkillPreviewInfo{},
		}
		gcgSkillPreviewNotify.ChangeOnstagePreviewList = append(gcgSkillPreviewNotify.ChangeOnstagePreviewList, gcgChangeOnstageInfo)
	}
	return gcgSkillPreviewNotify
}

// SendGCGMessagePackNotify 发送GCG游戏消息包通知
func (g *GameManager) SendGCGMessagePackNotify(controller *GCGController, serverSeq uint32, msgPackList []*proto.GCGMessagePack) {
	// 确保加载完成
	if controller.loadState != ControllerLoadState_InitFinish {
		return
	}
	// PacketGCGMessagePackNotify
	gcgMessagePackNotify := &proto.GCGMessagePackNotify{
		ServerSeq:   serverSeq,
		MsgPackList: msgPackList,
	}
	// 根据操控者的类型发送消息包
	switch controller.controllerType {
	case ControllerType_Player:
		GAME_MANAGER.SendMsg(cmd.GCGMessagePackNotify, controller.player.PlayerID, controller.player.ClientSeq, gcgMessagePackNotify)
	case ControllerType_AI:
		controller.ai.ReceiveGCGMessagePackNotify(gcgMessagePackNotify)
	default:
		logger.Error("controller type error, %v", controller.controllerType)
		return
	}
}

// PacketGCGGameBriefDataNotify GCG游戏简要数据通知
func (g *GameManager) PacketGCGGameBriefDataNotify(player *model.Player, businessType proto.GCGGameBusinessType, game *GCGGame) *proto.GCGGameBriefDataNotify {
	gcgGameBriefDataNotify := &proto.GCGGameBriefDataNotify{
		GcgBriefData: &proto.GCGGameBriefData{
			BusinessType:    businessType,
			PlatformType:    uint32(proto.PlatformType_PLATFORM_TYPE_PC), // TODO 根据玩家设备修改
			GameId:          game.gameId,
			PlayerBriefList: make([]*proto.GCGPlayerBriefData, 0, len(game.controllerMap)),
		},
		IsNewGame: true, // TODO 根据游戏修改
	}
	for _, controller := range game.controllerMap {
		gcgPlayerBriefData := &proto.GCGPlayerBriefData{
			ControllerId:   controller.controllerId,
			ProfilePicture: new(proto.ProfilePicture),
			CardIdList:     make([]uint32, 0, len(controller.charCardList)), // 这里展示给玩家的是角色牌
		}
		// 玩家信息
		if controller.player != nil {
			gcgPlayerBriefData.Uid = player.PlayerID
			gcgPlayerBriefData.ProfilePicture.AvatarId = player.TeamConfig.GetActiveAvatarId()
			gcgPlayerBriefData.ProfilePicture.CostumeId = player.AvatarMap[player.TeamConfig.GetActiveAvatarId()].Costume
			gcgPlayerBriefData.NickName = player.NickName
		}
		gcgGameBriefDataNotify.GcgBriefData.PlayerBriefList = append(gcgGameBriefDataNotify.GcgBriefData.PlayerBriefList)
	}
	return gcgGameBriefDataNotify
}

// PacketGCGTavernNpcInfoNotify GCG酒馆NPC信息通知
func (g *GameManager) PacketGCGTavernNpcInfoNotify(player *model.Player) *proto.GCGTavernNpcInfoNotify {
	gcgTavernNpcInfoNotify := &proto.GCGTavernNpcInfoNotify{
		Unk3300_FKAKHMMIEBC: make([]*proto.GCGTavernNpcInfo, 0, 0),
		Unk3300_BAMLNENDLCM: make([]*proto.GCGTavernNpcInfo, 0, 0),
		CharacterNpc: &proto.GCGTavernNpcInfo{
			Id:           0,
			ScenePointId: 0,
			LevelId:      0,
		},
	}
	return gcgTavernNpcInfoNotify
}

// PacketGCGTCTavernInfoNotify GCG酒馆信息通知
func (g *GameManager) PacketGCGTCTavernInfoNotify(player *model.Player) *proto.GCGTCTavernInfoNotify {
	gcgTCTavernInfoNotify := &proto.GCGTCTavernInfoNotify{
		LevelId:             0,
		Unk3300_IMFJBNFMCHM: false,
		Unk3300_MBGMHBNBKBK: false,
		PointId:             0,
		ElementType:         8,
		AvatarId:            10000007,
		CharacterId:         0,
	}
	return gcgTCTavernInfoNotify
}

// PacketGCGTCTavernChallengeDataNotify GCG酒馆挑战数据
func (g *GameManager) PacketGCGTCTavernChallengeDataNotify(player *model.Player) *proto.GCGTCTavernChallengeDataNotify {
	gcgTCTavernChallengeDataNotify := &proto.GCGTCTavernChallengeDataNotify{
		TavernChallengeList: make([]*proto.GCGTCTavernChallengeData, 0, 0),
	}
	for _, challenge := range player.GCGInfo.TavernChallengeMap {
		gcgTCTavernChallengeData := &proto.GCGTCTavernChallengeData{
			UnlockLevelIdList: challenge.UnlockLevelIdList,
			CharacterId:       challenge.CharacterId,
		}
		gcgTCTavernChallengeDataNotify.TavernChallengeList = append(gcgTCTavernChallengeDataNotify.TavernChallengeList, gcgTCTavernChallengeData)
	}
	return gcgTCTavernChallengeDataNotify
}

// PacketGCGBasicDataNotify GCG基础数据通知
func (g *GameManager) PacketGCGBasicDataNotify(player *model.Player) *proto.GCGBasicDataNotify {
	gcgBasicDataNotify := &proto.GCGBasicDataNotify{
		Level:                player.GCGInfo.Level,
		Exp:                  player.GCGInfo.Exp,
		LevelRewardTakenList: make([]uint32, 0, 0),
	}
	return gcgBasicDataNotify
}

// PacketGCGLevelChallengeNotify GCG等级挑战通知
func (g *GameManager) PacketGCGLevelChallengeNotify(player *model.Player) *proto.GCGLevelChallengeNotify {
	gcgLevelChallengeNotify := &proto.GCGLevelChallengeNotify{
		UnlockBossChallengeList:  make([]*proto.GCGBossChallengeData, 0, 0),
		UnlockWorldChallengeList: player.GCGInfo.UnlockWorldChallengeList,
		LevelList:                make([]*proto.GCGLevelData, 0, 0),
	}
	// Boss挑战信息
	for _, challenge := range player.GCGInfo.UnlockBossChallengeMap {
		gcgBossChallengeData := &proto.GCGBossChallengeData{
			UnlockLevelIdList: challenge.UnlockLevelIdList,
			Id:                challenge.Id,
		}
		gcgLevelChallengeNotify.UnlockBossChallengeList = append(gcgLevelChallengeNotify.UnlockBossChallengeList, gcgBossChallengeData)
	}
	// 等级挑战信息
	for _, challenge := range player.GCGInfo.LevelChallengeMap {
		gcgLevelData := &proto.GCGLevelData{
			FinishedChallengeIdList: challenge.FinishedChallengeIdList,
			LevelId:                 challenge.LevelId,
		}
		gcgLevelChallengeNotify.LevelList = append(gcgLevelChallengeNotify.LevelList, gcgLevelData)
	}
	return gcgLevelChallengeNotify
}

// PacketGCGDSBanCardNotify GCG禁止的卡牌通知
func (g *GameManager) PacketGCGDSBanCardNotify(player *model.Player) *proto.GCGDSBanCardNotify {
	gcgDSBanCardNotify := &proto.GCGDSBanCardNotify{
		CardList: player.GCGInfo.BanCardList,
	}
	return gcgDSBanCardNotify
}

// PacketGCGDSDataNotify GCG数据通知
func (g *GameManager) PacketGCGDSDataNotify(player *model.Player) *proto.GCGDSDataNotify {
	gcgDSDataNotify := &proto.GCGDSDataNotify{
		CurDeckId:            player.GCGInfo.CurDeckId,
		DeckList:             make([]*proto.GCGDSDeckData, 0, len(player.GCGInfo.DeckList)),
		UnlockCardBackIdList: player.GCGInfo.UnlockCardBackIdList,
		CardList:             make([]*proto.GCGDSCardData, 0, len(player.GCGInfo.CardList)),
		UnlockFieldIdList:    player.GCGInfo.UnlockFieldIdList,
		UnlockDeckIdList:     player.GCGInfo.UnlockDeckIdList,
	}
	// 卡组列表
	for i, deck := range player.GCGInfo.DeckList {
		gcgDSDeckData := &proto.GCGDSDeckData{
			CreateTime:        uint32(deck.CreateTime),
			FieldId:           deck.FieldId,
			CardBackId:        deck.CardBackId,
			CardList:          deck.CardList,
			CharacterCardList: deck.CharacterCardList,
			Id:                uint32(i),
			Name:              deck.Name,
			IsValid:           true, // TODO 校验卡组是否有效
		}
		gcgDSDataNotify.DeckList = append(gcgDSDataNotify.DeckList, gcgDSDeckData)
	}
	// 卡牌列表
	for _, card := range player.GCGInfo.CardList {
		gcgDSCardData := &proto.GCGDSCardData{
			Num:                           card.Num,
			FaceType:                      card.FaceType,
			CardId:                        card.CardId,
			ProficiencyRewardTakenIdxList: card.ProficiencyRewardTakenIdxList,
			UnlockFaceTypeList:            card.UnlockFaceTypeList,
			Proficiency:                   card.Proficiency,
		}
		gcgDSDataNotify.CardList = append(gcgDSDataNotify.CardList, gcgDSCardData)
	}
	return gcgDSDataNotify
}
