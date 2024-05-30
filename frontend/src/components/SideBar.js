import React, { useEffect, useState } from 'react'
import { useNavigate, useLocation } from 'react-router-dom';
import styled from 'styled-components'
import AccountCircleIcon from '@mui/icons-material/AccountCircle';
import DonutLargeIcon from '@mui/icons-material/DonutLarge';
import ChatIcon from '@mui/icons-material/Chat';
import MoreVertIcon from '@mui/icons-material/MoreVert';
import IconButton from '@mui/material/IconButton';
import SearchOutlinedIcon from '@mui/icons-material/SearchOutlined';
import SideBarChat from './SideBarChat';
import axios from 'axios';

function SideBar({ searchQuery, setSearchQuery, rooms, setRooms }) {
  const [chatList, setChatList] = useState([]);
  const navigate = useNavigate();
  const location = useLocation();
  const fetchRelationList = async () => {
    let user = sessionStorage.getItem("user")
    let params = {
      sender: user
    }
    await axios.get('http://localhost:8080/fetchRelations', { params })
      .then(resp => {
        setRooms(new Set(resp.data.data))
      })
      .catch(err => {
        console.log(err)
      })
  }

  const fetchUserList = async () => {
    let user = sessionStorage.getItem("user")
    let params = {
      user: user
    }
    await axios.get('http://localhost:8080/userList', { params })
      .then(resp => {
        setChatList(resp.data.data)
        console.log(resp.data.data, "chatlist")
      })
      .catch(err => {
        console.log(err)
      })
  }

  const handleSearchChange = event => {
    setSearchQuery(event.target.value);
  };

  const checkUser = () => {
    let user = sessionStorage.getItem("user")
    if (user == undefined || user == null) {
      navigate('/login');
      window.location.reload()
    }
  }

  useEffect(() => {
    checkUser()
    fetchRelationList()
    fetchUserList()
    console.log("location changed")
  }, [])


  return (
    <MainComp>
      <Header>
        <AccountCircleIcon />
        <HeaderRight>
          <IconButton>
          </IconButton>
          <IconButton>
          </IconButton>
          <IconButton onClick={() => {sessionStorage.removeItem("user"); checkUser()} }>
            <MoreVertIcon />
          </IconButton>
        </HeaderRight>
      </Header>
      <Searchbar>
        <SearchContainer>
          <SearchOutlinedIcon />
          <Input
            type="text"
            placeholder="Search chats..."
            value={searchQuery}
            onChange={handleSearchChange}
          />
        </SearchContainer>
      </Searchbar>
      {
        searchQuery != '' && (
          <Chats>
            {chatList.filter(chat => chat.toLowerCase().startsWith(searchQuery.toLowerCase()))
              .map((chat) => (
                <SideBarChat roomName={chat} setSearchQuery={setSearchQuery} setRooms={setRooms} rooms={rooms} />
              ))}
          </Chats>
        )
      }
      <Chats>
        {searchQuery == '' &&
          Array.from(rooms).map(room => (
            <SideBarChat roomName={room} setSearchQuery={setSearchQuery} setRooms={setRooms} rooms={rooms} />
          ))
        }
      </Chats>
    </MainComp>
  )
}

export default SideBar

const MainComp = styled.div`
  display: flex;
  flex-direction: column;
  flex: 0.35;
`
const Header = styled.div`
  display: flex;
  justify-content: space-between;
  padding: 20px;
  border-right: 1px solid lightgray;
`
const HeaderRight = styled.div`
  display: flex;
  align-items: center;
  justify-content: space-between;
  min-width: 10vw;
`
const Searchbar = styled.div`
  display: flex;
  align-items: center;
  background-color: #f6f6f6;
  height: 39px;
  padding: 10px;
`
const SearchContainer = styled.div`
  display: flex;
  align-items: center;
  /* padding: 5px;   */
  width: 100%;
  height: 35px;
  border-radius: 20px;
`
const Chats = styled.div`
  flex: 1;
  background-color: white;
  /* overflow: auto; */
`

const Input = styled.input`
  width: 100%;
  padding: 10px;
  font-size: 16px;
  border: none;
  border-radius: 20px;
  outline: none;
  box-sizing: border-box;
`;