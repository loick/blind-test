export const ME_ADD_INFOS = 'ME_ADD_INFOS'

export function addUserInfos(infos) {
  return { type: ME_ADD_INFOS, ...infos }
}
