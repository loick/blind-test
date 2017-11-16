import React, { Component } from 'react'
import { StyleSheet, Text, TextInput, ScrollView, Button } from 'react-native'

export default class Home extends Component {
  onSubmit() {
    console.log('coucou')
  }

  render() {
    return (
      <ScrollView contentContainerStyle={styles.container}>
        <Text>Blindary</Text>
        <TextInput style={styles.input} placeholder="Nickname" />
        <TextInput style={styles.input} autoCapitalize="characters" maxLength={5} autoCorrect={false} placeholder="Room number" />
        <Button onPress={this.onSubmit} title="Let's play" color="#841584" />
      </ScrollView>
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
  input: {
    margin: '3%',
    padding: '5% 3%',
    backgroundColor: '#FFF',
    borderRadius: 15,
    fontSize: 30,
    textAlign: 'center',
  }
})
