import 'server-only'
import { getCurrentGame } from '@/lib/matchmaker'

const Home = async () => {
  getCurrentGame()

  return (
    <div>

    </div>
  );
}

export default Home
