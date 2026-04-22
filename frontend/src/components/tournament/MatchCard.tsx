import { GamePairing } from '../../types';

interface Props {
  pairing: GamePairing;
  onReport?: (matchId: string) => void;
}

export function MatchCard({ pairing, onReport }: Props) {
  return (
    <div className="bg-white rounded-lg shadow p-4 border border-gray-200">
      <div className="flex justify-between items-center mb-2">
        <span className="font-bold">Table {pairing.tableNumber}</span>
        <span className={`px-2 py-1 rounded text-sm ${
          pairing.status === 'completed'
            ? 'bg-green-100 text-green-800'
            : 'bg-yellow-100 text-yellow-800'
        }`}>
          {pairing.status}
        </span>
      </div>

      <div className="space-y-1">
        {pairing.playerIds.map((playerId, idx) => (
          <div key={playerId} className="flex items-center gap-2">
            <span className={`w-4 h-4 rounded-full ${
              pairing.seatColors[idx] === 'yellow' ? 'bg-yellow-400' :
              pairing.seatColors[idx] === 'green' ? 'bg-green-400' :
              pairing.seatColors[idx] === 'blue' ? 'bg-blue-400' :
              'bg-red-400'
            }`} />
            <span className="text-sm">Player {playerId.slice(0, 8)}</span>
          </div>
        ))}
      </div>

      {pairing.status === 'pending' && onReport && (
        <button
          onClick={() => onReport(pairing.gameId)}
          className="mt-3 w-full bg-blue-600 text-white rounded px-3 py-1 text-sm hover:bg-blue-700"
        >
          Report Result
        </button>
      )}
    </div>
  );
}