import Phaser from 'phaser';
import config from './config';
import GameScene from './scenes/Game';
import * as WASM from './wasm'

WASM.init()

new Phaser.Game(
  Object.assign(config, {
    scene: [GameScene]
  })
);
