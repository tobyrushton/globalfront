import 'server-only'
import { getCurrentGame } from '@/lib/matchmaker'

const Home = async () => {
  const game = await getCurrentGame()

  return (
    <div>

    </div>
  );
}

export default Home
