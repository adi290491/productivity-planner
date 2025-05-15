import './App.css'
import { BrowserRouter as Router, Navigate, Route, Routes } from 'react-router-dom'
import Signup from './pages/Signup';
import Login from './pages/Login';
import Dashboard from './pages/Dashboard';
import PrivateRoute from './components/PrivateRoute';

function App() {
  return (
    <Router>
      <Routes>
        <Route path="/signup" element={<Signup />} />
        <Route path="/login" element={<Login />} />
        <Route path="/"
        element= {
          <PrivateRoute>
            <Dashboard />
          </PrivateRoute>
        }
        />
        <Route path="*" element={<Navigate to="/" />} />
      </Routes>
    </Router>
  )
}

export default App;
