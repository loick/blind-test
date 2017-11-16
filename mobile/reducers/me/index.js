import { ME_ADD_INFOS } from '../../actions/me'

export const initialState = {
  roomNumber: null,
  nickname: null,
}

export default function meReducer(state = initialState, action) {
  switch (action.type) {
    case ME_ADD_INFOS:
      return {
        ...state,
        nickname: action.nickname,
        roomNumber: action.roomNumber
      }

    default:
      return state
  }
}
