import React, { useState } from 'react'
import { Navigate, useNavigate, useParams } from 'react-router-dom'

function EmailVerification() {
  const {verify_id} = useParams();
  const [code, setCode] = useState('');
  const [error, setError] = useState(null);
  const navigate = useNavigate()

  const handleInputChange = async (e) => {
    const value = e.target.value;
    if (/^\d*$/.test(value) && value.length <= 6) { // Only allow numbers and restrict to 6 digits
      setCode(value);
    }

    if (value.length === 6){
      try {
        const response = await fetch(`${process.env.REACT_APP_API_HOST}/api/email_verification/${verify_id}`, {
          method: 'POST',
          headers: {
            'Content-Type': 'application/json',
          },
          credentials: 'include',
          body: JSON.stringify({
            "verif_code": value,
          }),
        });
  
        const result = await response.json();
  
        if (response.ok) {
          console.log('Email Verification Successful:', result);
          if (result.used_for === 'Verify Email'){
            navigate('/')
          }
          else if (result.used_for === 'Change Password'){
            navigate(`/change_password/${result.username}/${result.emailverify_id}`)
          }
        } else {
          console.error('Email Verification Failed:', result);
          setError(response.status + " " + result.error)
        }
      } catch (err) {
        console.error('Error:', err);
        setError(err.error)
      }
    }
  };

  return (
    <div>
      <h1>Email Verification</h1>
      <input
        type="text"
        value={code}
        onChange={handleInputChange}
        maxLength={6}
        placeholder="Enter 6-digit code"
      />
      {error && <p style={{ color: 'red' }}>{error}</p>}
    </div>
  )
}

export default EmailVerification