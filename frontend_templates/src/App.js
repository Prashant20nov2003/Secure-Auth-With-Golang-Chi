import logo from './logo.svg';
import './App.css';
import {
  Link,
  Route,
  BrowserRouter as Router,
  Routes,
  useLocation
} from "react-router-dom"
import Home from './page/Home';
import Register from './page/Register';
import Login from './page/Login';
import EmailVerification from './page/EmailVerification';
import ChangePassword from './page/ChangePassword';
import PostProduct from './page/PostProduct';
import { useEffect, useState } from 'react';

function App() {
  return (
    <Router>
      <AppContent />
    </Router>
  );
}

function AppContent() {
  const [username, setUsername] = useState(null)
  const location = useLocation();

  useEffect(() => {
    fetchUsername();
  }, [location]); 

  const fetchUsername = async () => {
    try {
      const response = await fetch(`${process.env.REACT_APP_API_HOST}/api/getusername`,{
        method: "GET",
        headers: {
          'Content-Type': 'application/json',
        },
        credentials:'include'
      });
      const data = await response.json();

      if (data && !data.error) {
        setUsername(data.username);
      } else {
        throw new Error(data?.error || "Failed to fetch username.");
      }

    } catch (err) {
      console.error('Error fetching username:', err);
    }
  };

  return (
      <div className="App">
        <div className="page-container">
          <div className="navbar">
            <div className='navbar-box'>
              <Link to="/">Home</Link>
              <Link to="/register">Register</Link>
              <Link to="/login">Login</Link>
              <p>{username? username:"No User"}</p>
            </div>
          </div>
          <Routes>
            <Route path="/" element={<Home/>}/>
            <Route path="/register" element={<Register/>}/>
            <Route path="/login" element={<Login/>}/>
            <Route path="/email_verification/:verify_id" element={<EmailVerification/>}/>
            <Route path="/change_password/:username/:verify_id" element={<ChangePassword/>}/>
            <Route path="/post_product" element={<PostProduct/>}/>
          </Routes>
        </div>
    </div>
  );
}

export default App;
