import 'server-only'
import { getCurrentGame } from '@/lib/matchmaker'
import { JoinButton } from '@/components/join-button'
import { Reloader } from '@/components/reloader'

const Home = async () => {
  const game = await getCurrentGame()

  return (
    <>
      <Reloader />
      <div className='absolute flex w-full h-full'>
        <div className="flex w-full h-full justify-center items-center">
          <JoinButton game={game} />
        </div>
      </div>
    </>
  )
}

export default Home
