import 'server-only'
import { getCurrentGame } from '@/lib/matchmaker'

const Home = async () => {
  const game = await getCurrentGame()
  const gameObj = game.toObject()

  return (
    <div className='absolute flex w-full h-full'>
      <div className="flex w-full h-full justify-center items-center">
        <button className='flex p-2 rounded-xl border-white cursor-pointer border-2'>
          <div className='flex flex-col'>
            <p>game: {gameObj.id}</p>
            <p>{gameObj.player_count} / {gameObj.max_players}</p>
          </div>
        </button>
      </div>
    </div>
  );
}

export default Home
