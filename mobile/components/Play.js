import React, { Component } from 'react'
import { connect } from 'react-redux'
import { StyleSheet, Text, TextInput, View, TouchableOpacity } from 'react-native'
import { API_URL } from '../config'
import { addUserInfos } from '../actions/me'

export class Play extends Component {
  state = {
    word: '',
  }

  updateWord = (word) => {
    this.setState((prevState) => ({ word }))
  }

  render() {
    const NUM_LETTERS = 6
    const placeholder = Array.from({ length: NUM_LETTERS }, () => '_').join(' ')
    return (
      <View style={styles.container}>
        <Text style={styles.title}>What is this song? ({NUM_LETTERS}) letters</Text>
        <TextInput
          style={styles.input}
          autoCapitalize="characters"
          onChangeText={this.updateWord}
          maxLength={NUM_LETTERS}
          autoCorrect={false}
          placeholder={placeholder}
        />
        <TouchableOpacity style={styles.submitButtonContainer} title="Let's play" onPress={this.onSubmit}>
          <Text style={styles.submitButton}>Send</Text>
        </TouchableOpacity>
      </View>
    )
  }
}

const styles = StyleSheet.create({
  container: {
    flex: 1,
    backgroundColor: '#26D0CE',
    alignItems: 'center',
    justifyContent: 'center',
  },
  title: {
    fontSize: 20,
    color: '#FFF',
  },
  input: {
    width: '80%',
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

const mapStateToProps = ({ me: { roomNumber } }) => ({ roomNumber })

export default connect(mapStateToProps)(Play)