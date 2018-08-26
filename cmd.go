package padchat

import (
	"errors"
	"regexp"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/json-iterator/go"
)

func (bot *Bot) sendCommand(cmd string, data interface{}) CommandResp {
	c := make(chan CommandResp)
	id := uuid.New().String()
	bot.retProcMap.Store(id, func(resp CommandResp) {
		c <- resp
	})
	bot.ws.WriteJSON(WSReq{Type: "user", Cmd: cmd, CmdID: id, Data: data})
	select {
	case d := <-c:
		return d
	case <-time.After(bot.reqTimeout):
		bot.retProcMap.Delete(id)
		close(c)
		return CommandResp{Success: false, Msg: "timeout"}
	}
}

// Init 执行初始化, 必须在登录前调用
func (bot *Bot) Init() CommandResp {
	return bot.sendCommand("init", nil)
}

// Close 关闭微信实例（不退出登陆）
func (bot *Bot) Close() CommandResp {
	return bot.sendCommand("close", nil)
}

// QRLogin 二维码登录
func (bot *Bot) QRLogin() CommandResp {
	return bot.sendCommand("login", LoginReq{LoginType: "qrcode"})
}

// RequestLogin 二次登陆, 手机端会弹出确认框, 点击后登陆, 不容易封号
func (bot *Bot) RequestLogin(wxData, token string) CommandResp {
	return bot.sendCommand("login", LoginReq{
		LoginType: "request",
		WXData:    wxData,
		Token:     token,
	})
}

// TokenLogin 断线重连, 用于短时间使用 `wxData` 和 `token` 再次登录
// `token`有效期很短, 如果登陆失败, 建议使用二次登陆方式
func (bot *Bot) TokenLogin(wxData, token string) CommandResp {
	return bot.sendCommand("login", LoginReq{
		LoginType: "token",
		WXData:    wxData,
		Token:     token,
	})
}

// UserLogin 账号密码登录
func (bot *Bot) UserLogin(wxData, username, password string) CommandResp {
	return bot.sendCommand("login", LoginReq{
		LoginType: "user",
		WXData:    wxData,
		UserName:  username,
		Password:  password,
	})
}

// PhoneLogin 手机验证码登录
func (bot *Bot) PhoneLogin(wxData, phone, code string) CommandResp {
	return bot.sendCommand("login", LoginReq{
		LoginType: "phone",
		WXData:    wxData,
		Phone:     phone,
		Code:      code,
	})
}

// GetWXData 获取设备62数据
func (bot *Bot) GetWXData() (string, error) {
	data := &struct {
		WXData string `json:"wx_data"`
	}{}
	resp := bot.sendCommand("getWxData", nil)
	if !resp.Success {
		return "", errors.New(resp.Msg)
	}
	err := jsoniter.Unmarshal(resp.Data, data)
	if err != nil {
		return "", err
	}
	return data.WXData, nil
}

