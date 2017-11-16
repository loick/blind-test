import React, { Component } from 'react'
import { StyleSheet, Text, TextInput, View } from 'react-native'
import Home from './components/Home'

export default class App extends Component {
  render() {
    return <Home />
  }
}

const styles = StyleSheet.create({
  container: {
    flex: 1,
    backgroundColor: '#fff',
    alignItems: 'center',
    justifyContent: 'center',
  },
});
