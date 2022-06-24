import Phaser from 'phaser'

export default class Piranha extends Phaser.GameObjects.Sprite {
    constructor(scene: Phaser.Scene, x: number, y: number, width: number = 80, height: number = 80) {
        const textureName = 'piranha'
        super(scene, x, y, textureName)
        
        this.setDisplaySize(width, height)
        
        // const frameNames = this.textures.get('piranha').getFrameNames()
        this.anims.create({
            key: 'swim',
            frames: [
                { 
                    "key": textureName,
                    "frame": "__piranha_swim_001.png"
                },
                {
                    "key": textureName,
                    "frame": "__piranha_swim_002.png"
                },
                {
                    "key": textureName,
                    "frame": "__piranha_swim_003.png"
                },
                {
                    "key": textureName,
                    "frame": "__piranha_swim_004.png"
                },
                {
                    "key": textureName,
                    "frame": "__piranha_swim_005.png"
                },
                {
                    "key": textureName,
                    "frame": "__piranha_swim_007.png"
                }
            ],
            frameRate: 20,
            repeat: -1
        })

        this.play('swim')

        this.scene.add.existing(this)
    }

    setRotation(radians?: number | undefined): this {
        if (radians) {
            if (radians > Math.PI / 2 || radians < -Math.PI / 2) {
                this.setFlipY(true)
            } else {
                this.setFlipY(false)
            }
        }

        return super.setRotation.bind(this)(radians)
    }
}