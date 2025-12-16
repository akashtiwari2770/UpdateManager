import { AppRouter } from './router';
import { ToastContainer } from '@/components/notifications/ToastContainer';

function App() {
  return (
    <ToastContainer>
      <AppRouter />
    </ToastContainer>
  );
}

export default App;

