'use client';

import { useCallback, useEffect, useRef, useState } from 'react';
import {
  getSocketClient,
  connectSocket,
  disconnectSocket,
  joinRoom,
  leaveRoom,
} from '@shared/lib/socket-client';
import { useVoteStore } from '../stores/vote-store';
import type { Participant, VoteResult, TournamentRound } from '../types';

interface UseSocketReturn {
  isConnected: boolean;
  emit: (event: string, data: unknown) => void;
}

export function useSocket(roomId: string): UseSocketReturn {
  const [isConnected, setIsConnected] = useState(false);
  const roomIdRef = useRef(roomId);
  roomIdRef.current = roomId;

  const {
    addParticipant,
    receiveVoteUpdate,
    startVoting,
    setResults,
    setCountdown,
    setCurrentRound,
  } = useVoteStore();

  useEffect(() => {
    const socket = getSocketClient();

    function onConnect() {
      setIsConnected(true);
      joinRoom(roomIdRef.current);
    }

    function onDisconnect() {
      setIsConnected(false);
    }

    function onParticipantJoined(data: { participant: Participant }) {
      addParticipant(data.participant);
    }

    function onVoteSubmitted(data: { participantName: string }) {
      receiveVoteUpdate(data.participantName);
    }

    function onRoundComplete(data: { round: TournamentRound }) {
      setCurrentRound(data.round);
    }

    function onVotingStarted() {
      startVoting();
    }

    function onVotingFinished(data: { results: VoteResult[] }) {
      setResults(data.results);
    }

    function onTimerTick(data: { remaining: number }) {
      setCountdown(data.remaining);
    }

    socket.on('connect', onConnect);
    socket.on('disconnect', onDisconnect);
    socket.on('participant-joined', onParticipantJoined);
    socket.on('vote-submitted', onVoteSubmitted);
    socket.on('round-complete', onRoundComplete);
    socket.on('voting-started', onVotingStarted);
    socket.on('voting-finished', onVotingFinished);
    socket.on('timer-tick', onTimerTick);

    connectSocket();

    return () => {
      socket.off('connect', onConnect);
      socket.off('disconnect', onDisconnect);
      socket.off('participant-joined', onParticipantJoined);
      socket.off('vote-submitted', onVoteSubmitted);
      socket.off('round-complete', onRoundComplete);
      socket.off('voting-started', onVotingStarted);
      socket.off('voting-finished', onVotingFinished);
      socket.off('timer-tick', onTimerTick);
      leaveRoom(roomIdRef.current);
      disconnectSocket();
    };
  }, [addParticipant, receiveVoteUpdate, startVoting, setResults, setCountdown, setCurrentRound]);

  const emit = useCallback((event: string, data: unknown) => {
    const socket = getSocketClient();
    socket.emit(event, data);
  }, []);

  return { isConnected, emit };
}
