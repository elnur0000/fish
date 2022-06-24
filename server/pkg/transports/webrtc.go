package transports

import (
	"encoding/base64"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/pion/webrtc/v3"
)

type webRTCTransport struct {
	onConnHandler func(conn Connection)
}

var WebRTCTransport webRTCTransport = webRTCTransport{}

type webRTCConn struct {
	peerConnection *webrtc.PeerConnection
	dataChannel    *webrtc.DataChannel
	recvBuffer     chan []byte
}

func NewWebRTCConn(dataChannel *webrtc.DataChannel, peerConn *webrtc.PeerConnection) webRTCConn {
	recvBuffer := make(chan []byte, 100)

	go func() {
		dataChannel.OnMessage(func(msg webrtc.DataChannelMessage) {
			recvBuffer <- msg.Data
		})
	}()
	return webRTCConn{
		recvBuffer:     recvBuffer,
		dataChannel:    dataChannel,
		peerConnection: peerConn,
	}
}

func (w webRTCConn) Send(b []byte) {
	w.dataChannel.Send(b)
}

func (w webRTCConn) Recv() chan []byte {
	return w.recvBuffer
}

func (w webRTCConn) getConnId() *uint16 {
	return w.dataChannel.ID()
}

func (w webRTCConn) OnDisconnect(cb func()) {
	w.peerConnection.OnConnectionStateChange(func(pcs webrtc.PeerConnectionState) {
		if pcs == webrtc.PeerConnectionStateDisconnected {
			cb()
		}
	})
}

func (webRTCTransport) CreatePeerConn(offerStr string) (*webrtc.PeerConnection, error) {
	config := webrtc.Configuration{
		ICEServers: []webrtc.ICEServer{
			{
				URLs: []string{"stun:stun.l.google.com:19302"},
			},
		},
	}

	peerConnection, err := webrtc.NewPeerConnection(config)

	if err != nil {
		return nil, err
	}

	// Register data channel creation handling
	peerConnection.OnDataChannel(func(d *webrtc.DataChannel) {
		d.OnOpen(func() {
			log.Printf("Data channel '%s'-'%d' is open.", d.Label(), d.ID())
			WebRTCTransport.onConnHandler(NewWebRTCConn(d, peerConnection))
		})
	})

	// Set the remote SessionDescription
	offer := webrtc.SessionDescription{}
	Decode(offerStr, &offer)
	err = peerConnection.SetRemoteDescription(offer)
	if err != nil {
		return nil, err
	}

	// Create an answer
	answer, err := peerConnection.CreateAnswer(nil)
	if err != nil {
		return nil, err
	}

	// Create channel that is blocked until ICE Gathering is complete
	gatherComplete := webrtc.GatheringCompletePromise(peerConnection)

	// Sets the LocalDescription, and starts our UDP listeners
	err = peerConnection.SetLocalDescription(answer)
	if err != nil {
		return nil, err
	}

	// Block until ICE Gathering is complete, disabling trickle ICE
	// we do this because we only can exchange one signaling message
	// in a production application you should exchange ICE Candidates via OnICECandidate
	<-gatherComplete

	return peerConnection, nil
}

func (wt *webRTCTransport) RegisterConnHandler(handler func(conn Connection)) {
	wt.onConnHandler = handler
}

func (wt webRTCTransport) SignalHandler(w http.ResponseWriter, r *http.Request) {

	res, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(nil)
	}

	offer := string(res)

	peerConnection, err := wt.CreatePeerConn(offer)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(nil)
	}

	answerBase64 := Encode(*peerConnection.LocalDescription())
	w.Write([]byte(answerBase64))
}

func CORS(h func(http.ResponseWriter, *http.Request)) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
		h(w, r)
	}
}

func (wt webRTCTransport) WaitForConnection() {
	http.HandleFunc("/signal", CORS(WebRTCTransport.SignalHandler))
}

// Encode encodes the input in base64
func Encode(obj interface{}) string {
	b, err := json.Marshal(obj)
	if err != nil {
		panic(err)
	}

	return base64.StdEncoding.EncodeToString(b)
}

func Decode(in string, obj interface{}) {
	b, err := base64.StdEncoding.DecodeString(in)
	if err != nil {
		panic(err)
	}

	err = json.Unmarshal(b, obj)
	if err != nil {
		panic(err)
	}
}
