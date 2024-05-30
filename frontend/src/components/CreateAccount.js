import React from 'react'
import styled from 'styled-components';
import { Link } from 'react-router-dom';

function CreateAccount() {
  return (
    <MainContainer>
      <CreateAccountForm>
        {/* <Logo src={logo} alt="Logo" /> */}
        <h2>Create Account</h2>
        <InputField type="text" placeholder="Username" />
        <InputField type="password" placeholder="Password" />
        <InputField type="password" placeholder="Confirm Password" />
        <SubmitButton>Create Account</SubmitButton>
        <LoginLink to="/login">Back to Login</LoginLink>
      </CreateAccountForm>
    </MainContainer>
  )
}

export default CreateAccount

const MainContainer = styled.div`
  display: flex;
  justify-content: center;
  align-items: center;
  height: 100vh;
`;

const CreateAccountForm = styled.form`
  width: 300px;
  padding: 20px;
  border: 1px solid #ccc;
  border-radius: 5px;
  box-shadow: 0 0 10px rgba(0, 0, 0, 0.1);
  text-align: center;
`;
const InputField = styled.input`
  width: 100%;
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

const CreateAccountLink = styled(Link)`
  display: block;
  margin-top: 10px;
  color: #007bff;
  text-decoration: none;
`;
const LoginLink = styled(Link)`
  display: block;
  margin-top: 10px;
  color: #007bff;
  text-decoration: none;
`;