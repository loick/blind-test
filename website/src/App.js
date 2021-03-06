import React, { Component } from 'react'
import MuiThemeProvider from 'material-ui/styles/MuiThemeProvider'
import AppBar from 'material-ui/AppBar'
import {
  Table,
  TableBody,
  TableHeader,
  TableHeaderColumn,
  TableRow,
  TableRowColumn,
} from 'material-ui/Table'
import { API_URL, SOCKET_URL } from './config'
import spotify from './spotify.svg'

const serialize = (obj) => {
  var str = [];
  for(var p in obj)
    if (obj.hasOwnProperty(p)) {
      str.push(encodeURIComponent(p) + "=" + encodeURIComponent(obj[p]))
    }
  return str.join("&")
}

const options = {
  client_id: '0bd627ea35b2418893459b4ee1568575',
  response_type: 'token',
  redirect_uri: window.location.origin + '/spotify-connect',
}

class App extends Component {
  ws = null

  state = {
    access_token: null,
    expires_in: null,
    token_type: null,
    roomNumber: null,
    answers: [],
  }

  componentDidMount() {
    if (window.location.hash) {
      const hash = window.location.hash.substring(1, window.location.hash.length)
      this.setState(hash.split('&').reduce((acc, curr) => {
        const [key, value] = curr.split('=')
        return ({ ...acc, [key]: value })
      }, {}), this.getRoom)
    }
  }

  componentDidUpdate() {
    if (!this.ws && this.state.roomNumber) {
      this.ws = new WebSocket(`${SOCKET_URL}${this.state.roomNumber}/1`)
      this.ws.onmessage = this.onReceiveAnswer
    }
  }

  onReceiveAnswer(event) {
    console.log(event)
  }

  async getRoom() {
    try {
      const createRoom = await fetch(`${API_URL}/roomnumber`, { method: 'GET' })
      const createRoomJson = await createRoom.json()
      const { roomNumber } = createRoomJson
      this.setState({ roomNumber })
      this.sendToken()
    } catch(e) {
      console.log(e)
    }
  }

  async sendToken() {
    try {
      const sendToken = await fetch(`${API_URL}/tokens`,
        {
          method: 'POST',
          body: JSON.stringify({
            roomNumber: this.state.roomNumber,
            token: this.state.access_token,
          })
        }
      )
    } catch(e) {
      console.log(e)
    }
  }

  render() {
    return (
      <MuiThemeProvider>
        <div>
          <AppBar
            title={`BLINDARY ${this.state.roomNumber ? `(${this.state.roomNumber})` : ''}`}
            iconElementLeft={<a href={`https://accounts.spotify.com/authorize?${serialize(options)}`}><img fill="#FFF" src={spotify} /></a>}
          />
          <Table>
            <TableHeader displaySelectAll={false}>
              <TableRow>
                <TableHeaderColumn>Nickname</TableHeaderColumn>
                <TableHeaderColumn>Points</TableHeaderColumn>
              </TableRow>
            </TableHeader>
            <TableBody displayRowCheckbox={false}>
              <TableRow>
                <TableRowColumn>Loick</TableRowColumn>
                <TableRowColumn>6</TableRowColumn>
              </TableRow>
              <TableRow>
                <TableRowColumn>Remy</TableRowColumn>
                <TableRowColumn>0</TableRowColumn>
              </TableRow>
            </TableBody>
          </Table>
        </div>
      </MuiThemeProvider>
    )
  }
}

export default App 
