import React, { useState } from 'react';
import { useNavigate, useParams } from 'react-router-dom';

function ChangePassword() {
  const {username, verify_id} = useParams();
  const [newPassword, setNewPassword] = useState('');
  const [error, setError] = useState('');
  const navigate = useNavigate()

  const handleChangePassword = async (e) => {
    e.preventDefault();

    try {
      const response = await fetch(`http://localhost:8000/api/forgot_password/${username}/${verify_id}`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        credentials: 'include', 
        body: JSON.stringify({
          "password": newPassword,
        }),
      });

      const result = await response.json();
  
      if (response.ok) {
        console.log('Password Successful Changed:', result);
        navigate('/login')
      } else {
        console.error('Email Verification Failed:', result);
        setError(response.status + " " + result.error)
      }

    } catch (err) {
      console.error('Error:', err.message);
      setError(err.message)
    }
  };

  return (
    <div>
      <h2>Change Password</h2>
      <form onSubmit={handleChangePassword}>
        <div>
          <label>New Password</label>
          <input
            type="password"
            name="newPassword"
            value={newPassword}
            onChange={(e) => setNewPassword(e.target.value)}
            placeholder="Enter new password"
            required
          />
        </div>
        {error && <p style={{ color: 'red' }}>{error}</p>}
        <button type="submit">Change Password</button>
      </form>
    </div>
  );
}

export default ChangePassword;
