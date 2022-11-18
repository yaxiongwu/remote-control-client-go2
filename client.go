package engine

import (
	"github.com/pion/ice/v2"
	log "github.com/pion/ion-log"
	"github.com/pion/webrtc/v3"
	"github.com/yaxiongwu/remote-control-client-go2/pkg/proto/rtc"
)

// Client is pub/sub transport
type Client struct {
	Id  string
	api *webrtc.DataChannel
	rtc *RTC
	pc  *webrtc.PeerConnection
	//role           Target
	role                       rtc.Role
	SendCandidates             []*webrtc.ICECandidate
	RecvCandidates             []webrtc.ICECandidateInit
	config                     *RTCConfig
	OnIceConnectionStateChange func(webrtc.ICEConnectionState, *webrtc.PeerConnection)
}

// NewTransport create a transport
func NewClient(uid string, rtc *RTC, role rtc.Role) *Client {
	c := &Client{
		Id:     uid,
		config: &DefaultConfig,
		rtc:    rtc,
		role:   role,
	}

	c.SendCandidates = []*webrtc.ICECandidate{}

	var err error
	var api *webrtc.API
	var me *webrtc.MediaEngine
	c.config.WebRTC.Setting.SetICEMulticastDNSMode(ice.MulticastDNSModeDisabled)
	// if role == Target_PUBLISHER {
	// 	me, err = getPublisherMediaEngine(rtc.config.WebRTC.VideoMime)
	// } else {
	me, err = getSubscriberMediaEngine()
	//}

	if err != nil {
		log.Errorf("getPublisherMediaEngine error: %v", err)
		return nil
	}

	api = webrtc.NewAPI(webrtc.WithMediaEngine(me), webrtc.WithSettingEngine(c.config.WebRTC.Setting))
	c.pc, err = api.NewPeerConnection(c.config.WebRTC.Configuration)

	if err != nil {
		log.Errorf("NewPeerConnection error: %v", err)
		return nil
	}

	// if role == Target_PUBLISHER {
	// 	log.Debugf("t.pc.CreateDataChannel(API_CHANNEL)")
	// 	_, err = t.pc.CreateDataChannel(API_CHANNEL, &webrtc.DataChannelInit{})

	// 	if err != nil {
	// 		log.Errorf("error creating data channel: %v", err)
	// 		return nil
	// 	}
	// }

	c.pc.OnICECandidate(func(i *webrtc.ICECandidate) {
		log.Debugf("t.pc.OnICECandidate,myid:%v,%v", uid, i)
		if i == nil {
			// Gathering done
			log.Infof("gather candidate done")
			return
		}
		//append before join session success
		if c.pc.CurrentRemoteDescription() == nil {
			c.SendCandidates = append(c.SendCandidates, i)
		} else {
			for _, cand := range c.SendCandidates {
				c.rtc.SendTrickle2(cand, uid)
			}
			c.SendCandidates = []*webrtc.ICECandidate{}
			c.rtc.SendTrickle2(i, uid)
		}
	})
	c.pc.OnICEConnectionStateChange(func(state webrtc.ICEConnectionState) {
		if rtc.OnIceConnectionStateChange != nil {
			rtc.OnIceConnectionStateChange(state, c.pc)
		}
		log.Debugf("IECConnectionStateChange to %v", state)
	})

	return c
}

func (t *Client) GetPeerConnection() *webrtc.PeerConnection {
	return t.pc
}
