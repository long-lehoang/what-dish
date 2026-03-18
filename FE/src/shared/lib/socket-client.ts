import { io, type Socket } from 'socket.io-client';
import { WS_BASE_URL } from '@shared/constants';

let socketInstance: Socket | null = null;

export function getSocketClient(): Socket {
  if (socketInstance) {
    return socketInstance;
  }

  socketInstance = io(WS_BASE_URL, {
    autoConnect: false,
    reconnection: true,
    reconnectionAttempts: 5,
    reconnectionDelay: 1000,
    reconnectionDelayMax: 5000,
    transports: ['websocket', 'polling'],
  });

  return socketInstance;
}

export function getConnectionState(): boolean {
  return socketInstance?.connected ?? false;
}

export function connectSocket(): void {
  const socket = getSocketClient();
  if (!socket.connected) {
    socket.connect();
  }
}

export function disconnectSocket(): void {
  if (socketInstance) {
    socketInstance.disconnect();
    socketInstance = null;
  }
}

export function joinRoom(roomId: string): void {
  const socket = getSocketClient();
  if (!socket.connected) {
    socket.connect();
  }
  socket.emit('room:join', { roomId });
}

export function leaveRoom(roomId: string): void {
  const socket = getSocketClient();
  socket.emit('room:leave', { roomId });
}

export function onEvent<T>(event: string, callback: (data: T) => void): () => void {
  const socket = getSocketClient();
  socket.on(event, callback as (...args: unknown[]) => void);

  return () => {
    socket.off(event, callback as (...args: unknown[]) => void);
  };
}

export function emitEvent<T>(event: string, data: T): void {
  const socket = getSocketClient();
  socket.emit(event, data);
}
