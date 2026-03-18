'use client';

import { useParams } from 'next/navigation';
import { useVoteRoom, RoomLobby, VoteStrategy, ResultsScreen } from '@/features/vote';
import { Skeleton } from '@/shared/ui';

export default function VoteRoomPage() {
  const params = useParams<{ roomId: string }>();
  const {
    room,
    status,
    participants,
    results,
    startVoting,
    submitVote,
    countdown,
    isLoading,
    error,
    playAgain,
    changeVoteType,
  } = useVoteRoom(params.roomId);

  if (isLoading) {
    return (
      <main className="flex min-h-screen items-center justify-center">
        <div className="space-y-4 text-center">
          <Skeleton className="mx-auto h-12 w-48" />
          <Skeleton className="mx-auto h-6 w-32" />
        </div>
      </main>
    );
  }

  if (error || !room) {
    return (
      <main className="flex min-h-screen items-center justify-center px-4 text-center">
        <div>
          <div className="mb-4 text-4xl">😵</div>
          <h1 className="mb-2 font-heading text-xl font-bold">Không tìm thấy phòng</h1>
          <p className="text-sm text-gray-500">{error || 'Phòng không tồn tại hoặc đã hết hạn.'}</p>
        </div>
      </main>
    );
  }

  return (
    <main className="min-h-screen">
      {status === 'waiting' && (
        <RoomLobby
          roomCode={room.code}
          hostName={room.hostName}
          voteType={room.voteType}
          participants={participants}
          isHost={true}
          onStart={startVoting}
          onChangeVoteType={changeVoteType}
        />
      )}

      {status === 'voting' && room.dishes && (
        <VoteStrategy
          voteType={room.voteType}
          dishes={room.dishes}
          onVoteComplete={submitVote}
          timeRemaining={countdown}
        />
      )}

      {status === 'finished' && <ResultsScreen results={results} onPlayAgain={playAgain} />}
    </main>
  );
}
