import Phaser from 'phaser';
import Piranha from '../entities/piranha';
import WebRTCConnection from '../network/webrtc';
import config from '../config'

type PlayersMap = {
  [id: string]: Phaser.GameObjects.Sprite
}

const generateBackground = (config, scene: Phaser.Scene) => {
  let graphics = new Phaser.GameObjects.Graphics(scene);
  scene.add.existing(graphics);

  let xMax = config.columns * config.lineSpacing;
  let yMax = config.rows * config.lineSpacing;

  graphics.lineStyle(config.lineThickness, config.lineColor, 1);
  graphics.beginPath();

  for(let x = 0; x <= xMax; x += config.lineSpacing) {
      graphics.moveTo(x, 0);
      graphics.lineTo(x, yMax);
  }

  for(let y = 0; y <= yMax; y += config.lineSpacing) {
      graphics.moveTo(0, y);
      graphics.lineTo(xMax, y);
  }

  graphics.closePath();
  graphics.strokePath();
}
export default class Demo extends Phaser.Scene {

  private players!: PlayersMap
  private webRTCConnection!: WebRTCConnection

  constructor() {
    super('GameScene');
  }

  preload() {
    this.players = {}
    this.load.image('background', 'assets/background.png');
    this.load.atlas("piranha", ['assets/pirana-sprite-0.png', 'assets/pirana-sprite-1.png'], 'assets/pirana-sprite.json')
  }
  
  create() {
    this.webRTCConnection = new WebRTCConnection()
    this.webRTCConnection.init()

    this.webRTCConnection.onMessage((data: ArrayBuffer) => {
      window.processServerMessage(new Uint8Array(data))
    })
    
    // generateBackground({
    //   columns: 50,
    //   rows: 50,
    //   lineSpacing: 100,
    //   lineThickness: 5,
    //   lineColor: 0xff00ff,
    // }, this)
    this.add.image(5000 / 2, 5000 / 2, 'background').setDisplaySize(5000, 5000)

  }

  update(time: number, delta: number): void {
    const localPlayerId = window.getLocalPlayerId()
    if(!localPlayerId) return
  
    const playerStates = window.getPlayerStates()

    const playerChangeTrack: {
      [id: string]: boolean
    } = {}

    for(const state of playerStates) {
      let player = this.players[state.ID]

      const isLocalPlayer = state.ID === localPlayerId

      playerChangeTrack[String(state.ID)] = true

      if (!player) {
        player = new Piranha(this, state.x, state.y, state.width, state.height)
        isLocalPlayer && this.cameras.main.startFollow(player, true)
        this.players[state.ID] = player
      }
      player.setPosition(state.x, state.y)
      player.setRotation(state.rotation)
      player.setSize(state.width, state.height)
      
      if(isLocalPlayer) {
        const mouseX = this.game.input.activePointer.x
        const mouseY = this.game.input.activePointer.y
        const newAngle = Phaser.Math.Angle.Between(mouseX, mouseY, this.cameras.main.centerX, this.cameras.main.centerY)
        window.processLocalPlayerInput(newAngle, delta, (msg) => this.webRTCConnection.sendMessage(msg))
      }
    }

    for(let key in this.players) {
      if(!(key in playerChangeTrack)){
        this.players[key].destroy()
        delete(this.players[key])
      }
    }

  }
}
