import React, { useEffect, useState } from 'react'
import { Navigate, useNavigate, useParams } from 'react-router-dom'

function EmailVerification() {
  const {verify_id} = useParams();
  const [code, setCode] = useState('');
  const [error, setError] = useState(null);
  const [verificationData, setVerificationData] = useState(null);
  const [timeLeft, setTimeLeft] = useState(null);
  const navigate = useNavigate()

  useEffect(() => {
    const fetchVerificationData = async () => {
      try {
        const response = await fetch(`${process.env.REACT_APP_API_HOST}/api/email_verification/${verify_id}`, {
          method: 'GET',
          headers: {
            'Content-Type': 'application/json',
          },
          credentials: 'include',
        });

        const result = await response.json();

        if (response.ok) {
          console.log(result)
          setVerificationData(result);
          setTimeLeft(result.time_left);
          console.log('Fetched Verification Data:', result);
        } else {
          setError(response.status + " " + result.error);
        }
      } catch (err) {
        console.error('Error:', err);
        setError(err.message);
      }
    };

    fetchVerificationData();
  }, [verify_id]);

  useEffect(() => {
    if (timeLeft === null) return;

    const timer = setInterval(() => {
      setTimeLeft((prev) => {
        if (prev <= 1) {
          clearInterval(timer);
          return 0;
        }
        return prev - 1;
      });
    }, 1000);

    return () => clearInterval(timer); // Cleanup the interval on component unmount
  }, [timeLeft]);

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

  const handleResendCode = async () => {
    try {
      const response = await fetch(`${process.env.REACT_APP_API_HOST}/api/resend_email_code/${verify_id}`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        credentials: 'include',
      });

      const result = await response.json();

      if (response.ok) {
        console.log('Email Verification Has Been Sent', result);
        setTimeLeft(result.time_left)
      } else {
        console.error('Failed to sent verification code:', result);
        setError(response.status + " " + result.error)
      }
    } catch (err) {
      console.error('Error:', err);
      setError(err.error)
    }
  }

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
      {!verificationData ? (
        <p>Loading...</p>
      ) : (
        <div>{timeLeft <= 0 ? <button onClick={handleResendCode}>Resend Code</button> : "time left: "+timeLeft}</div>
      )}
      {error && <p style={{ color: 'red' }}>{error}</p>}
    </div>
  )
}

export default EmailVerification