package padchat

import (
	"errors"

	"github.com/google/uuid"
	"github.com/json-iterator/go"
)

func (bot *Bot) sendCommand(cmd string, data interface{}) CommandResp {
	c := make(chan CommandResp)
	id := uuid.New().String()
	bot.retProcMap.Store(id, func(resp CommandResp) {
		c <- resp
	})
	bot.ws.WriteJSON(WSReq{Type: "user", CMD: cmd, CMDID: id, Data: data})
	return <-c
}

func (bot *Bot) Init() CommandResp {
	return bot.sendCommand("init", nil)
}

func (bot *Bot) QRLogin() CommandResp {
	return bot.sendCommand("login", LoginReq{LoginType: "qrcode"})
}

func (bot *Bot) RequestLogin(wxData, token string) CommandResp {
	return bot.sendCommand("login", LoginReq{
		LoginType: "request",
		WXData:    wxData,
		Token:     token,
	})
}

func (bot *Bot) TokenLogin(wxData, token string) CommandResp {
	return bot.sendCommand("login", LoginReq{
		LoginType: "token",
		WXData:    wxData,
		Token:     token,
	})
}

func (bot *Bot) UserLogin(wxData, username, password string) CommandResp {
	return bot.sendCommand("login", LoginReq{
		LoginType: "user",
		WXData:    wxData,
		UserName:  username,
		Password:  password,
	})
}

func (bot *Bot) PhoneLogin(wxData, phone, code string) CommandResp {
	return bot.sendCommand("login", LoginReq{
		LoginType: "phone",
		WXData:    wxData,
		Phone:     phone,
		Code:      code,
	})
}

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

func (bot *Bot) SyncContact() CommandResp {
	return bot.sendCommand("syncContact", nil)
}

func (bot *Bot) SendMsg(req SendMsgReq) (*SendMsgResp, error) {
	data := &SendMsgResp{}
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

func (bot *Bot) SNSUpload(file string) (*ImgResp, error) {
	resp := bot.sendCommand("snsUpload", struct {
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
