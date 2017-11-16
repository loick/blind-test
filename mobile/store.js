import { applyMiddleware, compose, createStore } from 'redux'
import createSagaMiddleware, { END } from 'redux-saga'

import rootSaga from 'sagas/'
import kairosReducers from 'reducers/'
import trackingMiddleware from 'utils/trackingMiddleware'

const sagaMiddleware = createSagaMiddleware()

export const enableDevTools = typeof window !== 'undefined' && window.devToolsExtension && __DEV__
  ? window.devToolsExtension()
  : x => x

const enhancer = compose(
  applyMiddleware(sagaMiddleware, trackingMiddleware),
  enableDevTools,
)

export const getStore = (state = {}) => {
  const store = createStore(
    kairosReducers,
    state,
    enhancer,
  )
  sagaMiddleware.run(rootSaga)

  store.runSaga = sagaMiddleware.run
  store.close = () => store.dispatch(END)

  return store
}

export default getStore
