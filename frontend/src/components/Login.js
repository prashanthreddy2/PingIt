import React, { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import styled from 'styled-components';
import axios from 'axios';
import { Button } from '@mui/material';

const Login = ({setUser}) => {
  const [userName, setUserName] = useState()
  const [password, setPassword] = useState()
  const navigate = useNavigate();
  const verifyLogin = async(e) => {
    e.preventDefault();
    const params = {
      "username": userName,
      "password": password, 
    }
    await axios.get('http://localhost:8080/login', {params})
      .then(resp => {
        setUser(resp.data.data.sender)
        sessionStorage.setItem("user", resp.data.data.sender)
        navigate('/')
      })
      .catch(err => {
        console.log(err)
      })
  }
  return (
    <MainContainer>
      <LoginForm>
        <h2>Login</h2>
        <InputField type="text" placeholder="Username" value={userName} onChange={(e) => setUserName(e.target.value)}/>
        <InputField type="password" placeholder="Password" value={password} onChange={(e) => setPassword(e.target.value)}/>
        <SubmitButton onClick={(e) => verifyLogin(e)}>Login</SubmitButton>
        <Button onClick={() => {navigate('/signup')}}>Signup</Button>
      </LoginForm>
    </MainContainer>
  );
};

export default Login;

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

const Logo = styled.img`
  width: 100px;
  height: 100px;
  margin-bottom: 20px;
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