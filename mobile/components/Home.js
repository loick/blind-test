import React, { Component } from 'react'
import { connect } from 'react-redux'
import { StyleSheet, Text, TextInput, ScrollView, View, TouchableOpacity } from 'react-native'
import { API_URL } from '../config'
import { addUserInfos } from '../actions/me'

export class Home extends Component {
  state = {
    nickname: null,
    roomNumber: null,
    error: null,
  }

  onSubmit = async () => {
    this.setState({ error: null })
    try {
      // Get a room
      const createRoom = await fetch(`${API_URL}/roomnumber`, { method: 'GET' })
      const createRoomJson = await createRoom.json()
      const { roomNumber } = createRoomJson

      // Valid the room
      const response = await fetch(`${API_URL}/roomnumber`, { method: 'POST', body: JSON.stringify({ roomNumber }) })
      this.props.dispatch(addUserInfos(this.state))
    } catch(error) {
      this.setState({ error })
    }
  }

  render() {
    return (
      <ScrollView contentContainerStyle={styles.container}>
        <Text style={styles.title}>BLINDARY</Text>
        <View style={styles.inputContainer}>
          <TextInput style={styles.input} placeholder="Nickname" onChangeText={(nickname) => this.setState({ nickname })} />
          <TextInput style={styles.input} autoCapitalize="characters" onChangeText={(roomNumber) => this.setState({ roomNumber })} maxLength={5} autoCorrect={false} placeholder="Room number" />
        </View>
        <TouchableOpacity style={styles.submitButtonContainer} title="Let's play" onPress={this.onSubmit}>
          <Text style={styles.submitButton}>Play</Text>
        </TouchableOpacity>
        { this.state.error && <Text>{this.state.error}</Text> }
      </ScrollView>
    )
  }
}

const styles = StyleSheet.create({
  container: {
    flex: 1,
    backgroundColor: '#26D0CE',
    alignItems: 'center',
    justifyContent: 'space-around',
  },
  title: {
    fontSize: 20,
    color: '#FFF',
  },
  inputContainer: {
    width: '80%',
  },
  input: {
    margin: '3%',
    padding: '5% 3%',
    backgroundColor: '#FFF',
    borderRadius: 15,
    fontSize: 30,
    textAlign: 'center',
  },
  submitButtonContainer: {
    width: '80%',
    borderRadius: 50,
    backgroundColor: '#fff',
    borderWidth: 2,
    borderColor: '#1ECD97',
  },
  submitButton: {
    width: '100%',
    fontSize: 20,
    lineHeight: 40,
    borderRadius: 50,
    color: '#1ECD97',
    textAlign: 'center',
  }
})

export default connect()(Home)