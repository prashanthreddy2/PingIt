import React, { useEffect, useState, useRef } from 'react'
import styled from 'styled-components'
import AccountCircleIcon from '@mui/icons-material/AccountCircle';
import MoreVertIcon from '@mui/icons-material/MoreVert';
import IconButton from '@mui/material/IconButton';
import MoodOutlinedIcon from '@mui/icons-material/MoodOutlined';
import { useParams } from 'react-router-dom';
import axios from 'axios';
import { useNavigate } from 'react-router-dom';


function ChatBar({ setSearchQuery, setRooms, rooms }) {
  const [inputMessage, setInputMessage] = useState("")
  const { roomId } = useParams()
  const [currentUser, setCurrentUser] = useState()
  const [webSocket, setWebSocket] = useState(null)
  const now = new Date();
  const [messages, setMessages] = useState([])
  const messagesEndRef = useRef(null);
  const navigate = useNavigate();
  useEffect(() => {
    let user = sessionStorage.getItem("user")
    if (user == undefined || user == null) {
      navigate('/login');
    }
  }, [])

  useEffect(() => {
    scrollToBottom();
  }, [messages]);


  useEffect(() => {
    let user = sessionStorage.getItem("user")
    setCurrentUser(user)
    fetchMessages()
    setSearchQuery('')
    console.log(rooms)
    setRooms((prev) => {
      const updatedRooms = new Set(prev);
      updatedRooms.add(roomId)
      console.log(updatedRooms, "urooms")
      return updatedRooms;
    })
  }, [roomId])

  const scrollToBottom = () => {
    if (messagesEndRef.current) {
      messagesEndRef.current.scrollIntoView({ behavior: 'smooth' });
    }
  };

  const sendMessage = async (e) => {
    e.preventDefault();
    if (inputMessage == "") {
      return
    }
    const currentMessage = {
      "sender": currentUser,
      "message": inputMessage,
      "timestamp": now.toUTCString(),
      "reciever": roomId
    }
    webSocket.send(JSON.stringify(currentMessage))
    setMessages(messages => [...messages, currentMessage])
    setInputMessage("")
  }

  const fetchMessages = async () => {
    const params = {
      "sender": sessionStorage.getItem("user"),
      "reciever": roomId,
    }
    await axios.get('http://localhost:8080/fetchMessages', { params })
      .then(res => {
        setMessages(res.data.data)
      })
      .catch(err => {
        console.log(err, "Messages fetch err")
      })
  }

  useEffect(() => {
    let user = sessionStorage.getItem("user")
    setCurrentUser(user)
    let ws = new WebSocket(`ws://localhost:8080/ws?clientID=${user}`)
    setWebSocket(ws)
    ws.onmessage = (event) => {
      const Data = JSON.parse(event.data)
      const currentMessage = {
        "sender": Data.Sender,
        "message": Data.Message,
        "timestamp": Data.Timestamp,
        "reciever": Data.Reciever
      }
      setMessages(messages => [...messages, currentMessage])
    }
    ws.onclose = () => {
      console.log("Connection closed, Please check internet!")``
    }
    // return () => {
    //   ws.close();
    // }
    // eslint-disable-next-line
  }, [])

  return (
    <Container>
      <ChatHeader>
        <AccountCircleIcon />
        <HeaderInfo>
          <h3>{roomId}</h3>
        </HeaderInfo>
        <HeaderRight>
          <IconButton>
            <MoreVertIcon />
          </IconButton>
        </HeaderRight>
      </ChatHeader>
      <ChatBody>
        {
          messages.map(message => (
            message.sender === currentUser ? (
              <SenderMessage>
                <ChatName>{message.sender}</ChatName>
                {message.message}
                <ChatTime> {message.timestamp} </ChatTime>
              </SenderMessage>
            ) : (
              <RecieverMessage>
                <ChatName>{message.sender}</ChatName>
                {message.message}
                <ChatTime>{message.timestamp}</ChatTime>
              </RecieverMessage>
            )
          ))
        }
        <div ref={messagesEndRef} />
      </ChatBody >
      <ChatFooter>
        <MoodOutlinedIcon />
        <form onSubmit={sendMessage}>
          <input type='text' value={inputMessage} onChange={(e) => setInputMessage(e.target.value)} />
        </form>
      </ChatFooter>
    </Container >
  )
}

export default ChatBar

const Container = styled.div`
  flex: 0.65;
  display: flex;
  flex-direction: column;
`
const ChatHeader = styled.div`
  padding: 20px;
  display: flex;
  align-items: center;
  border-bottom: 1px solid lightgray;
`
const HeaderInfo = styled.div`
  flex: 1;
  padding-left: 20px;
`
const HeaderRight = styled.div`
  display: flex;
  justify-content: space-between;
  min-width: 100px;
`
const ChatBody = styled.div`
  flex: 1;
  background-image: url("https://w0.peakpx.com/wallpaper/744/548/HD-wallpaper-whatsapp-ma-doodle-pattern-thumbnail.jpg");
  background-repeat: repeat;
  background-position: center;
  padding: 30px;
  overflow: auto;
`
const RecieverMessage = styled.p`
  position: relative;
  font-size: 16px;
  padding: 10px;
  background-color: #ffffff;
  border-radius: 10px;
  width: fit-content;
  margin-bottom: 15px;
`
const SenderMessage = styled.p`
  position: relative;
  font-size: 16px;
  padding: 10px;
  background-color: #ffffff;
  border-radius: 10px;
  width: fit-content;
  margin-left: auto;
  background-color: #dcf8c6;
  margin-bottom: 15px;
`
const ChatName = styled.span`
  position: absolute;
  top: -15px;
  font-weight: 800;
  font-size: xx-small;
`
const ChatTime = styled.span`
  margin-left: 10px;
  font-size: xx-small;
`
const ChatFooter = styled.div`
  display: flex;
  justify-content: space-between;
  align-items: center;
  height: 62px;
  border-top: 1px solid lightgray;
  form{
    flex: 1;
    display: flex;
    input{
      flex: 1;
      border-radius: 30px;
      padding: 10px;
      border: none;
    }
    button{
      display: none;
    }
  }
`