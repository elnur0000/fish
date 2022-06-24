import Phaser from 'phaser';

export default {
  type: Phaser.AUTO,
  parent: 'game',
  scale: {
      // Fit to window
      width: '100%',
      height:'100%',
      mode: Phaser.Scale.FIT
  },
  // physics: {
  //   default: 'arcade',
  //   arcade: {
  //       debug: true,
  //       fps: 60,
  //   }
  // },
  signalUrl: 'http://localhost:5000/signal'
};
