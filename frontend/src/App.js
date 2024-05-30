import './App.css';
import SideBar from './components/SideBar';
import ChatBar from './components/ChatBar';
import Login from './components/Login';
import styled from 'styled-components';
import SignUp from './components/SignUp';
import { BrowserRouter, Routes, Route, useNavigate } from 'react-router-dom';
import { useEffect, useState } from 'react';

function App() {
  const [user, setUser] = useState()
  const [searchQuery, setSearchQuery] = useState('')
  const [rooms, setRooms] = useState(new Set())
  const navigate = useNavigate()
  useEffect(() => {
    if (user == undefined || user == null) {
      navigate('/login');
    }
  }, [user])
  return (
    <div className="App">
      {!user ? (
        <>
          <Routes>
            <Route path='/login' element={<Login setUser={setUser} />} />
            <Route path='/signup' element={<SignUp setUser={setUser} />} />
          </Routes>
        </>
      ) : (

        <MainComp>
          <SideBar searchQuery={searchQuery} setSearchQuery={setSearchQuery} setRooms={setRooms} rooms={rooms} />
          <Routes>
            <Route path="/rooms/:roomId" element={<ChatBar searchQuery={searchQuery} setSearchQuery={setSearchQuery} setRooms={setRooms} rooms={rooms}/>} />
          </Routes>
        </MainComp>
      )
      }
    </div>
  );
}

export default App;

const MainComp = styled.div`
  display: flex;
  background-color: #ededed;
  margin-top: -50px;
  height: 90vh;
  width: 90vw;
  box-shadow: -1px 4px 20px -6px rgba(0,0,0,0.7);
`
