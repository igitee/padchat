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