// GetLoginToken 获取二次登陆数据
func (bot *Bot) GetLoginToken() (*LoginTokenResp, error) {
	data := &LoginTokenResp{}
	resp := bot.sendCommand("getLoginToken", nil)
	if !resp.Success {
		return nil, errors.New(resp.Msg)
	}
	err := jsoniter.Unmarshal(resp.Data, data)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// GetMyInfo 获取Bot微信号信息
func (bot *Bot) GetMyInfo() (*MyInfoResp, error) {
	data := &MyInfoResp{}
	resp := bot.sendCommand("getMyInfo", nil)
	if !resp.Success {
		return nil, errors.New(resp.Msg)
	}
	err := jsoniter.Unmarshal(resp.Data, data)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// SyncMsg 同步消息, 使用此接口手动触发同步消息, 一般用于刚登陆后调用, 可立即开始同步消息.
// 否则会在有新消息时才开始同步消息.
func (bot *Bot) SyncMsg() CommandResp {
	return bot.sendCommand("syncMsg", nil)
}

// Logout 退出登录
func (bot *Bot) Logout() CommandResp {
	return bot.sendCommand("logout", nil)
}

// SyncContact 同步通讯录
func (bot *Bot) SyncContact() CommandResp {
	return bot.sendCommand("syncContact", nil)
}

// SendMsg 发送文字信息
func (bot *Bot) SendMsg(req *SendMsgReq) (*SendMsgResp, error) {
	data := &SendMsgResp{}
	MkAtContent(req)
	resp := bot.sendCommand("sendMsg", req)
	if !resp.Success {
		return nil, errors.New(resp.Msg)
	}
	err := jsoniter.Unmarshal(resp.Data, data)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// MkAtContent prefix 'at' symbol for req content if necessary
func MkAtContent(req *SendMsgReq) {
	if len(req.AtList) > 0 {
		arr := regexp.MustCompile(`@`).
			FindAllString(req.Content, -1)
		sub := len(req.AtList) - len(arr)
		if sub > 0 {
			req.Content = strings.Repeat("@", sub) +
				"\n" + req.Content
		}
	}
}

// ShareCard 分享名片
func (bot *Bot) ShareCard(toUserName, content, userId string) (*SendMsgResp, error) {
	data := &SendMsgResp{}
	resp := bot.sendCommand("shareCard", struct {
		ToUserName string `json:"toUserName"`
		Content    string `json:"content"`
		UserID     string `json:"userId"`
	}{
		ToUserName: toUserName,
		Content:    content,
		UserID:     userId,
	})
	if !resp.Success {
		return nil, errors.New(resp.Msg)
	}
	err := jsoniter.Unmarshal(resp.Data, data)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// SendImage 发送图片消息, file 为图片 base64 数据
func (bot *Bot) SendImage(req SendMsgReq) (*SendMsgResp, error) {
	data := &SendMsgResp{}
	resp := bot.sendCommand("sendImage", req)
	if !resp.Success {
		return nil, errors.New(resp.Msg)
	}
	err := jsoniter.Unmarshal(resp.Data, data)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// GetRoomMembers 获取群成员信息
func (bot *Bot) GetRoomMembers(groupID string) (*ChatroomInfo, error) {
	resp := bot.sendCommand("getRoomMembers", struct {
		GroupID string `json:"groupId"`
	}{GroupID: groupID})
	if !resp.Success {
		return nil, errors.New(resp.Msg)
	}
	chatroomInfo := &ChatroomInfo{}
	err := jsoniter.Unmarshal(resp.Data, chatroomInfo)
	if err != nil {
		return nil, err
	}
	var ms []ChatMemberInfo
	err = jsoniter.Unmarshal([]byte(chatroomInfo.Member), &ms)
	if err != nil {
		return nil, err
	}
	chatroomInfo.Members = ms
	return chatroomInfo, nil
}

// GetContact 获取用户/群信息
func (bot *Bot) GetContact(userID string) (*Contact, error) {
	resp := bot.sendCommand("getContact", struct {
		UserID string `json:"userId"`
	}{UserID: userID})
	if !resp.Success {
		return nil, errors.New(resp.Msg)
	}
	contact := &Contact{}
	err := jsoniter.Unmarshal(resp.Data, contact)
	if err != nil {
		return nil, err
	}
	return contact, nil
}

// GetMsgImage 获取消息原始图片
// mType = 3
func (bot *Bot) GetMsgImage(rawMsgData Msg) (*MsgImageResp, error) {
	rawMsgData.Data = ""
	resp := bot.sendCommand("getMsgImage", struct {
		RawMsgData Msg `json:"rawMsgData"`
	}{RawMsgData: rawMsgData})
	if !resp.Success {
		return nil, errors.New(resp.Msg)
	}
	imgData := &MsgImageResp{}
	err := jsoniter.Unmarshal(resp.Data, imgData)
	if err != nil {
		return nil, err
	}
	return imgData, nil
}

// GetMsgVideo 获取消息原始视频
// mType = 43
func (bot *Bot) GetMsgVideo(rawMsgData Msg) (*MsgVideoResp, error) {
	rawMsgData.Data = ""
	resp := bot.sendCommand("getMsgVideo", struct {
		RawMsgData Msg `json:"rawMsgData"`
	}{RawMsgData: rawMsgData})
	if !resp.Success {
		return nil, errors.New(resp.Msg)
	}
	videoData := &MsgVideoResp{}
	err := jsoniter.Unmarshal(resp.Data, videoData)
	if err != nil {
		return nil, err
	}
	return videoData, nil
}

// GetMsgVoice 获取消息语音数据
// mType = 34
func (bot *Bot) GetMsgVoice(rawMsgData Msg) (*MsgVoiceResp, error) {
	rawMsgData.Data = ""
	resp := bot.sendCommand("getMsgVoice", struct {
		RawMsgData Msg `json:"rawMsgData"`
	}{RawMsgData: rawMsgData})
	if !resp.Success {
		return nil, errors.New(resp.Msg)
	}
	voiceData := &MsgVoiceResp{}
	err := jsoniter.Unmarshal(resp.Data, voiceData)
	if err != nil {
		return nil, err
	}
	return voiceData, nil
}

// CreateRoom 创建群
func (bot *Bot) CreateRoom(userList []string) (*CreateRoomResp, error) {
	resp := bot.sendCommand("createRoom", struct {
		UserList []string `json:"userList"`
	}{UserList: userList})
	if !resp.Success {
		return nil, errors.New(resp.Msg)
	}
	data := &CreateRoomResp{}
	err := jsoniter.Unmarshal(resp.Data, data)
	if err != nil {
		return nil, err
	}
	if data.UserName == "" {
		return nil, errors.New(data.Message)
	}
	return data, nil
}

// AddRoomMember 添加群成员
func (bot *Bot) AddRoomMember(groupID, userID string) (*MsgAndStatus, error) {
	resp := bot.sendCommand("addRoomMember", struct {
		GroupID string `json:"groupId"`
		UserID  string `json:"userId"`
	}{
		GroupID: groupID,
		UserID:  userID,
	})
	if !resp.Success {
		return nil, errors.New(resp.Msg)
	}
	data := &MsgAndStatus{}
	err := jsoniter.Unmarshal(resp.Data, data)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// InviteRoomMember 邀请群成员, 会给对方发送一条邀请消息, 无法判断对方是否真的接收到
func (bot *Bot) InviteRoomMember(groupID, userID string) (*MsgAndStatus, error) {
	resp := bot.sendCommand("inviteRoomMember", struct {
		GroupID string `json:"groupId"`
		UserID  string `json:"userId"`
	}{
		GroupID: groupID,
		UserID:  userID,
	})
	if !resp.Success {
		return nil, errors.New(resp.Msg)
	}
	data := &MsgAndStatus{}
	err := jsoniter.Unmarshal(resp.Data, data)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// DeleteRoomMember 删除群成员
func (bot *Bot) DeleteRoomMember(groupID, userID string) (*MsgAndStatus, error) {
	resp := bot.sendCommand("deleteRoomMember", struct {
		GroupID string `json:"groupId"`
		UserID  string `json:"userId"`
	}{
		GroupID: groupID,
		UserID:  userID,
	})
	if !resp.Success {
		return nil, errors.New(resp.Msg)
	}
	data := &MsgAndStatus{}
	err := jsoniter.Unmarshal(resp.Data, data)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// SetRoomAnnouncement 设置群公告
func (bot *Bot) SetRoomAnnouncement(groupID, content string) (*MsgAndStatus, error) {
	resp := bot.sendCommand("setRoomAnnouncement", struct {
		GroupID string `json:"groupId"`
		Content string `json:"content"`
	}{
		GroupID: groupID,
		Content: content,
	})
	if !resp.Success {
		return nil, errors.New(resp.Msg)
	}
	data := &MsgAndStatus{}
	err := jsoniter.Unmarshal(resp.Data, data)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// SetRoomName 设置群名称
func (bot *Bot) SetRoomName(groupID, content string) (*MsgAndStatus, error) {
	resp := bot.sendCommand("setRoomName", struct {
		GroupID string `json:"groupId"`
		Content string `json:"content"`
	}{
		GroupID: groupID,
		Content: content,
	})
	if !resp.Success {
		return nil, errors.New(resp.Msg)
	}
	data := &MsgAndStatus{}
	err := jsoniter.Unmarshal(resp.Data, data)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// QuitRoom 退出群
func (bot *Bot) QuitRoom(groupID string) (*MsgAndStatus, error) {
	resp := bot.sendCommand("quitRoom", struct {
		GroupID string `json:"groupId"`
	}{
		GroupID: groupID,
	})
	if !resp.Success {
		return nil, errors.New(resp.Msg)
	}
	data := &MsgAndStatus{}
	err := jsoniter.Unmarshal(resp.Data, data)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// GetRoomQRCode 获取微信群二维码
func (bot *Bot) GetRoomQRCode(groupID string) (*QRCodeResp, error) {
	resp := bot.sendCommand("getRoomQrcode", struct {
		GroupID string `json:"groupId"`
		Style   int    `json:"style"`
	}{
		GroupID: groupID,
	})
	if !resp.Success {
		return nil, errors.New(resp.Msg)
	}
	data := &QRCodeResp{}
	err := jsoniter.Unmarshal(resp.Data, data)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// SearchContact 搜索用户
func (bot *Bot) SearchContact(userID string) (*Contact, error) {
	resp := bot.sendCommand("searchContact", struct {
		UserID string `json:"userId"`
	}{
		UserID: userID,
	})
	if !resp.Success {
		return nil, errors.New(resp.Msg)
	}
	data := &Contact{}
	err := jsoniter.Unmarshal(resp.Data, data)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// DeleteContact 删除好友
func (bot *Bot) DeleteContact(userID string) (*MsgAndStatus, error) {
	resp := bot.sendCommand("deleteContact", struct {
		UserID string `json:"userId"`
	}{
		UserID: userID,
	})
	if !resp.Success {
		return nil, errors.New(resp.Msg)
	}
	data := &MsgAndStatus{}
	err := jsoniter.Unmarshal(resp.Data, data)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// GetUserQRCode 获取用户二维码, 仅限获取自己的二维码, 无法获取其他人的二维码
func (bot *Bot) GetUserQRCode(userID string, style int) (*QRCodeResp, error) {
	resp := bot.sendCommand("getRoomQrcode", struct {
		UserID string `json:"userId"`
		Style  int    `json:"style"`
	}{
		UserID: userID,
		Style:  style,
	})
	if !resp.Success {
		return nil, errors.New(resp.Msg)
	}
	data := &QRCodeResp{}
	err := jsoniter.Unmarshal(resp.Data, data)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// AcceptUser 通过好友验证
func (bot *Bot) AcceptUser(stranger, ticket string) (*MsgAndStatus, error) {
	resp := bot.sendCommand("acceptUser", struct {
		Stranger string `json:"stranger"`
		Ticket   string `json:"ticket"`
	}{
		Stranger: stranger,
		Ticket:   ticket,
	})
	if !resp.Success {
		return nil, errors.New(resp.Msg)
	}
	data := &MsgAndStatus{}
	err := jsoniter.Unmarshal(resp.Data, data)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// AddContact 添加好友
// 0: 通过微信号搜索
// 1: 搜索QQ号
// 3: 通过微信号搜索
// 4: 通过QQ好友添加
// 8: 通过群聊
// 12: 来自QQ好友
// 14: 通过群聊
// 15: 通过搜索手机号
// 17: 通过名片分享
// 22: 通过摇一摇打招呼方式
// 25: 通过漂流瓶
// 30: 通过二维码方式
func (bot *Bot) AddContact(stranger, ticket, content string, Type int) (*MsgAndStatus, error) {
	resp := bot.sendCommand("addContact", struct {
		Stranger string `json:"stranger"`
		Ticket   string `json:"ticket"`
		Type     int    `json:"type"`
		Content  string `json:"content"`
	}{
		Stranger: stranger,
		Ticket:   ticket,
		Type:     Type,
		Content:  content,
	})
	if !resp.Success {
		return nil, errors.New(resp.Msg)
	}
	data := &MsgAndStatus{}
	err := jsoniter.Unmarshal(resp.Data, data)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// SayHello 打招呼,如果已经是好友, 会收到由系统自动发送, 来自对方的一条文本信息
// "xx已通过你的朋友验证请求，现在可以开始聊天了"
func (bot *Bot) SayHello(stranger, ticket, content string) (*MsgAndStatus, error) {
	resp := bot.sendCommand("sayHello", struct {
		Stranger string `json:"stranger"`
		Ticket   string `json:"ticket"`
		Content  string `json:"content"`
	}{
		Stranger: stranger,
		Ticket:   ticket,
		Content:  content,
	})
	if !resp.Success {
		return nil, errors.New(resp.Msg)
	}
	data := &MsgAndStatus{}
	err := jsoniter.Unmarshal(resp.Data, data)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// SetRemark 设置备注
func (bot *Bot) SetRemark(userID, remark string) (*MsgAndStatus, error) {
	resp := bot.sendCommand("setRemark", struct {
		UserID string `json:"userId"`
		Remark string `json:"remark"`
	}{
		UserID: userID,
		Remark: remark,
	})
	if !resp.Success {
		return nil, errors.New(resp.Msg)
	}
	data := &MsgAndStatus{}
	err := jsoniter.Unmarshal(resp.Data, data)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// SetHeadImg 设置头像
func (bot *Bot) SetHeadImg(file string) (*ImgResp, error) {
	resp := bot.sendCommand("setHeadImg", struct {
		File string `json:"file"`
	}{
		File: file,
	})
	if !resp.Success {
		return nil, errors.New(resp.Msg)
	}
	data := &ImgResp{}
	err := jsoniter.Unmarshal(resp.Data, data)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// SNSUpload 上传图片到朋友圈
// 此接口只能上传图片，并不会将图片发到朋友圈中
func (bot *Bot) SNSUpload(file string) (*SNSUploadResp, error) {
	resp := bot.sendCommand("snsUpload", struct {
		File string `json:"file"`
	}{
		File: file,
	})
	if !resp.Success {
		return nil, errors.New(resp.Msg)
	}
	data := &SNSUploadResp{}
	err := jsoniter.Unmarshal(resp.Data, data)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// SNSObjectOperation 操作朋友圈
// type - 操作类型，1为删除朋友圈，4为删除评论，5为取消赞
// commentType - 操作类型，当删除评论时可用，需与评论type字段一致
func (bot *Bot) SNSObjectOperation(momentID, commentID string,
	Type, commentType int) (*MsgAndStatus, error) {
	resp := bot.sendCommand("snsobjectOp", struct {
		MomentID    string `json:"momentId"`
		Type        int    `json:"type"`
		CommentID   string
		CommentType int `json:"commentType"`
	}{
		MomentID:    momentID,
		Type:        Type,
		CommentID:   commentID,
		CommentType: commentType,
	})
	if !resp.Success {
		return nil, errors.New(resp.Msg)
	}
	data := &MsgAndStatus{}
	err := jsoniter.Unmarshal(resp.Data, data)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// SNSSendMoment 发朋友圈
// content - 文本内容或 TimeLineObject 结构体文本
func (bot *Bot) SNSSendMoment(content string) (*MomentResp, error) {
	resp := bot.sendCommand("snsSendMoment", struct {
		Content string `json:"content"`
	}{
		Content: content,
	})
	if !resp.Success {
		return nil, errors.New(resp.Msg)
	}
	data := &MomentResp{}
	err := jsoniter.Unmarshal(resp.Data, data)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// SNSUserPage 查看用户朋友圈
// momentID 首次传入空即获取第一页, 以后传入上次拉取的最后一条信息ID
func (bot *Bot) SNSUserPage(userID, momentID string) (*MomentListResp, error) {
	resp := bot.sendCommand("snsUserPage", struct {
		UserID   string `json:"userId"`
		MomentID string `json:"momentId"`
	}{
		UserID:   userID,
		MomentID: momentID,
	})
	if !resp.Success {
		return nil, errors.New(resp.Msg)
	}
	data := &MomentListResp{}
	err := jsoniter.Unmarshal(resp.Data, data)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// SNSTimeLine 查看朋友圈动态
// momentID 首次传入空即获取第一页, 以后传入上次拉取的最后一条信息ID
func (bot *Bot) SNSTimeLine(momentID string) (*MomentListResp, error) {
	resp := bot.sendCommand("snsTimeline", struct {
		MomentID string `json:"momentId"`
	}{
		MomentID: momentID,
	})
	if !resp.Success {
		return nil, errors.New(resp.Msg)
	}
	data := &MomentListResp{}
	err := jsoniter.Unmarshal(resp.Data, data)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// SNSGetObject 获取朋友圈信息详情
func (bot *Bot) SNSGetObject(momentID string) (*MomentDetailResp, error) {
	resp := bot.sendCommand("snsGetObject", struct {
		MomentID string `json:"momentId"`
	}{
		MomentID: momentID,
	})
	if !resp.Success {
		return nil, errors.New(resp.Msg)
	}
	data := &MomentDetailResp{}
	err := jsoniter.Unmarshal(resp.Data, data)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// SNSComment 评论朋友圈
func (bot *Bot) SNSComment(userID, momentID, content string) (*MomentDetailResp, error) {
	resp := bot.sendCommand("snsComment", struct {
		UserID   string `json:"userId"`
		MomentID string `json:"momentId"`
		Content  string `json:"content"`
	}{
		UserID:   userID,
		MomentID: momentID,
		Content:  content,
	})
	if !resp.Success {
		return nil, errors.New(resp.Msg)
	}
	data := &MomentDetailResp{}
	err := jsoniter.Unmarshal(resp.Data, data)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// SNSLike 朋友圈点赞
func (bot *Bot) SNSLike(userID, momentID string) (*MomentDetailResp, error) {
	resp := bot.sendCommand("snsLike", struct {
		UserID   string `json:"userId"`
		MomentID string `json:"momentId"`
	}{
		UserID:   userID,
		MomentID: momentID,
	})
	if !resp.Success {
		return nil, errors.New(resp.Msg)
	}
	data := &MomentDetailResp{}
	err := jsoniter.Unmarshal(resp.Data, data)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// SyncFav 同步收藏消息
func (bot *Bot) SyncFav(favKey string) (*FavListResp, error) {
	resp := bot.sendCommand("syncFav", struct {
		FavKey string `json:"favKey"`
	}{
		FavKey: favKey,
	})
	if !resp.Success {
		return nil, errors.New(resp.Msg)
	}
	data := &FavListResp{}
	err := jsoniter.Unmarshal(resp.Data, data)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// AddFav 添加收藏
func (bot *Bot) AddFav(content string) (*AddFavResp, error) {
	resp := bot.sendCommand("addFav", struct {
		Content string `json:"content"`
	}{
		Content: content,
	})
	if !resp.Success {
		return nil, errors.New(resp.Msg)
	}
	data := &AddFavResp{}
	err := jsoniter.Unmarshal(resp.Data, data)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// GetFav 获取收藏消息详情
func (bot *Bot) GetFav(favID int) (*FavResp, error) {
	resp := bot.sendCommand("getFav", struct {
		FavID int `json:"favId"`
	}{
		FavID: favID,
	})
	if !resp.Success {
		return nil, errors.New(resp.Msg)
	}
	data := &FavResp{}
	err := jsoniter.Unmarshal(resp.Data, data)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// DeleteFav 删除收藏
func (bot *Bot) DeleteFav(favID int) (*FavResp, error) {
	resp := bot.sendCommand("deleteFav", struct {
		FavID int `json:"favId"`
	}{
		FavID: favID,
	})
	if !resp.Success {
		return nil, errors.New(resp.Msg)
	}
	data := &FavResp{}
	err := jsoniter.Unmarshal(resp.Data, data)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// GetLabelList 获取所有标签
func (bot *Bot) GetLabelList() (*LabelListResp, error) {
	resp := bot.sendCommand("getLabelList", nil)
	if !resp.Success {
		return nil, errors.New(resp.Msg)
	}
	data := &LabelListResp{}
	err := jsoniter.Unmarshal(resp.Data, data)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// AddLabel 添加标签
func (bot *Bot) AddLabel(label string) (*MsgAndStatus, error) {
	resp := bot.sendCommand("addLabel", struct {
		Label string `json:"label"`
	}{Label: label})
	if !resp.Success {
		return nil, errors.New(resp.Msg)
	}
	data := &MsgAndStatus{}
	err := jsoniter.Unmarshal(resp.Data, data)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// DeleteLabel 删除标签
func (bot *Bot) DeleteLabel(labelID int) (*MsgAndStatus, error) {
	resp := bot.sendCommand("deleteLabel", struct {
		LabelID int `json:"labelId"`
	}{LabelID: labelID})
	if !resp.Success {
		return nil, errors.New(resp.Msg)
	}
	data := &MsgAndStatus{}
	err := jsoniter.Unmarshal(resp.Data, data)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// SetLabel 设置用户标签
func (bot *Bot) SetLabel(userID string, labelID int) (*MsgAndStatus, error) {
	resp := bot.sendCommand("setLabel", struct {
		LabelID int    `json:"labelId"`
		UserID  string `json:"userId"`
	}{LabelID: labelID})
	if !resp.Success {
		return nil, errors.New(resp.Msg)
	}
	data := &MsgAndStatus{}
	err := jsoniter.Unmarshal(resp.Data, data)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// ReceiveRedPacket 接收红包
// mType = 49
// 使用 bytes.Contains(msg.Content, []byte("<![CDATA[微信红包]]>")) 与转账区分
func (bot *Bot) ReceiveRedPacket(rawMsgData Msg) (*ExternalMsgResp, error) {
	rawMsgData.Data = ""
	resp := bot.sendCommand("receiveRedPacket", struct {
		RawMsgData Msg `json:"rawMsgData"`
	}{RawMsgData: rawMsgData})
	if !resp.Success {
		return nil, errors.New(resp.Msg)
	}
	data := &ExternalMsgResp{}
	err := jsoniter.Unmarshal(resp.Data, data)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// QueryRedPacket 查看红包信息, 如果是别人发的红包, 未领取且未领取完毕时, 无法取到红包信息
func (bot *Bot) QueryRedPacket(rawMsgData Msg, index int) (*ExternalMsgResp, error) {
	rawMsgData.Data = ""
	resp := bot.sendCommand("queryRedPacket", struct {
		RawMsgData Msg `json:"rawMsgData"`
		Index      int `json:"index"`
	}{
		Index:      index,
		RawMsgData: rawMsgData,
	})
	if !resp.Success {
		return nil, errors.New(resp.Msg)
	}
	data := &ExternalMsgResp{}
	err := jsoniter.Unmarshal(resp.Data, data)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// OpenRedPacket 领取红包
func (bot *Bot) OpenRedPacket(rawMsgData Msg, key string) (*ExternalMsgResp, error) {
	rawMsgData.Data = ""
	resp := bot.sendCommand("openRedPacket", struct {
		RawMsgData Msg    `json:"rawMsgData"`
		Key        string `json:"key"`
	}{
		Key:        key,
		RawMsgData: rawMsgData,
	})
	if !resp.Success {
		return nil, errors.New(resp.Msg)
	}
	data := &ExternalMsgResp{}
	err := jsoniter.Unmarshal(resp.Data, data)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// QueryTransfer 查看转账消息
func (bot *Bot) QueryTransfer(rawMsgData Msg) (*ExternalMsgResp, error) {
	rawMsgData.Data = ""
	resp := bot.sendCommand("queryTransfer", struct {
		RawMsgData Msg `json:"rawMsgData"`
	}{
		RawMsgData: rawMsgData,
	})
	if !resp.Success {
		return nil, errors.New(resp.Msg)
	}
	data := &ExternalMsgResp{}
	err := jsoniter.Unmarshal(resp.Data, data)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// AcceptTransfer 接受转账
// mType = 49
// 使用 bytes.Contains(msg.Content, []byte("<![CDATA[微信转账]]>")) 与红包区分
func (bot *Bot) AcceptTransfer(rawMsgData Msg) (*ExternalMsgResp, error) {
	rawMsgData.Data = ""
	resp := bot.sendCommand("acceptTransfer", struct {
		RawMsgData Msg `json:"rawMsgData"`
	}{
		RawMsgData: rawMsgData,
	})
	if !resp.Success {
		return nil, errors.New(resp.Msg)
	}
	data := &ExternalMsgResp{}
	err := jsoniter.Unmarshal(resp.Data, data)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// SearchMp 搜索公众号
func (bot *Bot) SearchMp(content string) (*SearchMPResp, error) {
	resp := bot.sendCommand("searchMp", struct {
		Content string `json:"content"`
	}{
		Content: content,
	})
	if !resp.Success {
		return nil, errors.New(resp.Msg)
	}
	data := &SearchMPResp{}
	err := jsoniter.Unmarshal(resp.Data, data)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// GetSubscriptionInfo 获取公众号信息
func (bot *Bot) GetSubscriptionInfo(ghName string) (*SearchMPResp, error) {
	resp := bot.sendCommand("getSubscriptionInfo", struct {
		GhName string `json:"ghName"`
	}{
		GhName: ghName,
	})
	if !resp.Success {
		return nil, errors.New(resp.Msg)
	}
	data := &SearchMPResp{}
	err := jsoniter.Unmarshal(resp.Data, data)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// OperateSubscription 操作公众号菜单
func (bot *Bot) OperateSubscription(ghName string, menuId int, menuKey string) (*MsgAndStatus, error) {
	resp := bot.sendCommand("operateSubscription", struct {
		GhName  string `json:"ghName"`
		MenuID  int    `json:"menuId"`
		MenuKey string `json:"menuKey"`
	}{
		GhName:  ghName,
		MenuID:  menuId,
		MenuKey: menuKey,
	})
	if !resp.Success {
		return nil, errors.New(resp.Msg)
	}
	data := &MsgAndStatus{}
	err := jsoniter.Unmarshal(resp.Data, data)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// GetRequestToken 获取网页访问授权
func (bot *Bot) GetRequestToken(ghName, url string) (*RequestTokenResp, error) {
	resp := bot.sendCommand("getRequestToken", struct {
		GhName string `json:"ghName"`
		URL    string `json:"url"`
	}{
		GhName: ghName,
		URL:    url,
	})
	if !resp.Success {
		return nil, errors.New(resp.Msg)
	}
	data := &RequestTokenResp{}
	err := jsoniter.Unmarshal(resp.Data, data)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// RequestUrl 访问网页
func (bot *Bot) RequestUrl(url, xKey, xUin string) (*RequestUrlResp, error) {
	resp := bot.sendCommand("requestUrl", struct {
		URL  string `json:"url"`
		XKey string `json:"xKey"`
		XUin string `json:"xUin"`
	}{
		URL:  url,
		XKey: xKey,
		XUin: xUin,
	})
	if !resp.Success {
		return nil, errors.New(resp.Msg)
	}
	data := &RequestUrlResp{}
	err := jsoniter.Unmarshal(resp.Data, data)
	if err != nil {
		return nil, err
	}
	return data, nil
}
