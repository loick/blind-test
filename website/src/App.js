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

const serialize = (obj) => {
  var str = [];
  for(var p in obj)
    if (obj.hasOwnProperty(p)) {
      str.push(encodeURIComponent(p) + "=" + encodeURIComponent(obj[p]));
    }
  return str.join("&");
}

class App extends Component {
  componentDidMount() {
    const options = {
      client_id: '0bd627ea35b2418893459b4ee1568575',
      response_type: 'token',
      redirect_uri: 'http://127.0.0.1:3000',
    }

    fetch(`https://accounts.spotify.com/authorize?${serialize(options)}`)
      .then(response => response.json())
      .then(response => console.log(response))
      .catch(error => console.log(error))
  }

  render() {
    return (
      <MuiThemeProvider>
        <div>
          <AppBar title="BLINDARY" iconClassNameLeft="none" />
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
