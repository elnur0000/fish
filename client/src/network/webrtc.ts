import config from '../config'

export default class WebRTCConnection {
    private _pc: RTCPeerConnection
    private _channel: Phaser.Events.EventEmitter
    private _dataChannel!: RTCDataChannel

    constructor() {
        this._pc = new RTCPeerConnection({
            iceServers: [
              {
                urls: 'stun:stun.l.google.com:19302'
              }
            ]
        })

        this._channel = new Phaser.Events.EventEmitter()

    }

    init () {
        return new Promise(async (resolve, reject) => {
            this._dataChannel = this._pc.createDataChannel('chan')
            this._dataChannel.onerror = () => reject(new Error('failed to create a Webrtc datachannel'))
            this._dataChannel.onopen = () => resolve(null)
            this._dataChannel.onmessage = (ev) => this._channel.emit('data', ev.data)
            const sdp = await this.getLocalSdp()
            const remoteSdp = await this.signal(sdp)
            
            this._pc.setRemoteDescription(new RTCSessionDescription(JSON.parse(atob(remoteSdp))))
        })
    }

    onMessage (msgHandler: (data: ArrayBuffer) => void) {
        return this._channel.on('data', msgHandler)
    }

    sendMessage (msg: ArrayBuffer) {
        this._dataChannel.send(msg)
    }

    private getLocalSdp (): Promise<string> {
        return new Promise((resolve, reject) => {
            this._pc.oniceconnectionstatechange = e => console.log(this._pc.iceConnectionState)
            this._pc.onnegotiationneeded = e => this._pc.createOffer().then(d => this._pc.setLocalDescription(d)).catch(console.log)
            this._pc.onicecandidate = event => {
                console.log(event)
                if (event.candidate === null) {
                    resolve(btoa(JSON.stringify(this._pc.localDescription)))
                }
            }
        })
    }

    private signal (sdp: string): Promise<string> {
        return new Promise((resolve, reject) => {
            const http = new XMLHttpRequest()
            http.open('POST', config.signalUrl, true)
        
            http.setRequestHeader('Content-type', 'text/plain')
        
            http.onreadystatechange = function() {
                if(http.readyState == 4 && http.status == 200) {
                    resolve(http.responseText)
                }
            }

            http.send(sdp)
        }) 
    }

}