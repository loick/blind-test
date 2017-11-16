import React, { Component } from 'react'
import { StyleSheet, Text, TextInput, View, StatusBar } from 'react-native'
import Home from './components/Home'
import Play from './components/Play'
import { createStore } from 'redux'
import { Provider, connect } from 'react-redux'
import reducers from './reducers'

const store = createStore(reducers)

StatusBar.setHidden(true)

export const App = ({ roomNumber }) => roomNumber ? <Play /> : <Home />

const mapStateToProps = ({ me: { roomNumber } }) => ({ roomNumber })
const Test = connect(mapStateToProps)(App)

export default class AppContainer extends Component {
  render() {
    return (
      <Provider store={store}>
        <Test />
      </Provider>
    )
  }
}

const styles = StyleSheet.create({
  container: {
    flex: 1,
    backgroundColor: '#fff',
    alignItems: 'center',
    justifyContent: 'center',
  },
})


