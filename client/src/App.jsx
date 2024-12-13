import './index.css';
import TopBar from './components/top-bar/TopBar';
import { BrowserRouter, Route, Routes } from 'react-router-dom';
import SideBar from './components/sidebar/SideBar';
import Dashboard from './components/features/dashboard/Dashboard';
import Inventory from './components/features/inventory/Inventory';

function App() {
  return (
    <div className="h-screen overflow-hidden">
      <header className="fixed top-0 left-0 right-0 z-50 opacity-90">
        <TopBar />
      </header>
      <BrowserRouter>
        <aside className="fixed left-0 top-16 bottom-0 w-64">
          <SideBar />
        </aside>
        <main className="ml-64 flex-1 overflow-auto p-4 bg-gray-100 dark:bg-gray-900">
          <Routes>
            <Route path="/dashboard" element={<Dashboard />} />
            <Route path="/inventory" element={<Inventory />} />
          </Routes>
        </main>
      </BrowserRouter>
    </div>
  );
}

export default App;
