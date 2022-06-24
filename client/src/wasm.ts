type PlayerState = {
    ID: number
    rotation: number
    x: number
    y: number
    velocity: number
    height: number
    width: number
}

declare global { interface Window { 
    processServerMessage: (msg: ArrayBuffer) => void
    getPlayerStates: () => PlayerState[] 
    processLocalPlayerInput: (rotation: number, duration: number, sendMessage: (msg: ArrayBuffer) => void) => void 
    getLocalPlayerId: () => number | undefined 
}}

export const init = () => {
    const go = new Go(); // Defined in wasm_exec.js
    const WASM_URL = 'http://localhost:5000/main.wasm';

    var wasm;

    if ('instantiateStreaming' in WebAssembly) {
        WebAssembly.instantiateStreaming(fetch(WASM_URL), go.importObject).then(function (obj) {
            wasm = obj.instance;
            go.run(wasm);
        })
    } else {
        fetch(WASM_URL).then(resp =>
            resp.arrayBuffer()
        ).then(bytes =>
            WebAssembly.instantiate(bytes, go.importObject).then(function (obj) {
                wasm = obj.instance;
                go.run(wasm);
            })
        )
    }
}