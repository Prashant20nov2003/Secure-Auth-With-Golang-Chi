import React, { useState } from 'react'
import { useNavigate } from 'react-router-dom';

function Login() {
  const [formData, setFormData] = useState({
    username: '',
    password: ''
  });
  const [error, setError] = useState('')
  const navigate = useNavigate();

  const handleChange = (e) => {
    setFormData({
      ...formData,
      [e.target.name]: e.target.value,
    });
  };

  const handleSubmit = async (e) => {
    e.preventDefault();
    console.log('Form submitted:', formData);
    
    try {
      const response = await fetch(`${process.env.REACT_APP_API_HOST}/api/login`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({
          "username": formData.username,
          "password": formData.password
        }),
        credentials: 'include',
      });

      const result = await response.json(); // Parse the JSON response

      if (response.ok) {
        console.log('Successful:', result);
        if (result.message === "Verify link first"){
          navigate(`/email_verification/${result.emailverify_id}`)
        }
        else{
          navigate(`/`)
        }
      } else {
        console.error('Login failed:', result);
        setError(response.status + " " + result.error)
      }
    } catch (err) {
      console.error('Error:', err);
      setError(err.toString())
    }
  };

  const handleForgotPassword = async (e) => {
    e.preventDefault();

    if (formData.username !== ""){
      try {
        const response = await fetch(`${process.env.REACT_APP_API_HOST}/api/forgot_password/${formData.username}`, {
          method: 'POST',
          headers: {
            'Content-Type': 'application/json',
          },
          body: JSON.stringify({
          }),
        });
  
        const result = await response.json(); // Parse the JSON response
  
        if (response.ok) {
          console.log('Successful:', result);
          navigate(`/email_verification/${result.emailverify_id}`)
        } else {
          console.error('Failed:', result);
          setError(response.status + " " + result.error)
        }
      } catch (err) {
        console.error('Error:', err);
        setError(err.message || "Can't connect to API")
      }
    }else{
      setError("Input your username to recover your password")
    }
  }

  return (
    <div>
      <h2>Login</h2>
      <form onSubmit={handleSubmit}>
        <div>
          <label>Username</label>
          <input
            type="text"
            name="username"
            value={formData.username}
            onChange={handleChange}
            placeholder="Enter your username"
            required
          />
        </div>
        <div>
          <label>Password</label>
          <input
            type="password"
            name="password"
            value={formData.password}
            onChange={handleChange}
            placeholder="Enter your password"
            required
          />
        </div>
        {error && <p style={{ color: 'red' }}>{error}</p>}
        <button type="submit">Login</button>
      </form>
      <button onClick={handleForgotPassword}>Forgot Password</button>
    </div>
  );
}

export default Login