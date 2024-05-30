import React from 'react'
import styled from 'styled-components'
import PersonIcon from '@mui/icons-material/Person';
import { Link } from 'react-router-dom';

function SideBarChat({ roomName, lastMessage, setSearchQuery, rooms, setRooms}) {
    return (
        <StyledLink to={`/rooms/${roomName}`}>
            <Component>
                <PersonIcon />
                <Info>
                    <h2>{roomName}</h2>
                    <p>{lastMessage}</p>
                </Info>
            </Component>
        </StyledLink>
    )
}

export default SideBarChat

const StyledLink = styled(Link)`
    text-decoration: none;
    color: black;
`
const Component = styled.div`
    display: flex;
    padding: 20px;
    cursor: pointer;
    border-bottom: 1px solid #f6f6f6;
    &:hover{
        background-color: #ebebeb;
    }
`
const Info = styled.div`
    margin-left: 15px;
`