import './globals.css';
import { Route, Routes } from 'react-router-dom';
import LoginForm from './_auth/forms/LoginForm';
import SignupForm from './_auth/forms/SignupForm';
import { Home } from './_root/pages';
import RootLayout from './_root/RootLayout';
import AuthLayout from './_auth/AuthLayout';

function App() {
  return (
    <main className="flex h-screen">
      <Routes>
        {/* public routes */}
        <Route element={<AuthLayout />}>
          <Route path="/login" element={<LoginForm />} />
          <Route path="/signup" element={<SignupForm />} />
        </Route>
        {/* private routes */}
        <Route element={<RootLayout />}>
          <Route path="/home" element={<Home />} />
        </Route>
      </Routes>
    </main>
  );
}

export default App;
