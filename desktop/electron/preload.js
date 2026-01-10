const { contextBridge, ipcRenderer } = require('electron');

// Minimal API — expose nothing sensitive, but allow future IPC if needed
contextBridge.exposeInMainWorld('electronAPI', {
    ping: () => 'pong'
});
