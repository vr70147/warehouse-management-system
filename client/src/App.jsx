import './index.css';
import TopBar from './components/top-bar/TopBar';
import { BrowserRouter, Routes, Route } from 'react-router-dom';
import SideBar from './components/sidebar/SideBar';
import DashboardPage from './pages/DashboardPage';
import { useState } from 'react';
import { Provider } from 'react-redux';
import { store } from './redux/store';
import InventoryPage from './pages/InventoryPage';
import { Toaster } from 'sonner';
import OrdersPage from './pages/OrdersPage';

function App() {
  const [isCollapsed, setIsCollapsed] = useState(false);

  return (
    <Provider store={store}>
      <BrowserRouter>
        <div className="h-screen bg-gray-100 dark:bg-gray-900 flex flex-col">
          <header className="fixed top-0 left-0 right-0 z-50">
            <TopBar />
          </header>

          <div className="flex flex-grow pt-16">
            <SideBar
              isCollapsed={isCollapsed}
              setIsCollapsed={setIsCollapsed}
            />
            <Toaster position="bottom-right" expand />
            <main
              className={`flex-grow overflow-y-auto transition-all duration-300 `}
            >
              <Routes>
                <Route path="/dashboard" element={<DashboardPage />} />
                <Route path="/inventory" element={<InventoryPage />} />
                <Route path="/order" element={<OrdersPage />} />
              </Routes>
            </main>
          </div>
        </div>
      </BrowserRouter>
    </Provider>
  );
}

export default App;
