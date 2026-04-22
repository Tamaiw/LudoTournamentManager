import { GamePairing } from '../../types';

interface Props {
  round: number;
  pairings: GamePairing[];
}

export function TableAssignment({ round, pairings }: Props) {
  return (
    <div className="mb-8">
      <h3 className="text-lg font-bold mb-4">Round {round}</h3>
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4">
        {pairings.map((pairing) => (
          <div key={pairing.gameId} className="bg-gray-50 rounded-lg p-3">
            <div className="font-semibold mb-2">Table {pairing.tableNumber}</div>
            <div className="text-sm space-y-1">
              {pairing.playerIds.map((playerId, idx) => (
                <div key={playerId} className="flex items-center gap-2">
                  <span className={`w-3 h-3 rounded-full ${
                    pairing.seatColors[idx] === 'yellow' ? 'bg-yellow-400' :
                    pairing.seatColors[idx] === 'green' ? 'bg-green-400' :
                    pairing.seatColors[idx] === 'blue' ? 'bg-blue-400' :
                    'bg-red-400'
                  }`} />
                  <span>{playerId.slice(0, 8)}...</span>
                </div>
              ))}
            </div>
          </div>
        ))}
      </div>
    </div>
  );
}