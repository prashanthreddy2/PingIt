import React, { useEffect, useState } from 'react';
import styled from 'styled-components';
import axios from 'axios';
import { useNavigate } from 'react-router-dom';

function SignUp() {
  const [userName, setUserName] = useState()
  const [password, setPassword] = useState()
  const [errorMsg, setErrorMsg] = useState()
  const navigate = useNavigate()
  const handleSubmit = async (e) => {
    e.preventDefault();
    const params = {
      "username": userName,
      "password": password,
    }
    await axios.get('http://localhost:8080/signup', { params })
      .then(res => {
        setErrorMsg(null)
        navigate('/login')
      })
      .catch(err => {
        console.log(err.response.data.message, "<- error creating account")
        setErrorMsg(err.response.data.message)
      })
      setUserName('')
      setPassword('')
  }
  useEffect(() => {
    const timeoutId = setTimeout(() => {
      setErrorMsg(null);
    }, 5000);
    return () => clearTimeout(timeoutId);
  }, [errorMsg])
  return (
    <MainContainer>
      <LoginForm>
        <h2>Signup</h2>
        <InputField type="text" placeholder="Username" value={userName} onChange={(e) => setUserName(e.target.value)}  required/>
        <InputField type="password" placeholder="Password" value={password} onChange={(e) => setPassword(e.target.value)} required />
        {errorMsg && <h5> {errorMsg} </h5>}
        <SubmitButton onClick={(e) => handleSubmit(e)}>Signup</SubmitButton>
        {/* <Button onClick={() => {navigate('/login')}}>Login</Button> */}
      </LoginForm>
    </MainContainer>
  )
}

export default SignUp

const MainContainer = styled.div`
  display: flex;
  justify-content: center;
  align-items: center;
  height: 100vh;
`;

const LoginForm = styled.form`
  width: 300px;
  padding: 20px;
  border: 1px solid #ccc;
  border-radius: 5px;
  box-shadow: 0 0 10px rgba(0, 0, 0, 0.1);
  text-align: center;
`;

const InputField = styled.input`
  width: 90%;
  margin-bottom: 15px;
  padding: 10px;
  border: 1px solid #ccc;
  border-radius: 5px;
  outline: none;
`;

const SubmitButton = styled.button`
  width: 100%;
  padding: 10px;
  background-color: #007bff;
  color: #fff;
  border: none;
  border-radius: 5px;
  cursor: pointer;
  transition: background-color 0.3s;
  margin-bottom: 10px;

  &:hover {
    background-color: #0056b3;
  }
`;